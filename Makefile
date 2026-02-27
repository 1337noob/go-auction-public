migrate:
	migrate -path migrations -database "postgresql://auction:auction@127.0.0.1:5432/auction?sslmode=disable" -verbose up

rollback:
	migrate -path migrations -database "postgresql://auction:auction@127.0.0.1:5432/auction?sslmode=disable" -verbose down

make-migration:
	migrate create -ext sql -dir migrations -seq $(name)

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto
.PHONY: migrate rollback make-migration proto
