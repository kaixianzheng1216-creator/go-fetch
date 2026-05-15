package handler

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/model"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

type WebsiteHandler struct {
	websites           service.Website
	currentUser        func(context.Context) model.User
	websiteLookupError func(error) error
}

func NewWebsite(websites service.Website, currentUser func(context.Context) model.User, websiteLookupError func(error) error) WebsiteHandler {
	return WebsiteHandler{websites: websites, currentUser: currentUser, websiteLookupError: websiteLookupError}
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

func (handler WebsiteHandler) ListWebsites(ctx context.Context, _ *emptyRequest) (*WebsiteListOutput, error) {
	websites, err := handler.websites.List(ctx, handler.currentUser(ctx).ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("加载站点列表失败")
	}

	return NewWebsiteListOutput(ToWebsites(websites)), nil
}

func (handler WebsiteHandler) CreateWebsite(ctx context.Context, input *websiteRequest) (*WebsiteOutput, error) {
	website, err := handler.websites.Create(ctx, handler.currentUser(ctx).ID, input.Body.Name, input.Body.Domain)
	if err != nil {
		if errors.Is(err, service.ErrInvalidWebsiteName) {
			return nil, huma.Error400BadRequest("站点名称不能为空")
		}
		return nil, huma.Error500InternalServerError("创建站点失败")
	}

	return NewWebsiteOutput(ToWebsite(website)), nil
}

func (handler WebsiteHandler) GetWebsite(ctx context.Context, input *websiteIDRequest) (*WebsiteOutput, error) {
	website, err := handler.websites.Get(ctx, handler.currentUser(ctx).ID, input.WebsiteID)
	if err != nil {
		return nil, handler.websiteLookupError(err)
	}

	return NewWebsiteOutput(ToWebsite(website)), nil
}

func (handler WebsiteHandler) UpdateWebsite(ctx context.Context, input *updateWebsiteRequest) (*WebsiteOutput, error) {
	website, err := handler.websites.Update(ctx, handler.currentUser(ctx).ID, input.WebsiteID, input.Body.Name, input.Body.Domain)
	if err != nil {
		if errors.Is(err, service.ErrInvalidWebsiteName) {
			return nil, huma.Error400BadRequest("站点名称不能为空")
		}
		return nil, handler.websiteLookupError(err)
	}

	return NewWebsiteOutput(ToWebsite(website)), nil
}

func (handler WebsiteHandler) DeleteWebsite(ctx context.Context, input *websiteIDRequest) (*OKOutput, error) {
	if err := handler.websites.Delete(ctx, handler.currentUser(ctx).ID, input.WebsiteID); err != nil {
		return nil, handler.websiteLookupError(err)
	}

	return NewOKOutput(), nil
}
