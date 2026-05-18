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

type createWebsiteInput struct {
	Body struct {
		Name   string `json:"name" required:"true" minLength:"1" maxLength:"100"`
		Domain string `json:"domain,omitempty" maxLength:"500"`
	}
}

type websiteIDInput struct {
	WebsiteID uuid.UUID `path:"websiteID" format:"uuid"`
}

type updateWebsiteInput struct {
	WebsiteID uuid.UUID `path:"websiteID" format:"uuid"`
	Body      struct {
		Name   string `json:"name" required:"true" minLength:"1" maxLength:"100"`
		Domain string `json:"domain,omitempty" maxLength:"500"`
	}
}

type websiteListOutput struct {
	Body websiteListBody
}

type websiteOutput struct {
	Body struct {
		ID        uuid.UUID `json:"id" format:"uuid"`
		Name      string    `json:"name"`
		Domain    string    `json:"domain"`
		CreatedAt time.Time `json:"createdAt"`
	}
}

type websiteListBody []struct {
	ID        uuid.UUID `json:"id" format:"uuid"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain"`
	CreatedAt time.Time `json:"createdAt"`
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

	return newWebsiteListOutput(websites), nil
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

	return newWebsiteOutput(website), nil
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

	return newWebsiteOutput(website), nil
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

	return newWebsiteOutput(website), nil
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

func newWebsiteListOutput(websites []domain.Website) *websiteListOutput {
	output := &websiteListOutput{Body: make(websiteListBody, len(websites))}
	for i, website := range websites {
		output.Body[i].ID = website.ID
		output.Body[i].Name = website.Name
		output.Body[i].Domain = website.Domain
		output.Body[i].CreatedAt = website.CreatedAt
	}
	return output
}

func newWebsiteOutput(website domain.Website) *websiteOutput {
	output := &websiteOutput{}
	output.Body.ID = website.ID
	output.Body.Name = website.Name
	output.Body.Domain = website.Domain
	output.Body.CreatedAt = website.CreatedAt
	return output
}

func websiteLookupError(err error) error {
	if err == nil {
		return nil
	}
	if isNotFound(err) {
		return huma.Error404NotFound(errorMessageWebsiteNotFound)
	}
	return huma.Error500InternalServerError(errorMessageWebsiteLoadFailed)
}

func websiteMutationError(err error, fallbackMessage string) error {
	if errors.Is(err, service.ErrInvalidWebsiteName) {
		return huma.Error400BadRequest(errorMessageWebsiteNameCannotEmpty)
	}
	return websiteLookupErrorWithFallback(err, fallbackMessage)
}

func websiteLookupErrorWithFallback(err error, fallbackMessage string) error {
	if err == nil {
		return nil
	}
	if isNotFound(err) {
		return huma.Error404NotFound(errorMessageWebsiteNotFound)
	}
	return huma.Error500InternalServerError(fallbackMessage)
}
