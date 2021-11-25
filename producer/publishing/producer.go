package publishing

import (
	"context"

	"github.com/rick-and-morty-character-migration/producer/model"
)

type Publisher interface {
	Publish(ctx context.Context, message model.Message) error
}
