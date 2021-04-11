package generate

import (
	"context"
	"github.com/wojnosystems/flick/pkg/generate/dsl"
	"io"
)

type Maker interface {
	Generate(ctx context.Context, definition *dsl.Document, output io.Writer) error
}
