package store

import (
	"context"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/collector"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	storedb "github.com/kaixianzheng1216-creator/go-fetch/internal/store/db"

	"github.com/google/uuid"
)

func (s *Store) SaveEvent(ctx context.Context, input domain.EventInput) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	websiteUUID, err := uuid.Parse(input.WebsiteID)
	if err != nil {
		return err
	}

	sessionID, err := uuid.Parse(input.SessionID)
	if err != nil {
		return err
	}

	visitID, err := uuid.Parse(input.VisitID)
	if err != nil {
		return err
	}

	qtx := s.queries.WithTx(tx)
	if err := qtx.InsertSession(ctx, storedb.InsertSessionParams{
		ID:        sessionID,
		WebsiteID: websiteUUID,
		Browser:   input.Browser,
		Os:        input.OS,
		Device:    input.Device,
		Screen:    input.Screen,
		Language:  input.Language,
		Country:   input.Country,
		CreatedAt: input.CreatedAt,
	}); err != nil {
		return err
	}

	eventID := uuid.New()
	if err := qtx.InsertEvent(ctx, storedb.InsertEventParams{
		ID:             eventID,
		WebsiteID:      websiteUUID,
		SessionID:      sessionID,
		VisitID:        visitID,
		EventType:      int32(input.EventType),
		EventName:      input.EventName,
		UrlPath:        input.URLPath,
		UrlQuery:       input.URLQuery,
		ReferrerPath:   input.ReferrerPath,
		ReferrerDomain: input.ReferrerDomain,
		PageTitle:      input.PageTitle,
		Hostname:       input.Hostname,
		UtmSource:      input.UTMSource,
		UtmMedium:      input.UTMMedium,
		UtmCampaign:    input.UTMCampaign,
		UtmContent:     input.UTMContent,
		UtmTerm:        input.UTMTerm,
		Browser:        input.Browser,
		Os:             input.OS,
		Device:         input.Device,
		Screen:         input.Screen,
		Language:       input.Language,
		Country:        input.Country,
		CreatedAt:      input.CreatedAt,
	}); err != nil {
		return err
	}

	for _, item := range collector.FlattenData(input.Data) {
		if err := qtx.InsertEventData(ctx, storedb.InsertEventDataParams{
			ID:          uuid.New(),
			WebsiteID:   websiteUUID,
			EventID:     eventID,
			DataKey:     item.Key,
			StringValue: item.StringValue,
			NumberValue: pgFloat(item.NumberValue),
			CreatedAt:   input.CreatedAt,
		}); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
