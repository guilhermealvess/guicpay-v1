package event

import (
	"context"
)

type EventNotification interface {
	PublishEntity(ctx context.Context, entity any) error
}
