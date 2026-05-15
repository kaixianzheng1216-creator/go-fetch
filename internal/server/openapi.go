package server

import (
	"encoding/json"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/session"
)

func OpenAPIJSON() ([]byte, error) {
	r := chi.NewRouter()
	api := humachi.New(r, humaConfig())
	registerAPIRoutes(api, &App{})
	doc := api.OpenAPI()
	localizeOpenAPIErrorSchemas(doc)
	localizeOpenAPIErrorResponses(doc)
	return json.MarshalIndent(doc, "", "  ")
}

func humaConfig() huma.Config {
	cfg := huma.DefaultConfig("go-fetch Analytics API", "0.1.0")
	cfg.DocsPath = "/api/docs"
	cfg.SchemasPath = ""
	cfg.CreateHooks = nil
	cfg.Servers = []*huma.Server{{URL: "/"}}
	cfg.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"sessionCookie": {
			Type: "apiKey",
			In:   "cookie",
			Name: session.CookieName,
		},
	}
	return cfg
}

func localizeOpenAPIErrorResponses(doc *huma.OpenAPI) {
	for _, path := range doc.Paths {
		operations := []*huma.Operation{
			path.Get,
			path.Put,
			path.Post,
			path.Delete,
			path.Options,
			path.Head,
			path.Patch,
			path.Trace,
		}

		for _, op := range operations {
			if op == nil {
				continue
			}

			for status, response := range op.Responses {
				if response == nil {
					continue
				}

				description, ok := localizedOpenAPIResponseDescription(status)
				if ok {
					response.Description = description
				}
			}
		}
	}
}

func localizeOpenAPIErrorSchemas(doc *huma.OpenAPI) {
	if doc.Components == nil || doc.Components.Schemas == nil {
		return
	}

	schemas := doc.Components.Schemas.Map()
	if errorModel := schemas["ErrorModel"]; errorModel != nil {
		localizeSchemaProperty(errorModel, "type", "错误类型文档地址。")
		localizeSchemaProperty(errorModel, "title", "错误类型摘要。", "请求错误")
		localizeSchemaProperty(errorModel, "status", "HTTP 状态码。")
		localizeSchemaProperty(errorModel, "detail", "本次错误的具体说明。", "字段 foo 为必填项。")
		localizeSchemaProperty(errorModel, "instance", "本次错误实例的标识地址。")
		localizeSchemaProperty(errorModel, "errors", "字段级错误详情列表。")
	}

	if errorDetail := schemas["ErrorDetail"]; errorDetail != nil {
		localizeSchemaProperty(errorDetail, "location", "错误发生位置，例如 body.items[3].tags 或 path.thing-id。")
		localizeSchemaProperty(errorDetail, "message", "错误消息。")
		localizeSchemaProperty(errorDetail, "value", "触发错误的值。")
	}
}

func localizeSchemaProperty(schema *huma.Schema, propertyName, description string, examples ...any) {
	if schema == nil || schema.Properties == nil {
		return
	}

	property := schema.Properties[propertyName]
	if property == nil {
		return
	}

	property.Description = description
	if len(examples) > 0 {
		property.Examples = examples
	}
}

func localizedOpenAPIResponseDescription(status string) (string, bool) {
	if status == "default" {
		return "错误", true
	}

	code, err := strconv.Atoi(status)
	if err != nil || code < 400 {
		return "", false
	}

	return localizedStatusText(code), true
}
