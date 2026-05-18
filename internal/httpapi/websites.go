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
	Name   string `json:"name" required:"true" minLength:"1" maxLength:"100"`
	Domain string `json:"domain,omitempty" maxLength:"500"`
}

type WebsiteResponse struct {
	ID        uuid.UUID `json:"id" format:"uuid"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain"`
	CreatedAt time.Time `json:"createdAt"`
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

type websiteListOutput struct {
	Body []WebsiteResponse
}

type websiteOutput struct {
	Body WebsiteResponse
}

func (srv server) registerWebsiteRoutes(humaAPI huma.API, authMiddleware huma.Middlewares) {
	huma.Register(
		humaAPI,
		securedOperation(http.MethodGet, "/api/websites", "listWebsites", "List websites", "Websites", authMiddleware),
		srv.listWebsites,
	)

	createOperation := securedOperation(http.MethodPost, "/api/websites", "createWebsite", "Create website", "Websites", authMiddleware)
	createOperation.DefaultStatus = http.StatusCreated
	huma.Register(humaAPI, createOperation, srv.createWebsite)

	huma.Register(
		humaAPI,
		securedOperation(http.MethodGet, "/api/websites/{websiteID}", "getWebsite", "Get website", "Websites", authMiddleware),
		srv.getWebsite,
	)

	huma.Register(
		humaAPI,
		securedOperation(http.MethodPatch, "/api/websites/{websiteID}", "updateWebsite", "Update website", "Websites", authMiddleware),
		srv.updateWebsite,
	)

	huma.Register(
		humaAPI,
		securedOperation(http.MethodDelete, "/api/websites/{websiteID}", "deleteWebsite", "Delete website", "Websites", authMiddleware),
		srv.deleteWebsite,
	)
}

func (srv server) listWebsites(ctx context.Context, _ *emptyInput) (*websiteListOutput, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	websites, err := srv.websites.List(ctx, userID)
	if err != nil {
		return nil, huma.Error500InternalServerError(errorMessageWebsiteListLoadFailed)
	}

	return &websiteListOutput{Body: newWebsiteResponses(websites)}, nil
}

func (srv server) createWebsite(ctx context.Context, input *createWebsiteInput) (*websiteOutput, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	website, err := srv.websites.Create(ctx, userID, service.WebsiteInput{
		Name:   input.Body.Name,
		Domain: input.Body.Domain,
	})
	if err != nil {
		return nil, websiteMutationError(err, errorMessageWebsiteCreateFailed)
	}

	return &websiteOutput{Body: newWebsiteResponse(website)}, nil
}

func (srv server) getWebsite(ctx context.Context, input *websiteIDInput) (*websiteOutput, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	website, err := srv.websites.Get(ctx, userID, input.WebsiteID)
	if err != nil {
		return nil, websiteLookupError(err)
	}

	return &websiteOutput{Body: newWebsiteResponse(website)}, nil
}

func (srv server) updateWebsite(ctx context.Context, input *updateWebsiteInput) (*websiteOutput, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	website, err := srv.websites.Update(ctx, userID, input.WebsiteID, service.WebsiteInput{
		Name:   input.Body.Name,
		Domain: input.Body.Domain,
	})
	if err != nil {
		return nil, websiteMutationError(err, errorMessageWebsiteUpdateFailed)
	}

	return &websiteOutput{Body: newWebsiteResponse(website)}, nil
}

func (srv server) deleteWebsite(ctx context.Context, input *websiteIDInput) (*okOutput, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err := srv.websites.Delete(ctx, userID, input.WebsiteID); err != nil {
		return nil, websiteLookupError(err)
	}

	return newOKOutput(), nil
}

func newWebsiteResponse(website domain.Website) WebsiteResponse {
	return WebsiteResponse{
		ID:        website.ID,
		Name:      website.Name,
		Domain:    website.Domain,
		CreatedAt: website.CreatedAt,
	}
}

func newWebsiteResponses(websites []domain.Website) []WebsiteResponse {
	result := make([]WebsiteResponse, len(websites))
	for i, website := range websites {
		result[i] = newWebsiteResponse(website)
	}
	return result
}
