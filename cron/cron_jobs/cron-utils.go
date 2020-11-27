package cron_jobs

import (
	"context"
	"github.com/Harry-027/go-notify/api-server/config"
	"github.com/Harry-027/go-notify/api-server/models"
	"github.com/robfig/cron/v3"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

var KafkaConn *kafka.Conn

func InitKafkaConn() {
	config.LoadConfig()
	protocol := config.GetConfig(config.KAFKA_PROTOCOL)
	broker := config.GetConfig(config.BROKER)
	topic := config.GetConfig(config.TOPIC)
	conn, err := kafka.DialLeader(context.Background(), protocol, broker, topic, 0)
	if err != nil {
		log.Fatal("Failed to connect with kafka", err.Error())
	}
	_ = conn.SetWriteDeadline(time.Time{})
	log.Println("Connected with Kafka server successfully !!")
	KafkaConn = conn
}

// CronScheduler ...
type CronScheduler struct {
	c *cron.Cron
}

// cron constructor ...
func GetNewCron() CronScheduler {
	return CronScheduler{
		c: cron.New(),
	}
}

// Ass the cron jobs before start ...
func SetCronJob(cronSchedulerType CronScheduler, job models.Job, fn func(job2 models.Job) error) {
	log.Println("Setting jobs for various client preferences ...")
	entryId, err := cronSchedulerType.c.AddFunc(job.Type, func() {
		err := fn(job)
		if err != nil {
			log.Println("An error occurred while callback: ", err)
		}
	})
	if err != nil {
		log.Println("An error occurred while setting cron job: ", err)
	}
	log.Println("EntryId for cron job: ", entryId)
}
