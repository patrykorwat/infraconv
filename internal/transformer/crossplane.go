package transformer

import (
	"context"
	"fmt"
	"github.com/crossplane/upjet/pkg/resource/json"
	"github.com/patrykorwat/infraconv/internal/parser"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/upbound/provider-aws/apis"
	"github.com/upbound/provider-aws/config"
	"github.com/upbound/provider-aws/config/ec2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"
	"os"
	"reflect"
)

type crossplaneTransformer struct {
}

func (c crossplaneTransformer) Transform(ctx context.Context, cfg *parser.Config, directoryOutput string) error {
	awsProvider, err := config.GetProvider(ctx, false)
	if err != nil {
		return errors.Wrap(err, "parsing error")
	}

	ec2.Configure(awsProvider) // awsProvider size is substantial

	awsScheme := runtime.NewScheme()
	err = apis.AddToScheme(awsScheme)
	if err != nil {
		return errors.Wrap(err, "cannot add schemes")
	}

	// GOLANG 1.24 - Feat 3: SwissTable: 958 elements, almost will switch to SwissTable https://abseil.io/about/design/swisstables
	log.Info().Int("len", len(awsProvider.Resources)).Msg("Found resources")

	for _, resource := range cfg.Resources {
		if _, ok := awsProvider.Resources[resource.Type]; !ok {
			return errors.New("Couldn't find resource type: " + resource.Type)
		}
		cRes := awsProvider.Resources[resource.Type]
		gvk := schema.GroupVersionKind{
			Group:   fmt.Sprintf("%s.%s", cRes.ShortGroup, awsProvider.RootGroup),
			Version: cRes.Version,
			Kind:    cRes.Kind,
		}
		if knownType, ok := awsScheme.AllKnownTypes()[gvk]; ok {
			specField, _ := knownType.FieldByName("Spec")
			forProviderField, _ := specField.Type.FieldByName("ForProvider")
			parametersInstance := reflect.New(forProviderField.Type).Interface()

			marshal, err := json.TFParser.Marshal(resource.Attributes)
			if err != nil {
				return errors.Wrap(err, "cannot serialize TF json")
			}
			log.Info().Any("marshalled data", string(marshal)).Msg("pre-transformed resource")

			err = json.TFParser.Unmarshal(marshal, &parametersInstance)
			if err != nil {
				return errors.Wrap(err, "cannot convert TF json")
			}
			log.Info().Any("transformedResource", parametersInstance).Msg("transformed resource")

			transformedResourceBytes, _ := json.JSParser.Marshal(parametersInstance)
			log.Info().Any("transformedResource", string(transformedResourceBytes)).Msg("transformed resource")

			convertedResource := &unstructured.Unstructured{}
			convertedResource.SetUnstructuredContent(map[string]interface{}{
				"spec": map[string]interface{}{
					"forProvider": parametersInstance,
				},
			})
			convertedResource.SetGroupVersionKind(gvk)

			newFile, err := os.Create("converted-resources.yaml")
			if err != nil {
				return errors.Wrap(err, "Cannot create output file")
			}
			y := printers.YAMLPrinter{}
			defer newFile.Close()
			err = y.PrintObj(convertedResource, newFile)
			if err != nil {
				return errors.Wrap(err, "Cannot print converted resource")
			}
		} else {
			return errors.New("Couldn't find type: " + gvk.String())
		}
	}

	return nil
}

func NewCrossplaneTransformer() Transformer {
	return &crossplaneTransformer{}
}
