package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type Notification struct {
	UserID           string `json:"userId"`
	Title            string `json:"title"`
	Message          string `json:"message"`
	NotificationType string `json:"notificationType"`
}

func main() {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()
	ch.QueueDeclare("notification_queue", true, false, false, false, nil)

	r := gin.Default()
	r.POST("/notifications", func(c *gin.Context) {
		var notif Notification
		if err := c.ShouldBindJSON(&notif); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		body, _ := json.Marshal(notif)
		ch.Publish("", "notification_queue", false, false, amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
		c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
	})
	r.Run(":8080")
}
