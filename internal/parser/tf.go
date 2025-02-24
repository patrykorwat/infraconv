package parser

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/rs/zerolog/log"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	"os"
	"path/filepath"
	"strings"
)

type tfParser struct {
}

func parseConfig(file *hcl.File) *Config {
	requiredProviders := make([]*RequiredProvider, 0)
	providers := make([]*Provider, 0)
	resources := make([]*Resource, 0)
	modules := make([]*Module, 0)
	locals := make([]*Local, 0)

	for _, block := range file.Body.(*hclsyntax.Body).Blocks {
		switch block.Type {
		case "terraform":
			requiredProviders = parseTerraform(block)
		case "provider":
			providers = append(providers, parseProvider(block))
		case "module":
			modules = append(modules, parseModule(block))
		case "resource":
			resources = append(resources, parseResource(block))
		case "locals":
			locals = append(locals, parseLocals(block))
		}
	}

	return &Config{RequiredProviders: requiredProviders, Providers: providers, Resources: resources, Modules: modules, Locals: locals}
}

func parseHCLFile(file string, parser *hclparse.Parser) (*Config, error) {
	if filepath.Ext(file) == ".tf" {
		_, err := os.Stat(file)
		if !os.IsNotExist(err) {
			file, diags := parser.ParseHCLFile(file)
			if diags.HasErrors() {
				return nil, fmt.Errorf("failed to load config file %s: %s", file, diags.Errs())
			}

			return parseConfig(file), nil
		}
	}

	return &Config{}, nil
}

func parseSingleFile(file string, hclParser *hclparse.Parser, config *Config) error {
	parsedConfig, err := parseHCLFile(file, hclParser)
	if err != nil {
		return fmt.Errorf("error parsing HCL file: %w", err)
	}

	config.RequiredProviders = append(config.RequiredProviders, parsedConfig.RequiredProviders...)
	config.Providers = append(config.Providers, parsedConfig.Providers...)
	config.Modules = append(config.Modules, parsedConfig.Modules...)
	config.Resources = append(config.Resources, parsedConfig.Resources...)
	config.Locals = append(config.Locals, parsedConfig.Locals...)

	return nil
}

func (t tfParser) Parse(directory string) (*Config, error) {
	config := &Config{}
	hclParser := hclparse.NewParser()

	err := filepath.Walk(directory, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory: %w", err)
		}

		if info.IsDir() || strings.Contains(file, ".terraform/") || filepath.Ext(file) != ".tf" {
			return nil
		}

		return parseSingleFile(file, hclParser, config)
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing a directory: %w", err)
	}

	return config, nil
}

func NewTfParser() Parser {
	return &tfParser{}
}

func parseTerraform(block *hclsyntax.Block) []*RequiredProvider {
	requiredProviders := make([]*RequiredProvider, 0)

	for _, block := range block.Body.Blocks {
		switch block.Type {
		case "required_providers":
			for _, attribute := range block.Body.Attributes {
				rp := RequiredProvider{
					Name:       attribute.Name,
					Attributes: evaluateExpression(attribute.Expr).(map[string]any),
				}
				requiredProviders = append(requiredProviders, &rp)
			}
		}
	}

	return requiredProviders
}

func parseProvider(block *hclsyntax.Block) *Provider {
	provider := &Provider{
		Name:       block.Labels[0],
		Attributes: map[string]any{},
	}

	for _, attribute := range block.Body.Attributes {
		value := evaluateExpression(attribute.Expr)
		provider.Attributes[attribute.Name] = value
	}

	return provider
}

func parseModule(block *hclsyntax.Block) *Module {
	module := &Module{
		Labels:     block.Labels,
		Attributes: map[string]any{},
	}

	for _, attribute := range block.Body.Attributes {
		value := evaluateExpression(attribute.Expr)
		module.Attributes[attribute.Name] = value

		if attribute.Name == "source" {
			module.Source = value.(string)
		}
	}

	return module
}

