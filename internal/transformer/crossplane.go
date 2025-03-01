package transformer

import (
	"context"
	"fmt"
	upjetConfig "github.com/crossplane/upjet/pkg/config"
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

func (c *crossplaneTransformer) Transform(ctx context.Context, cfg *parser.Config, directoryOutput string) error {
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

	log.Info().Int("len", len(awsProvider.Resources)).Msg("Found resources")

	newFile, err := os.Create("converted-resources.yaml")
	defer newFile.Close()
	if err != nil {
		return errors.Wrap(err, "Cannot create output file")
	}
	yamlPrinter := &printers.YAMLPrinter{}

	region := cfg.Providers[0].Attributes["region"].(string)

	for _, resource := range cfg.Resources {
		err = c.convertResource(awsProvider, resource, awsScheme, region, yamlPrinter, newFile)
		if err != nil {
			log.Warn().Err(err).Str("resourceName", resource.Name).Str("resourceType", resource.Type).
				Msg("Cannot convert resource")
		}
	}

	return nil
}

func (c *crossplaneTransformer) convertResource(awsProvider *upjetConfig.Provider, resource *parser.Resource,
	awsScheme *runtime.Scheme, region string, yamlPrinter *printers.YAMLPrinter, newFile *os.File) error {
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
		newType := reflect.New(forProviderField.Type)
		parametersInstance := newType.Interface()
		regionField := newType.Elem().FieldByName("Region")
		if regionField.IsValid() && regionField.CanSet() {
			regionField.Set(reflect.ValueOf(&region))
		}

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
		c.setResourceName(resource, convertedResource)

		err = yamlPrinter.PrintObj(convertedResource, newFile)
		if err != nil {
			return errors.Wrap(err, "Cannot print converted resource")
		}
	} else {
		return errors.New("Couldn't find type: " + gvk.String())
	}
	return nil
}

func (c *crossplaneTransformer) setResourceName(resource *parser.Resource, convertedResource *unstructured.Unstructured) {
	c.setDefaultResourceName(convertedResource, resource)

	// in some cases, name needs to be overridden as it is used as target resource name as for S3 Buckets
	if cfg, ok := config.TerraformPluginSDKExternalNameConfigs[resource.Type]; ok {
		if len(cfg.IdentifierFields) == 1 {
			if _, ok := resource.Attributes[cfg.IdentifierFields[0]]; ok {
				convertedResource.SetName(resource.Attributes[cfg.IdentifierFields[0]].(string))
			}
		}
	}
}

func (c *crossplaneTransformer) setDefaultResourceName(convertedResource *unstructured.Unstructured, resource *parser.Resource) {
	name := resource.Name
	convertedResource.SetName(name)
	convertedResource.SetAnnotations(map[string]string{
		"infraconv-name": name,
	})
}

func NewCrossplaneTransformer() Transformer {
	return &crossplaneTransformer{}
}
