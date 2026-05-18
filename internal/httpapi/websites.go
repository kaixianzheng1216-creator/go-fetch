package httpapi

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

type WebsiteRequest struct {
	Name   string `json:"name" required:"true" minLength:"1" maxLength:"100"`
	Domain string `json:"domain,omitempty" maxLength:"500"`
}

type createWebsiteInput struct {
	Body WebsiteRequest
}

type websiteIDInput struct {
	WebsiteID string `path:"websiteID" format:"uuid"`
}

type updateWebsiteInput struct {
	WebsiteID string `path:"websiteID" format:"uuid"`
	Body      WebsiteRequest
}

type WebsiteResponse struct {
	ID        uuid.UUID `json:"id" format:"uuid"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain"`
	CreatedAt time.Time `json:"createdAt"`
}

type websiteListOutput struct {
	Body []WebsiteResponse
}

type websiteOutput struct {
	Body WebsiteResponse
}

func (apiServer server) registerWebsiteRoutes(humaAPI huma.API, authMiddleware huma.Middlewares) {
	huma.Register(
		humaAPI,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/api/websites",
			OperationID: "listWebsites",
			Summary:     "列出站点",
			Tags:        []string{"Websites"},
			Security:    []map[string][]string{{"sessionCookie": {}}},
			Middlewares: authMiddleware,
		},
		apiServer.listWebsites,
	)

	createOperation := huma.Operation{
		Method:      http.MethodPost,
		Path:        "/api/websites",
		OperationID: "createWebsite",
		Summary:     "创建站点",
		Tags:        []string{"Websites"},
		Security:    []map[string][]string{{"sessionCookie": {}}},
		Middlewares: authMiddleware,
	}
	createOperation.DefaultStatus = http.StatusCreated
	huma.Register(humaAPI, createOperation, apiServer.createWebsite)

	huma.Register(
		humaAPI,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/api/websites/{websiteID}",
			OperationID: "getWebsite",
			Summary:     "获取站点",
			Tags:        []string{"Websites"},
			Security:    []map[string][]string{{"sessionCookie": {}}},
			Middlewares: authMiddleware,
		},
		apiServer.getWebsite,
	)

	huma.Register(
		humaAPI,
		huma.Operation{
			Method:      http.MethodPatch,
			Path:        "/api/websites/{websiteID}",
			OperationID: "updateWebsite",
			Summary:     "更新站点",
			Tags:        []string{"Websites"},
			Security:    []map[string][]string{{"sessionCookie": {}}},
			Middlewares: authMiddleware,
		},
		apiServer.updateWebsite,
	)

	huma.Register(
		humaAPI,
		huma.Operation{
			Method:      http.MethodDelete,
			Path:        "/api/websites/{websiteID}",
			OperationID: "deleteWebsite",
			Summary:     "删除站点",
			Tags:        []string{"Websites"},
			Security:    []map[string][]string{{"sessionCookie": {}}},
			Middlewares: authMiddleware,
		},
		apiServer.deleteWebsite,
	)
}

func (apiServer server) listWebsites(ctx context.Context, _ *emptyInput) (*websiteListOutput, error) {
	websites, err := apiServer.websites.List(ctx, currentUser(ctx).ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("加载站点列表失败")
	}

	return &websiteListOutput{Body: toWebsiteResponses(websites)}, nil
}

func (apiServer server) createWebsite(ctx context.Context, input *createWebsiteInput) (*websiteOutput, error) {
	website, err := apiServer.websites.Create(ctx, currentUser(ctx).ID, input.Body.Name, input.Body.Domain)
	if err != nil {
		if errors.Is(err, service.ErrInvalidWebsiteName) {
			return nil, huma.Error400BadRequest("站点名称不能为空")
		}
		return nil, huma.Error500InternalServerError("创建站点失败")
	}

	return &websiteOutput{Body: toWebsiteResponse(website)}, nil
}

func (apiServer server) getWebsite(ctx context.Context, input *websiteIDInput) (*websiteOutput, error) {
	websiteID, err := parseUUID(input.WebsiteID, "websiteID")
	if err != nil {
		return nil, err
	}

	website, err := apiServer.websites.Get(ctx, currentUser(ctx).ID, websiteID)
	if err != nil {
		return nil, websiteLookupError(err)
	}

	return &websiteOutput{Body: toWebsiteResponse(website)}, nil
}

func (apiServer server) updateWebsite(ctx context.Context, input *updateWebsiteInput) (*websiteOutput, error) {
	websiteID, err := parseUUID(input.WebsiteID, "websiteID")
	if err != nil {
		return nil, err
	}

	website, err := apiServer.websites.Update(ctx, currentUser(ctx).ID, websiteID, input.Body.Name, input.Body.Domain)
	if err != nil {
		if errors.Is(err, service.ErrInvalidWebsiteName) {
			return nil, huma.Error400BadRequest("站点名称不能为空")
		}
		return nil, websiteLookupError(err)
	}

	return &websiteOutput{Body: toWebsiteResponse(website)}, nil
}

func (apiServer server) deleteWebsite(ctx context.Context, input *websiteIDInput) (*okOutput, error) {
	websiteID, err := parseUUID(input.WebsiteID, "websiteID")
	if err != nil {
		return nil, err
	}

	if err := apiServer.websites.Delete(ctx, currentUser(ctx).ID, websiteID); err != nil {
		return nil, websiteLookupError(err)
	}

	return toOKOutput(), nil
}

func toWebsiteResponse(website domain.Website) WebsiteResponse {
	return WebsiteResponse{
		ID:        website.ID,
		Name:      website.Name,
		Domain:    website.Domain,
		CreatedAt: website.CreatedAt,
	}
}

func toWebsiteResponses(websites []domain.Website) []WebsiteResponse {
	result := make([]WebsiteResponse, 0, len(websites))
	for _, website := range websites {
		result = append(result, toWebsiteResponse(website))
	}
	return result
}
