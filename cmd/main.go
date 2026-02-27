package main

import (
	"database/sql"
	"log"
	aucApi "main/auction/interfaces/api"
	aucModule "main/auction/module"
	metrApi "main/metrics/interfaces/api"
	metrModule "main/metrics/module"
	notifApi "main/notification/interfaces/api"
	notifModule "main/notification/module"
	"main/pkg"
	"main/pkg/integration_event_bus"
	"main/pkg/outbox"
	schedModule "main/scheduler/module"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "auction")
	dbPassword := getEnv("DB_PASSWORD", "auction")
	dbName := getEnv("DB_NAME", "auction")
	dbURL := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("failed to ping database", "error", err)
	}
	sqlTxManager := pkg.NewSQLTransactionManager(db)

	kafkaBroker := getEnv("KAFKA_BROKER", "localhost:9093")

	configOutbox := integration_event_bus.KafkaConfig{
		Brokers:            []string{kafkaBroker},
		TopicPrefix:        "integration-events",
		ConsumerGroup:      "outbox-consumer-group",
		TransactionManager: sqlTxManager,
	}
	kafkaIntegrationEventBusOutbox, err := integration_event_bus.NewKafkaIntegrationEventBus(configOutbox)
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaIntegrationEventBusOutbox.Close()

	configAuction := integration_event_bus.KafkaConfig{
		Brokers:            []string{kafkaBroker},
		TopicPrefix:        "integration-events",
		ConsumerGroup:      "auction-consumer-group",
		TransactionManager: sqlTxManager,
	}
	kafkaIntegrationEventBusAuction, err := integration_event_bus.NewKafkaIntegrationEventBus(configAuction)
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaIntegrationEventBusAuction.Close()

	configScheduler := integration_event_bus.KafkaConfig{
		Brokers:            []string{kafkaBroker},
		TopicPrefix:        "integration-events",
		ConsumerGroup:      "scheduler-consumer-group",
		TransactionManager: sqlTxManager,
	}
	kafkaIntegrationEventBusScheduler, err := integration_event_bus.NewKafkaIntegrationEventBus(configScheduler)
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaIntegrationEventBusScheduler.Close()

	configNotification := integration_event_bus.KafkaConfig{
		Brokers:            []string{kafkaBroker},
		TopicPrefix:        "integration-events",
		ConsumerGroup:      "notification-consumer-group",
		TransactionManager: sqlTxManager,
	}
	kafkaIntegrationEventBusNotification, err := integration_event_bus.NewKafkaIntegrationEventBus(configNotification)
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaIntegrationEventBusNotification.Close()

	configMetrics := integration_event_bus.KafkaConfig{
		Brokers:            []string{kafkaBroker},
		TopicPrefix:        "integration-events",
		ConsumerGroup:      "metrics-consumer-group",
		TransactionManager: sqlTxManager,
	}
	kafkaIntegrationEventBusMetrics, err := integration_event_bus.NewKafkaIntegrationEventBus(configMetrics)
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaIntegrationEventBusMetrics.Close()

	auctionOutboxTableName := "auction.outbox"
	auctionPostgresOutboxRepo := outbox.NewPostgresOutboxRepository(auctionOutboxTableName)

	schedulerOutboxTableName := "scheduler.outbox"
	schedulerPostgresOutboxRepo := outbox.NewPostgresOutboxRepository(schedulerOutboxTableName)

	marshaller := integration_event_bus.NewEventMarshaller()
	outboxPublisher := outbox.NewSimpleEventPublisher(kafkaIntegrationEventBusOutbox, marshaller)
	outboxInterval := time.Second * 1
	outboxLimit := 100

	auctionPostgresOutboxWorker := outbox.NewPostgresOutboxWorker(sqlTxManager, auctionPostgresOutboxRepo, outboxPublisher, outboxInterval, outboxLimit)
	auctionPostgresOutboxWorker.StartOutboxWorker()

	schedulerPostgresOutboxWorker := outbox.NewPostgresOutboxWorker(sqlTxManager, schedulerPostgresOutboxRepo, outboxPublisher, outboxInterval, outboxLimit)
	schedulerPostgresOutboxWorker.StartOutboxWorker()

	auctionModule := aucModule.NewAuctionModule(kafkaIntegrationEventBusAuction, auctionPostgresOutboxRepo, sqlTxManager, db)
	auctionHttpHandler := aucApi.NewHttpHandler(auctionModule)
	_ = schedModule.NewSchedulerModule(kafkaIntegrationEventBusScheduler, schedulerPostgresOutboxRepo, sqlTxManager)
	notifyModule := notifModule.NewNotificationModule(kafkaIntegrationEventBusNotification)
	notifyHttpHandler := notifApi.NewHttpHandler(notifyModule)
	metricsModule := metrModule.NewMetricsModule(kafkaIntegrationEventBusMetrics)
	metricsHttpHandler := metrApi.NewHttpHandler(metricsModule.GetQueryHandler())

	kafkaIntegrationEventBusAuction.StartConsuming()
	kafkaIntegrationEventBusScheduler.StartConsuming()
	kafkaIntegrationEventBusNotification.StartConsuming()
	kafkaIntegrationEventBusMetrics.StartConsuming()

	router := mux.NewRouter()

	// Lots endpoints
	router.HandleFunc("/lots", auctionHttpHandler.CreateLot).Methods("POST")

	// Auctions endpoints
	router.HandleFunc("/auctions", auctionHttpHandler.CreateAuction).Methods("POST")
	router.HandleFunc("/auctions/{id}", auctionHttpHandler.FindAuctionById).Methods("GET")
	router.HandleFunc("/auctions/{id}/bids", auctionHttpHandler.PlaceBid).Methods("POST")
	router.HandleFunc("/auctions/{id}/cancel", auctionHttpHandler.CancelAuction).Methods("POST")

	// WebSocket endpoints
	router.HandleFunc("/ws", notifyHttpHandler.HandleWebSocket)

	// Metrics endpoints
	router.HandleFunc("/metrics/global", metricsHttpHandler.GetGlobalMetrics).Methods("GET")
	router.HandleFunc("/metrics/auctions", metricsHttpHandler.GetAllAuctionsMetrics).Methods("GET")
	router.HandleFunc("/metrics/auction", metricsHttpHandler.GetAuctionMetrics).Methods("GET")
	router.HandleFunc("/metrics/users", metricsHttpHandler.GetAllUsersMetrics).Methods("GET")
	router.HandleFunc("/metrics/user", metricsHttpHandler.GetUserMetrics).Methods("GET")

	// pprof endpoints for profiling
	//router.HandleFunc("/debug/pprof/", pprof.Index)
	//router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	//router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	//router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	//router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	go func() {
		log.Println("time:", time.Now().String())
	}()

	log.Println("Сервер запущен на :8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
