package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	storesqlc "github.com/kaixianzheng1216-creator/go-fetch/internal/repository/sqlc"
)

func (store *Store) SaveEvent(ctx context.Context, input domain.EventInput) error {
	err := pgx.BeginFunc(ctx, store.databasePool, func(tx pgx.Tx) error {
		queries := store.queries.WithTx(tx)
		if err := queries.InsertSession(ctx, storesqlc.InsertSessionParams{
			ID:         input.SessionID,
			WebsiteID:  input.WebsiteID,
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
			return fmt.Errorf("insert session: %w", err)
		}

		eventID := uuid.New()
		if err := queries.InsertEvent(ctx, storesqlc.InsertEventParams{
			ID:             eventID,
			WebsiteID:      input.WebsiteID,
			SessionID:      input.SessionID,
			VisitID:        input.VisitID,
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
			return fmt.Errorf("insert event: %w", err)
		}

		for _, item := range domain.FlattenEventData(input.Data) {
			if err := queries.InsertEventData(ctx, storesqlc.InsertEventDataParams{
				ID:          uuid.New(),
				WebsiteID:   input.WebsiteID,
				EventID:     eventID,
				DataKey:     item.Key,
				StringValue: item.StringValue,
				NumberValue: pgFloat8(item.NumberValue),
				DateValue:   pgOptionalTime(item.DateValue),
				DataType:    int32(item.DataType),
				CreatedAt:   input.CreatedAt,
			}); err != nil {
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
