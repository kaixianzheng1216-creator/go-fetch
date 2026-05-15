package website

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	userdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/domain/user"
	websitedomain "github.com/kaixianzheng1216-creator/go-fetch/internal/domain/website"
)

type Store interface {
	ListWebsites(ctx context.Context, userID string) ([]websitedomain.Website, error)
	CreateWebsite(ctx context.Context, userID, name, domainName string) (websitedomain.Website, error)
	GetWebsite(ctx context.Context, userID, websiteID string) (websitedomain.Website, error)
	UpdateWebsite(ctx context.Context, userID, websiteID, name, domainName string) error
	DeleteWebsite(ctx context.Context, userID, websiteID string) error
}

type Handler struct {
	store              Store
	currentUser        func(context.Context) userdomain.User
	websiteLookupError func(error) error
}

func New(
	dataStore Store,
	currentUser func(context.Context) userdomain.User,
	websiteLookupError func(error) error,
) Handler {
	return Handler{
		store:              dataStore,
		currentUser:        currentUser,
		websiteLookupError: websiteLookupError,
	}
}

type websiteRequest struct {
	Body WebsiteRequest
}

type websiteIDRequest struct {
	WebsiteID string `path:"websiteID" format:"uuid"`
}

type updateWebsiteRequest struct {
	WebsiteID string `path:"websiteID" format:"uuid"`
	Body      WebsiteRequest
}

type emptyRequest struct{}

func (h Handler) List(ctx context.Context, _ *emptyRequest) (*listOutput, error) {
	websites, err := h.store.ListWebsites(ctx, h.currentUser(ctx).ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("load websites failed")
	}

	return newListOutput(ToWebsites(websites)), nil
}

func (h Handler) Create(ctx context.Context, request *websiteRequest) (*websiteOutput, error) {
	body := normalizeWebsiteRequest(request.Body)
	if body.Name == "" {
		return nil, huma.Error400BadRequest("name cannot be empty")
	}

	website, err := h.store.CreateWebsite(ctx, h.currentUser(ctx).ID, body.Name, body.Domain)
	if err != nil {
		return nil, huma.Error500InternalServerError("create website failed")
	}

	return newWebsiteOutput(ToWebsite(website)), nil
}

func (h Handler) Get(ctx context.Context, request *websiteIDRequest) (*websiteOutput, error) {
	website, err := h.store.GetWebsite(ctx, h.currentUser(ctx).ID, request.WebsiteID)
	if err != nil {
		return nil, h.websiteLookupError(err)
	}

	return newWebsiteOutput(ToWebsite(website)), nil
}

func (h Handler) Update(ctx context.Context, request *updateWebsiteRequest) (*websiteOutput, error) {
	body := normalizeWebsiteRequest(request.Body)
	if body.Name == "" {
		return nil, huma.Error400BadRequest("name cannot be empty")
	}

	user := h.currentUser(ctx)
	if err := h.store.UpdateWebsite(ctx, user.ID, request.WebsiteID, body.Name, body.Domain); err != nil {
		return nil, h.websiteLookupError(err)
	}

	website, err := h.store.GetWebsite(ctx, user.ID, request.WebsiteID)
	if err != nil {
		return nil, h.websiteLookupError(err)
	}

	return newWebsiteOutput(ToWebsite(website)), nil
}

func (h Handler) Delete(ctx context.Context, request *websiteIDRequest) (*okOutput, error) {
	if err := h.store.DeleteWebsite(ctx, h.currentUser(ctx).ID, request.WebsiteID); err != nil {
		return nil, h.websiteLookupError(err)
	}

	return newOKOutput(), nil
}

func normalizeWebsiteRequest(request WebsiteRequest) WebsiteRequest {
	request.Name = strings.TrimSpace(request.Name)
	request.Domain = strings.TrimSpace(request.Domain)
	return request
}