func parseResource(block *hclsyntax.Block) *Resource {
	resource := &Resource{
		Type:       block.Labels[0],
		Name:       block.Labels[1],
		Labels:     block.Labels,
		Attributes: map[string]any{},
	}

	for _, attribute := range block.Body.Attributes {
		value := evaluateExpression(attribute.Expr)
		resource.Attributes[attribute.Name] = value
	}

	for _, bodyBlock := range block.Body.Blocks {
		parseResourcesFromBlock(bodyBlock, resource)
	}

	return resource
}

func parseResourcesFromBlock(bodyBlock *hclsyntax.Block, resource *Resource) {
	blType := bodyBlock.Type
	if _, ok := resource.Attributes[blType]; !ok {
		resource.Attributes[blType] = map[string]any{}
	}

	data := resource.Attributes[blType].(map[string]any)

	for _, attribute := range bodyBlock.Body.Attributes {
		data[attribute.Name] = evaluateExpression(attribute.Expr)
	}
}

func parseLocals(block *hclsyntax.Block) *Local {
	local := &Local{
		Attributes: map[string]any{},
	}

	for _, attribute := range block.Body.Attributes {
		value := evaluateExpression(attribute.Expr)
		local.Attributes[attribute.Name] = value
	}

	return local
}

func buildVarExpressions(traversal hcl.Traversal) string {
	varExp := make([]string, 0, len(traversal))

	for _, v := range traversal {
		switch v := v.(type) {
		case hcl.TraverseRoot:
			if v.Name != "" {
				varExp = append(varExp, v.Name)
			}
		case hcl.TraverseAttr:
			if v.Name != "" {
				varExp = append(varExp, v.Name)
			}
		}
	}

	return strings.Join(varExp, ".")
}

func convertValue(val cty.Value) interface{} {
	switch val.Type() {
	case cty.Number:
		return val.AsBigFloat().String()
	case cty.String:
		return val.AsString()
	case cty.Bool:
		var v bool
		_ = gocty.FromCtyValue(val, &v)

		return v
	default:
		log.Warn().Str("type", val.Type().GoString()).Msg("unsupported type")
		return ""
	}
}

func evaluateExpression(expr hcl.Expression) any {
	resultString := ""
	resultMap := map[string]any{}

	switch expr := expr.(type) {
	case *hclsyntax.ScopeTraversalExpr:
		resultString += buildVarExpressions(expr.Traversal)
	case *hclsyntax.LiteralValueExpr:
		return convertValue(expr.Val)
	case *hclsyntax.TemplateExpr:
		parts := expr.Parts
		for _, part := range parts {
			resultString += evaluateExpression(part).(string)
		}
	case *hclsyntax.TupleConsExpr:
		for _, elem := range expr.Exprs {
			resultString += evaluateExpression(elem).(string) + ","
		}
	case *hclsyntax.ObjectConsKeyExpr:
		resultString += evaluateExpression(expr.Wrapped).(string)
	case *hclsyntax.ObjectConsExpr:
		for i := range expr.Items {
			item := expr.Items[i]

			resultMap[evaluateExpression(item.KeyExpr).(string)] = evaluateExpression(item.ValueExpr)
		}

		return resultMap
	case *hclsyntax.IndexExpr:
		resultString += evaluateExpression(expr.Collection).(string)
	case *hclsyntax.FunctionCallExpr:
		resultString += evaluateFunctionExpression(expr)
	default:
		log.Warn().Any("Expr", expr).Msg("unsupported expr")
	}

	return resultString
}

func evaluateFunctionExpression(expr *hclsyntax.FunctionCallExpr) string {
	var args string

	for i := range expr.Args {
		exp := evaluateExpression(expr.Args[i])

		// TODO: Implement other cases

		switch exp := exp.(type) {
		case string:
			args += exp
		case map[string]any:
			var values string

			for k, v := range exp {
				values += k

				switch v := v.(type) {
				case string:
					values += ":" + v
				default:
					log.Warn().Any("value", expr).Msg("unsupported function arg value")
				}
			}

			args = fmt.Sprintf("%s{%s}", args, values)
		default:
			log.Warn().Any("arg", expr).Msg("unsupported function arg value")
		}
	}

	return fmt.Sprintf("%s(%s)", expr.Name, args)
}
