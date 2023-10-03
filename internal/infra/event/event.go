package event

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/internal/domain/event"
)

type eventNotification struct {
	url string
}

func NewEventNotification(endpointUrl string) event.EventNotification {
	return &eventNotification{endpointUrl}
}

type Message struct {
	MessageID   string         `json:"message_id"`
	Timestamp   time.Time      `json:"timestmap"`
	Source      string         `json:"source"`
	ContentType string         `json:"content_type"`
	Data        string         `json:"data"`
	Attributes  map[string]any `json:"attributes"`
}

func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(*m)
}

func (e eventNotification) PublishEntity(ctx context.Context, entity any) error {
	raw, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	message := Message{
		MessageID:   uuid.NewString(),
		Timestamp:   time.Now().UTC(),
		Source:      "guic-pay",
		ContentType: "application/json",
		Data:        string(raw),
	}

	messageRaw, err := message.Marshal()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, e.url, bytes.NewReader(messageRaw))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return errors.New("TODO")
	}

	return nil
}
