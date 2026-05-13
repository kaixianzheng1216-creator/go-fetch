package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi"

	"github.com/danielgtaylor/huma/v2"
)

func registerWebsiteRoutes(api huma.API, app *App, auth huma.Middlewares) {
	huma.Register(api, authenticated(operation(http.MethodGet, "/api/websites", "listWebsites", "Websites", http.StatusUnauthorized, http.StatusInternalServerError), auth), app.listWebsites)

	createOp := authenticated(operation(http.MethodPost, "/api/websites", "createWebsite", "Websites", http.StatusBadRequest, http.StatusUnauthorized, http.StatusInternalServerError), auth)
	createOp.DefaultStatus = http.StatusCreated
	createOp.SkipValidateBody = true
	huma.Register(api, createOp, app.createWebsite)

	huma.Register(api, authenticated(operation(http.MethodGet, "/api/websites/{websiteID}", "getWebsite", "Websites", http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError), auth), app.getWebsite)

	updateOp := authenticated(operation(http.MethodPatch, "/api/websites/{websiteID}", "updateWebsite", "Websites", http.StatusBadRequest, http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError), auth)
	updateOp.SkipValidateBody = true
	huma.Register(api, updateOp, app.updateWebsite)

	huma.Register(api, authenticated(operation(http.MethodDelete, "/api/websites/{websiteID}", "deleteWebsite", "Websites", http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError), auth), app.deleteWebsite)
}

func (a *App) listWebsites(ctx context.Context, _ *emptyInput) (*jsonBody[[]httpapi.Website], error) {
	websites, err := a.store.ListWebsites(ctx, userFromContext(ctx).ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to list websites")
	}

	return &jsonBody[[]httpapi.Website]{Body: httpapi.WebsitesFromDomain(websites)}, nil
}

func (a *App) createWebsite(ctx context.Context, input *websiteBodyInput) (*jsonBody[httpapi.Website], error) {
	request := normalizeWebsiteRequest(input.Body)
	if request.Name == "" {
		return nil, huma.Error400BadRequest("name is required")
	}

	website, err := a.store.CreateWebsite(ctx, userFromContext(ctx).ID, request.Name, request.Domain)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to create website")
	}

	return &jsonBody[httpapi.Website]{Body: httpapi.WebsiteFromDomain(website)}, nil
}

func (a *App) getWebsite(ctx context.Context, input *websitePathInput) (*jsonBody[httpapi.Website], error) {
	website, err := a.store.GetWebsite(ctx, userFromContext(ctx).ID, input.WebsiteID)
	if err != nil {
		return nil, websiteLookupError(err)
	}

	return &jsonBody[httpapi.Website]{Body: httpapi.WebsiteFromDomain(website)}, nil
}

func (a *App) updateWebsite(ctx context.Context, input *updateWebsiteInput) (*jsonBody[httpapi.Website], error) {
	request := normalizeWebsiteRequest(input.Body)
	if request.Name == "" {
		return nil, huma.Error400BadRequest("name is required")
	}

	user := userFromContext(ctx)
	if err := a.store.UpdateWebsite(ctx, user.ID, input.WebsiteID, request.Name, request.Domain); err != nil {
		return nil, websiteLookupError(err)
	}

	website, err := a.store.GetWebsite(ctx, user.ID, input.WebsiteID)
	if err != nil {
		return nil, websiteLookupError(err)
	}

	return &jsonBody[httpapi.Website]{Body: httpapi.WebsiteFromDomain(website)}, nil
}

func (a *App) deleteWebsite(ctx context.Context, input *websitePathInput) (*jsonBody[httpapi.OK], error) {
	if err := a.store.DeleteWebsite(ctx, userFromContext(ctx).ID, input.WebsiteID); err != nil {
		return nil, websiteLookupError(err)
	}

	return &jsonBody[httpapi.OK]{Body: httpapi.OK{OK: true}}, nil
}

func normalizeWebsiteRequest(request httpapi.WebsiteRequest) httpapi.WebsiteRequest {
	request.Name = strings.TrimSpace(request.Name)
	request.Domain = strings.TrimSpace(request.Domain)
	return request
}
