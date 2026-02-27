# Аукцион

## Основные особенности
- Ограниченные контексты в виде модулей: аукцион, планировщик, уведомления, метрика
- Event Sourcing для агрегатов аукцион и лот
- CQRS: разделение команд и запросов, проекция для аукциона
- Event Driven: интеграция контекстов через события, слабая связанность
- Фронтенд получает события по websocket: старт, новая ставка, завершение

## Поток событий
- AuctionCreated -> Scheduler создает задачи start/end/timeout
- AuctionStartTimeReached -> Auction стартует
- AuctionStarted -> Notification уведомляет frontend
- BidPlaced -> Notification уведомляет frontend
- AuctionTimeoutReached -> Auction завершается
- AuctionEndTimeReached -> Auction завершается
- AuctionCompleted -> Notification уведомляет frontend

## Стек
#### Backend - Go, PostgreSQL, Kafka
#### Frontend - React

## Запуск
### Требования
- docker и docker compose
- golang-migrate для миграций

```bash
# Запуск всех сервисов
docker compose up -d

# Применение миграций
make migrate
```

## Сервисы
- Backend API:  http://localhost:8081
- Frontend:  http://localhost:3000
- Kafka UI:  http://localhost:8080

## API

#### Коллекция postman в файле
```
auction_postman.json
```


#### Auction
```
POST   /lots                    # создать лот
POST   /auctions                # создать аукцион
GET    /auctions/{id}           # получить аукцион по id
POST   /auctions/{id}/bids      # сделать ставку
POST   /auctions/{id}/cancel    # отменить аукцион
```

#### WebSocket для frontend
```
GET    /ws?auction_id={id}&user_id={id}  # подписаться на уведомления
```

#### Metrics
```
GET    /metrics/global                      # глобальные метрики
GET    /metrics/auction                     # метрики всех аукционов
GET    /metrics/auction?auction_id={id}     # метрики аукциона
GET    /metrics/user                        # метрики пользователей
GET    /metrics/user?user_id={id}           # метрики пользователя
```
