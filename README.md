# go-notify
---
An email automation tool written in [Golang](https://golang.org/).
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
* Swagger included built using Swago[https://github.com/swaggo/swag]
* Full proof auth features - signup, login, update password, forgot password, logout.
* Cors, Helmet, Api-rate limiter included as middleware from security perspective.
* Users can register their clients & custom HTML templates.
* Mail scheduling (daily, weekly, monthly) using cron jobs.
* Subscription & Payment (payment has been stubbed for now and can be replaced with any suitable payment gateway).
* Api-server can be consumed by any client-side API, however for easy interaction - CLI (built using cobra) has been included.

## Installation & setup :-
---
* Go,Docker,Docker compose & Make should be pre-installed.

* Clone the repository.
```bash
git clone https://github.com/Harry-027/go-notify.git
```
* Run the command `make download`.
* Run the command `make setup`.
* CLI & Opentts engine setup is now ready to use.
* Run the command `go-audio --help` to explore various operations.

## Generate the audio files :-
---
* Run the command `go-audio aud --input=PATH_TO_PDF --output=PATH_TO_OUTPUT --voice=male|female`.
* Default value for given flags ::
    * input=./sample_pdf/test.pdf
    * output=homeDir/audio-go/output
    * voice=female
