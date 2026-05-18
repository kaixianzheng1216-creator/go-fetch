package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

type WebsiteRequest struct {
	Name       string `json:"name" required:"true" minLength:"1" maxLength:"100"`
	DomainName string `json:"domain,omitempty" maxLength:"500"`
}

type createWebsiteInput struct {
	Body WebsiteRequest
}

type websiteIDInput struct {
	WebsiteID uuid.UUID `path:"websiteID" format:"uuid"`
}

type updateWebsiteInput struct {
	WebsiteID uuid.UUID `path:"websiteID" format:"uuid"`
	Body      WebsiteRequest
}

type WebsiteResponse struct {
	ID         uuid.UUID `json:"id" format:"uuid"`
	Name       string    `json:"name"`
	DomainName string    `json:"domain"`
	CreatedAt  time.Time `json:"createdAt"`
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
		securedOperation(http.MethodGet, "/api/websites", "listWebsites", "列出站点", "Websites", authMiddleware),
		apiServer.listWebsites,
	)

	createOperation := securedOperation(http.MethodPost, "/api/websites", "createWebsite", "创建站点", "Websites", authMiddleware)
	createOperation.DefaultStatus = http.StatusCreated
	huma.Register(humaAPI, createOperation, apiServer.createWebsite)

	huma.Register(
		humaAPI,
		securedOperation(http.MethodGet, "/api/websites/{websiteID}", "getWebsite", "获取站点", "Websites", authMiddleware),
		apiServer.getWebsite,
	)

	huma.Register(
		humaAPI,
		securedOperation(http.MethodPatch, "/api/websites/{websiteID}", "updateWebsite", "更新站点", "Websites", authMiddleware),
		apiServer.updateWebsite,
	)

	huma.Register(
		humaAPI,
		securedOperation(http.MethodDelete, "/api/websites/{websiteID}", "deleteWebsite", "删除站点", "Websites", authMiddleware),
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
	website, err := apiServer.websites.Create(ctx, currentUser(ctx).ID, service.WebsiteParams{
		Name:       input.Body.Name,
		DomainName: input.Body.DomainName,
	})
	if err != nil {
		return nil, websiteMutationError(err, "创建站点失败")
	}

	return &websiteOutput{Body: toWebsiteResponse(website)}, nil
}

func (apiServer server) getWebsite(ctx context.Context, input *websiteIDInput) (*websiteOutput, error) {
	website, err := apiServer.websites.Get(ctx, currentUser(ctx).ID, input.WebsiteID)
	if err != nil {
		return nil, websiteLookupError(err)
	}

	return &websiteOutput{Body: toWebsiteResponse(website)}, nil
}

func (apiServer server) updateWebsite(ctx context.Context, input *updateWebsiteInput) (*websiteOutput, error) {
	website, err := apiServer.websites.Update(ctx, currentUser(ctx).ID, input.WebsiteID, service.WebsiteParams{
		Name:       input.Body.Name,
		DomainName: input.Body.DomainName,
	})
	if err != nil {
		return nil, websiteMutationError(err, "更新站点失败")
	}

	return &websiteOutput{Body: toWebsiteResponse(website)}, nil
}

func (apiServer server) deleteWebsite(ctx context.Context, input *websiteIDInput) (*okOutput, error) {
	if err := apiServer.websites.Delete(ctx, currentUser(ctx).ID, input.WebsiteID); err != nil {
		return nil, websiteLookupError(err)
	}

	return toOKOutput(), nil
}

func toWebsiteResponse(website domain.Website) WebsiteResponse {
	return WebsiteResponse{
		ID:         website.ID,
		Name:       website.Name,
		DomainName: website.DomainName,
		CreatedAt:  website.CreatedAt,
	}
}

func toWebsiteResponses(websites []domain.Website) []WebsiteResponse {
	result := make([]WebsiteResponse, 0, len(websites))
	for _, website := range websites {
		result = append(result, toWebsiteResponse(website))
	}
	return result
}
