package events

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type EmailEvent struct {
	Event     string `json:"event"`
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Link      string `json:"link,omitempty"`
	ResetLink string `json:"reset_link,omitempty"`
}

type NatsClient struct {
	conn *nats.Conn
}

func NewNatsClient(url string) (*NatsClient, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsClient{conn: nc}, nil
}

func (nc *NatsClient) PublishUserRegistered(userID uint, email, name string) error {
	event := EmailEvent{
		Event:  "user.registered",
		UserID: userID,
		Email:  email,
		Name:   name,
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("❌ Error marshaling event: %v", err)
		return err
	}

	err = nc.conn.Publish("user.registered", data)
	if err != nil {
		log.Printf("❌ Error publishing to NATS: %v", err)
		return err
	}

	log.Printf("✅ Published registration event for user: %s", email)
	return nil
}
