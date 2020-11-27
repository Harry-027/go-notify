package utils

import (
	"context"
	"github.com/Harry-027/go-notify/api-server/config"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"time"
)

// Utility to Hash the password ...
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Utility to compare the password and hash ...
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Connect with Kafka ...
func InitKafkaConn() *kafka.Conn {
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
	return conn
}

// Get the cron job param based on preference type ...
func GetType(pref string) string {
	pref = strings.ToLower(pref)
	cronType := map[string]string{
		"daily":   "@daily",
		"weekly":  "@weekly",
		"monthly": "@monthly",
	}
	return cronType[pref]
}

func ApiResponse(ctx *fiber.Ctx, statusCode int, apiConst config.ApiResponse) error {
	return ctx.Status(statusCode).JSON(fiber.Map{"status": apiConst.Status, "message": apiConst.Message})
}

func ApiResponseWithCustomMsg(ctx *fiber.Ctx, statusCode int, resp interface{}) error {
	return ctx.Status(statusCode).JSON(resp)
}
