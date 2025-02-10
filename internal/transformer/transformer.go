package transformer

import "github.com/patrykorwat/infraconv/internal/parser"

type Transformer interface {
	Transform(config *parser.Config, directoryOutput string) error
}
