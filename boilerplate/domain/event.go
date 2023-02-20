package domain

import (
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/johnnrails/ddd_go/boilerplate/identity"
)

type RawEvent interface {
	GetType() string
}

type Event struct {
	ID         uuid.UUID      `json:"id"`
	Type       string         `json:"type"`
	StreamID   uuid.UUID      `json:"stream_id"`
	StreamName string         `json:"stream_name"`
	OccurredAt time.Time      `json:"occurred_at"`
	ExpiresAt  *time.Time     `json:"expires_at,omitempty"`
	Payload    interface{}    `json:"payload,omitempty"`
	Metadata   *EventMetadata `json:"metadata,omitempty"`
}

func NewEventFromRawEvent(streamID uuid.UUID, streamName string, rawEvent RawEvent) (*Event, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("could not generate event id: %w", err)
	}

	return &Event{
		ID:         id,
		Type:       rawEvent.GetType(),
		StreamID:   streamID,
		StreamName: streamName,
		OccurredAt: time.Now(),
		Payload:    rawEvent,
	}, nil
}

type EventMetadata struct {
	Identity  *identity.Identity `json:"identity,omitempty"`
	IPAddress net.IP             `json:"ip_address,omitempty"`
	UserAgent string             `json:"http_user_agent,omitempty"`
	Referer   string             `json:"http_referer,omitempty"`
}

func (m *EventMetadata) IsEmpty() bool {
	return m.IPAddress == nil && m.Identity == nil && m.UserAgent == "" && m.Referer == ""
}
