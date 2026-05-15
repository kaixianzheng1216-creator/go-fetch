package websites

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	userdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/user"
	websitedomain "github.com/kaixianzheng1216-creator/go-fetch/internal/website"
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

func (handler Handler) ListWebsites(ctx context.Context, _ *emptyRequest) (*listOutput, error) {
	websites, err := handler.store.ListWebsites(ctx, handler.currentUser(ctx).ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("加载站点列表失败")
	}

	return newListOutput(ToWebsites(websites)), nil
}

func (handler Handler) CreateWebsite(ctx context.Context, request *websiteRequest) (*websiteOutput, error) {
	body := normalizeWebsiteRequest(request.Body)
	if body.Name == "" {
		return nil, huma.Error400BadRequest("站点名称不能为空")
	}

	website, err := handler.store.CreateWebsite(ctx, handler.currentUser(ctx).ID, body.Name, body.Domain)
	if err != nil {
		return nil, huma.Error500InternalServerError("创建站点失败")
	}

	return newWebsiteOutput(ToWebsite(website)), nil
}

func (handler Handler) GetWebsite(ctx context.Context, request *websiteIDRequest) (*websiteOutput, error) {
	website, err := handler.store.GetWebsite(ctx, handler.currentUser(ctx).ID, request.WebsiteID)
	if err != nil {
		return nil, handler.websiteLookupError(err)
	}

	return newWebsiteOutput(ToWebsite(website)), nil
}

func (handler Handler) UpdateWebsite(ctx context.Context, request *updateWebsiteRequest) (*websiteOutput, error) {
	body := normalizeWebsiteRequest(request.Body)
	if body.Name == "" {
		return nil, huma.Error400BadRequest("站点名称不能为空")
	}

	user := handler.currentUser(ctx)
	if err := handler.store.UpdateWebsite(ctx, user.ID, request.WebsiteID, body.Name, body.Domain); err != nil {
		return nil, handler.websiteLookupError(err)
	}

	website, err := handler.store.GetWebsite(ctx, user.ID, request.WebsiteID)
	if err != nil {
		return nil, handler.websiteLookupError(err)
	}

	return newWebsiteOutput(ToWebsite(website)), nil
}

func (handler Handler) DeleteWebsite(ctx context.Context, request *websiteIDRequest) (*okOutput, error) {
	if err := handler.store.DeleteWebsite(ctx, handler.currentUser(ctx).ID, request.WebsiteID); err != nil {
		return nil, handler.websiteLookupError(err)
	}

	return newOKOutput(), nil
}

func normalizeWebsiteRequest(request WebsiteRequest) WebsiteRequest {
	request.Name = strings.TrimSpace(request.Name)
	request.Domain = strings.TrimSpace(request.Domain)
	return request
}
