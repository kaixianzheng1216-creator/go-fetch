package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

func registerWebsiteRoutes(api huma.API, app *App, auth huma.Middlewares) {
	listOp := newOperation(
		http.MethodGet,
		"/api/websites",
		"listWebsites",
		"Websites",
		http.StatusUnauthorized,
		http.StatusInternalServerError,
	)

	huma.Register(api, withAuth(listOp, auth), app.listWebsites)

	createOp := newOperation(
		http.MethodPost,
		"/api/websites",
		"createWebsite",
		"Websites",
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusUnprocessableEntity,
		http.StatusInternalServerError,
	)

	createOp = withAuth(createOp, auth)
	createOp.DefaultStatus = http.StatusCreated

	huma.Register(api, createOp, app.createWebsite)

	getOp := newOperation(
		http.MethodGet,
		"/api/websites/{websiteID}",
		"getWebsite",
		"Websites",
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusInternalServerError,
	)

	huma.Register(api, withAuth(getOp, auth), app.getWebsite)

	updateOp := newOperation(
		http.MethodPatch,
		"/api/websites/{websiteID}",
		"updateWebsite",
		"Websites",
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusUnprocessableEntity,
		http.StatusInternalServerError,
	)

	updateOp = withAuth(updateOp, auth)

	huma.Register(api, updateOp, app.updateWebsite)

	deleteOp := newOperation(
		http.MethodDelete,
		"/api/websites/{websiteID}",
		"deleteWebsite",
		"Websites",
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusInternalServerError,
	)

	huma.Register(api, withAuth(deleteOp, auth), app.deleteWebsite)
}

func (a *App) listWebsites(ctx context.Context, _ *emptyInput) (*jsonBody[[]Website], error) {
	websites, err := a.store.ListWebsites(ctx, userFromContext(ctx).ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("加载网站列表失败")
	}

	response := toWebsites(websites)

	return jsonResponse(response), nil
}

func (a *App) createWebsite(ctx context.Context, input *websiteInput) (*jsonBody[Website], error) {
	request := normalizeWebsiteRequest(input.Body)
	if request.Name == "" {
		return nil, huma.Error400BadRequest("名称不能为空")
	}

	website, err := a.store.CreateWebsite(ctx, userFromContext(ctx).ID, request.Name, request.Domain)
	if err != nil {
		return nil, huma.Error500InternalServerError("创建网站失败")
	}

	response := toWebsite(website)

	return jsonResponse(response), nil
}

func (a *App) getWebsite(ctx context.Context, input *websiteIDInput) (*jsonBody[Website], error) {
	website, err := a.store.GetWebsite(ctx, userFromContext(ctx).ID, input.WebsiteID)
	if err != nil {
		return nil, websiteLookupError(err)
	}

	response := toWebsite(website)

	return jsonResponse(response), nil
}

func (a *App) updateWebsite(ctx context.Context, input *updateWebsiteInput) (*jsonBody[Website], error) {
	request := normalizeWebsiteRequest(input.Body)
	if request.Name == "" {
		return nil, huma.Error400BadRequest("名称不能为空")
	}

	user := userFromContext(ctx)
	if err := a.store.UpdateWebsite(ctx, user.ID, input.WebsiteID, request.Name, request.Domain); err != nil {
		return nil, websiteLookupError(err)
	}

	website, err := a.store.GetWebsite(ctx, user.ID, input.WebsiteID)
	if err != nil {
		return nil, websiteLookupError(err)
	}

	response := toWebsite(website)

	return jsonResponse(response), nil
}

func (a *App) deleteWebsite(ctx context.Context, input *websiteIDInput) (*jsonBody[OK], error) {
	if err := a.store.DeleteWebsite(ctx, userFromContext(ctx).ID, input.WebsiteID); err != nil {
		return nil, websiteLookupError(err)
	}

	return okResponse(), nil
}

func normalizeWebsiteRequest(request WebsiteRequest) WebsiteRequest {
	request.Name = strings.TrimSpace(request.Name)
	request.Domain = strings.TrimSpace(request.Domain)
	return request
}
