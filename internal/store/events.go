package store

import (
	"context"
	"fmt"

	eventdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/event"
	storesqlc "github.com/kaixianzheng1216-creator/go-fetch/internal/store/sqlc"

	"github.com/google/uuid"
)

func (store *Store) SaveEvent(ctx context.Context, input eventdomain.EventInput) error {
	transaction, err := store.databasePool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("开启保存事件事务失败: %w", err)
	}
	defer func() {
		_ = transaction.Rollback(ctx)
	}()

	websiteUUID, err := uuid.Parse(input.WebsiteID)
	if err != nil {
		return fmt.Errorf("解析站点 ID 失败: %w", err)
	}

	sessionID, err := uuid.Parse(input.SessionID)
	if err != nil {
		return fmt.Errorf("解析会话 ID 失败: %w", err)
	}

	visitID, err := uuid.Parse(input.VisitID)
	if err != nil {
		return fmt.Errorf("解析访问 ID 失败: %w", err)
	}

	transactionQueries := store.queries.WithTx(transaction)
	if err := transactionQueries.InsertSession(ctx, storesqlc.InsertSessionParams{
		ID:         sessionID,
		WebsiteID:  websiteUUID,
		Browser:    input.Browser,
		Os:         input.OS,
		Device:     input.Device,
		Screen:     input.Screen,
		Language:   input.Language,
		Country:    input.Country,
		Region:     input.Region,
		City:       input.City,
		DistinctID: input.DistinctID,
		CreatedAt:  input.CreatedAt,
	}); err != nil {
		return fmt.Errorf("写入会话失败: %w", err)
	}

	eventID := uuid.New()
	if err := transactionQueries.InsertEvent(ctx, storesqlc.InsertEventParams{
		ID:             eventID,
		WebsiteID:      websiteUUID,
		SessionID:      sessionID,
		VisitID:        visitID,
		EventType:      int32(input.EventType),
		EventName:      input.EventName,
		UrlPath:        input.URLPath,
		UrlQuery:       input.URLQuery,
		ReferrerPath:   input.ReferrerPath,
		ReferrerQuery:  input.ReferrerQuery,
		ReferrerDomain: input.ReferrerDomain,
		PageTitle:      input.PageTitle,
		Hostname:       input.Hostname,
		UtmSource:      input.UTMSource,
		UtmMedium:      input.UTMMedium,
		UtmCampaign:    input.UTMCampaign,
		UtmContent:     input.UTMContent,
		UtmTerm:        input.UTMTerm,
		CreatedAt:      input.CreatedAt,
	}); err != nil {
		return fmt.Errorf("写入事件失败: %w", err)
	}

	for _, item := range eventdomain.FlattenEventData(input.Data) {
		if err := transactionQueries.InsertEventData(ctx, storesqlc.InsertEventDataParams{
			ID:          uuid.New(),
			WebsiteID:   websiteUUID,
			EventID:     eventID,
			DataKey:     item.Key,
			StringValue: item.StringValue,
			NumberValue: pgNumeric(item.NumberValue),
			DateValue:   pgOptionalTime(item.DateValue),
			DataType:    int32(item.DataType),
			CreatedAt:   input.CreatedAt,
		}); err != nil {
			return fmt.Errorf("写入事件数据失败: %w", err)
		}
	}

	if err := transaction.Commit(ctx); err != nil {
		return fmt.Errorf("提交保存事件事务失败: %w", err)
	}

	return nil
}
