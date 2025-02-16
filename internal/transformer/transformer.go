package transformer

import (
	"context"
	"github.com/patrykorwat/infraconv/internal/parser"
)

type Transformer interface {
	Transform(ctx context.Context, config *parser.Config, directoryOutput string) error
}
