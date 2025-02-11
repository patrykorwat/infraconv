package transformer

import (
	"context"
	"github.com/pkg/errors"

	"github.com/patrykorwat/infraconv/internal/parser"
	"github.com/rs/zerolog/log"
	"github.com/upbound/provider-aws/config"
	"github.com/upbound/provider-aws/config/ec2"
)

type crossplaneTransformer struct {
}

func (c crossplaneTransformer) Transform(cfg *parser.Config, directoryOutput string) error {
	//p := &config.Provider{
	//	Resources: make(map[string]*config.Resource),
	//}
	ctx := context.Background()
	awsProvider, err := config.GetProvider(ctx, false)
	if err != nil {
		return errors.Wrap(err, "parsing error")
	}

	ec2.Configure(awsProvider) // awsProvider size is substantial

	// GOLANG 1.24 - Feat 3: SwissTable: 958 elements, almost will switch to SwissTable https://abseil.io/about/design/swisstables
	log.Info().Int("len", len(awsProvider.Resources)).Msg("Found resources")

	return nil
}

func NewCrossplaneTransformer() Transformer {
	return &crossplaneTransformer{}
}
