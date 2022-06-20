package main

import (
	"os"
	"os/signal"

	"github.com/Harry-027/go-notify/api-server/database"
	"github.com/Harry-027/go-notify/cron/cron_jobs"
)

func main() {
	cron_jobs.InitKafkaConn()               // connect with Kafka
	database.ConnectDB()                    // connect with DB
	go cron_jobs.ScheduleDailyJob()         // Schedule the pending cron jobs
	go cron_jobs.ScheduleJobOnServerStart() // Schedule the active cron jobs
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
