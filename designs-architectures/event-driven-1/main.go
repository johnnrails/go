package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
	pubsubproto "github.com/vardius/pubsub/v2/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Event contains id, payload and metadata
type Event struct {
	ID            uuid.UUID   `json:"id"`
	Type          string      `json:"type"`
	StreamID      uuid.UUID   `json:"stream_id"`
	StreamName    string      `json:"stream_name"`
	StreamVersion int         `json:"stream_version"`
	OccurredAt    time.Time   `json:"occurred_at"`
	ExpiresAt     *time.Time  `json:"expires_at,omitempty"`
	Payload       interface{} `json:"payload,omitempty"`
}

type EventHandler func(ctx context.Context, event *Event) error

type EventBus interface {
	Publish(ctx context.Context, event *Event) error
	Subscribe(ctx context.Context, eventType string, fn EventHandler) error
	Unsubscribe(ctx context.Context, eventType string, fn EventHandler) error
	PublishAndAcknowledge(parentCtx context.Context, event *Event) error
}

type eventBus struct {
	handlerTimeout      time.Duration
	pubsub              pubsubproto.PubSubClient
	mtx                 sync.RWMutex
	unsubscribeChannels map[reflect.Value]chan struct{}
}

func NewEventBus(timeout time.Duration, pubsub pubsubproto.PubSubClient) EventBus {
	return &eventBus{
		handlerTimeout:      timeout,
		pubsub:              pubsub,
		unsubscribeChannels: make(map[reflect.Value]chan struct{}),
	}
}

// Subscrive registers handler to be notified of every event published
func (b *eventBus) Subscribe(ctx context.Context, eventType string, fn EventHandler) error {
	stream, err := b.pubsub.Subscribe(ctx, &pubsubproto.SubscribeRequest{
		Topic: eventType,
	})
	if err != nil {
		return err
	}

	unsubscribeCh := make(chan struct{}, 1)

	b.mtx.Lock()
	b.unsubscribeChannels[reflect.ValueOf(fn)] = unsubscribeCh
	b.mtx.Unlock()

	ctxDoneCh := ctx.Done()
	for {
		select {
		case <-ctxDoneCh:
			return ctx.Err()
		case <-unsubscribeCh:
			return nil
		default:
			resp, err := stream.Recv()
			if err != nil {
				return err
			}

			if err := b.dispatchEvent(resp.GetPayload(), fn); err != nil {
				return err
			}
		}
	}
}

func (b *eventBus) Publish(ctx context.Context, event *Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if _, err := b.pubsub.Publish(ctx, &pubsubproto.PublishRequest{
		Topic:   event.Type,
		Payload: payload,
	}); err != nil {
		return err
	}
	return nil
}

func (b *eventBus) Unsubscribe(ctx context.Context, eventType string, fn EventHandler) error {
	rv := reflect.ValueOf(fn)
	b.mtx.RLock()
	if ch, ok := b.unsubscribeChannels[rv]; ok {
		ch <- struct{}{}
	}
	b.mtx.RUnlock()
	return nil
}

func (b *eventBus) dispatchEvent(payload []byte, fn EventHandler) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.handlerTimeout)
	defer cancel()

	var e Event
	if err := json.Unmarshal(payload, &e); err != nil {
		return err
	}

	return fn(ctx, &e)
}

func (b *eventBus) PublishAndAcknowledge(ctx context.Context, event *Event) error {
	panic("not implemented")
}

func main() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial("localhost:9879", opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	bus := NewEventBus(300*time.Millisecond, pubsubproto.NewPubSubClient(conn))
	ctx := context.Background()
	err = bus.Publish(ctx, &Event{
		ID:   uuid.New(),
		Type: "add",
	})
	if err != nil {
		fmt.Println("ERR: Publish")
	}

	err = bus.Subscribe(ctx, "add", func(ctx context.Context, event *Event) error {
		fmt.Println(event.Payload)
		fmt.Println(event.ID)
		return nil
	})
	if err != nil {
		fmt.Println("ERR: Subscribe", err)
	}
}
