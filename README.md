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

![go-notify](https://github.com/Harry-027/go-notify/snapshots/system_diagram.png "go-notify")

## DB Schema:
---

![DB_Schema](https://github.com/Harry-027/go-notify/snapshots/dbSchema.PNG "go-notify")

## Installation & setup :-
---
* Go,Docker & Make should be pre-installed.
* Clone the repository.
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
