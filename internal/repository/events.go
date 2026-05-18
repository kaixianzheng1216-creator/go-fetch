package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	storesqlc "github.com/kaixianzheng1216-creator/go-fetch/internal/repository/sqlc"
)

func (store *Store) SaveEvent(ctx context.Context, event domain.EventRecord) error {
	err := pgx.BeginFunc(ctx, store.pool, func(tx pgx.Tx) error {
		queries := store.queries.WithTx(tx)
		if err := queries.InsertSession(ctx, insertSessionParams(event)); err != nil {
			return fmt.Errorf("insert session: %w", err)
		}

		eventID := uuid.New()
		if err := queries.InsertEvent(ctx, insertEventParams(eventID, event)); err != nil {
			return fmt.Errorf("insert event: %w", err)
		}

		for _, item := range domain.FlattenEventData(event.Data) {
			if err := queries.InsertEventData(ctx, insertEventDataParams(event, eventID, item)); err != nil {
				return fmt.Errorf("insert event data: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("save event transaction: %w", err)
	}

	return nil
}

func insertSessionParams(event domain.EventRecord) storesqlc.InsertSessionParams {
	return storesqlc.InsertSessionParams{
		ID:         event.SessionID,
		WebsiteID:  event.WebsiteID,
		Browser:    event.Browser,
		OS:         event.OS,
		Device:     event.Device,
		Screen:     event.Screen,
		Language:   event.Language,
		Country:    event.Country,
		Region:     event.Region,
		City:       event.City,
		DistinctID: event.DistinctID,
		CreatedAt:  event.CreatedAt,
	}
}

func insertEventParams(eventID uuid.UUID, event domain.EventRecord) storesqlc.InsertEventParams {
	return storesqlc.InsertEventParams{
		ID:             eventID,
		WebsiteID:      event.WebsiteID,
		SessionID:      event.SessionID,
		VisitID:        event.VisitID,
		EventType:      int32(event.EventType),
		EventName:      event.EventName,
		URLPath:        event.URLPath,
		URLQuery:       event.URLQuery,
		ReferrerPath:   event.ReferrerPath,
		ReferrerQuery:  event.ReferrerQuery,
		ReferrerDomain: event.ReferrerDomain,
		PageTitle:      event.PageTitle,
		Hostname:       event.Hostname,
		UTMSource:      event.UTMSource,
		UTMMedium:      event.UTMMedium,
		UTMCampaign:    event.UTMCampaign,
		UTMContent:     event.UTMContent,
		UTMTerm:        event.UTMTerm,
		CreatedAt:      event.CreatedAt,
	}
}

func insertEventDataParams(event domain.EventRecord, eventID uuid.UUID, item domain.FlatEventData) storesqlc.InsertEventDataParams {
	return storesqlc.InsertEventDataParams{
		ID:          uuid.New(),
		WebsiteID:   event.WebsiteID,
		EventID:     eventID,
		DataKey:     item.Key,
		StringValue: item.StringValue,
		NumberValue: pgFloat8(item.NumberValue),
		DateValue:   pgOptionalTime(item.DateValue),
		DataType:    int32(item.DataType),
		CreatedAt:   event.CreatedAt,
	}
}
