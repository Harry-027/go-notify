# go-notify
---
An email automation solution written in [Golang](https://golang.org/).
It facilitate users to register, send & schedule custom HTML mails for their clients.

## Built using rich tech-stack:
---
* Api-server built using [Go-fiber](https://gofiber.io/).
* Using [Apache Kafka](https://kafka.apache.org/) as a message broker.
* Using [Postgres](https://www.postgresql.org/) as database.
* Using [Redis](https://redis.io/) as cache.
* Client CLI built using [Cobra](https://github.com/spf13/cobra).
* Using [Mailgun](https://www.mailgun.com/) as Email service.
* Using [K6](https://k6.io/) for load testing.
* Using [Prometheus](https://prometheus.io/) & [Grafana](https://grafana.com/) for Api-server monitoring.

## Architecture diagram:
---

![Arch_Diagram](https://github.com/Harry-027/go-notify/blob/master/snapshots/system_diagram.png "Arch_Diagram")

## DB Schema:
---

![DB_Schema](https://github.com/Harry-027/go-notify/blob/master/snapshots/dbSchema.PNG "DB_Schema")

## Features Included:
---
* Authentication & authorization using [JWT](https://jwt.io/).
* Swagger included, built using [Swago](https://github.com/swaggo/swag)
* Full proof auth features - signup, login, update password, forgot password, logout.
* Cors, Helmet, Api-rate limiter included as middleware from security perspective.
* Users can register their clients & custom HTML templates.
* Mail scheduling (daily, weekly, monthly) using cron jobs.
* Subscription & Payment (payment has been stubbed for now and can be replaced with any suitable payment gateway).
* Api-server can be consumed by any client-side API, however for easy interaction - CLI (built using [Cobra](https://github.com/spf13/cobra)) has been included.

## Installation & setup :-
---
* Go,Docker,Docker compose & Make should be pre-installed.
* Clone the repository: `git clone https://github.com/Harry-027/go-notify.git`.
* Run the command `make download` (this will install go modules).
* Create a new file .env under root directory & copy the env variables from .sample-env.
  (Note that mailgun env variables should be replaced with original credentials. Rest may remain untouched)
* Run the command `make setup` (this will start the required docker containers - postgres, redis, apache kafka & zoo-keeper).
* Run the command `docker ps` to ensure all the four containers are up & running.
* Open a new terminal & run the command `make server` to spin up the api server.
* Open a new terminal & run the command `make consumer` to spin up the kafka consumer.
* Open a new terminal & run the command `make cronjob` to start the cron processes.
* Open a new terminal & run the command `cli-go`. This will install the go-notify cli on your machine.
* Cli is now ready to operate. Run the command `go-notify --help` to explore various commands.

![CLI](https://github.com/Harry-027/go-notify/blob/master/snapshots/cli_snapshot.PNG "CLI")

## Swagger :-
---
* Once the server starts listening on port 3001, visit http://localhost:3001/swagger/ on browser for swagger definition.

![Swagger](https://github.com/Harry-027/go-notify/blob/master/snapshots/swagger_snapshot.PNG "Swagger")

## Monitoring (using Prometheus & Grafana):-
---
* Before spinning up Prometheus & Grafana for monitoring, replace the HOST_IP variable (under monitoring/prometheus/config.yml) with your machine IP.
* Run the command `make monitor` to start Api-server monitoring.
* Once the containers - Prometheus & Grafana are up, visit http://localhost:3000 on browser for Grafana dashboard.
* Default credentials for Grafana: username - 'admin' , password - 'admin'
* Once logged into Grafana, visit settings to select prometheus data source as target to view the dashboard.

![Grafana](https://github.com/Harry-027/go-notify/blob/master/snapshots/grafana.PNG "Grafana")

## Load Testing :-
---
* Before running load tests, replace the hostip variable value with your machine ip, under loadtesting/tests/loadtests.js
* Run the command `make load-testing` to run the load tests.

![LoadTesting_results](https://github.com/Harry-027/go-notify/blob/master/snapshots/loadTestingResults.PNG "LoadTesting_results")

## Sample for a custom HTML received mail :-
---

![Mail](https://github.com/Harry-027/go-notify/blob/master/snapshots/mailSample.PNG "Mail")

# Contributing :beers:
---
* Performance improvements, bug fixes, better design approaches are welcome. Please discuss any changes by raising an issue, beforehand.

# Maintainer :sunglasses:
---
[Harish Bhawnani Linkedin](https://www.linkedin.com/in/harish-bhawnani-86728457)
[Email](harishmmp@gmail.com)

## License
---
[MIT](LICENSE) © Harish Bhawnani
