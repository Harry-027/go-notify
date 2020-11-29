download:
	go mod download
setup:
	docker-compose -f docker-compose.yml up -d
server:
	cd api-server && go run main.go
cronjob:
	cd cron && go run main.go
consumer:
	cd kafka-consumer && go run main.go
cli-go:
	cd cli/go-notify && go install
load-testing:
	cd loadtesting && docker-compose -f docker-compose.yml run k6 run ./tests/loadtests.js
monitor:
	cd monitoring && docker-compose -f docker-compose.yml up -d --build
