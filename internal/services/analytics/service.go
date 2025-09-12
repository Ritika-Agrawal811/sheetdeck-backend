package analytics

import (
	"context"
	"fmt"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/domain/dtos"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/repository"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/pkg/utils"
)

type AnalyticsService interface {
	RecordPageView(ctx context.Context, details dtos.PageviewRequest) error
}

type analyticsService struct {
	repo *repository.Queries
}

func NewAnalyticsService(repo *repository.Queries) AnalyticsService {
	return &analyticsService{
		repo: repo,
	}
}

/**
 * Record a page view
 * @param details dtos.PageviewRequest
 * @return error
 */
func (s *analyticsService) RecordPageView(ctx context.Context, details dtos.PageviewRequest) error {
	if details.Route == "" || details.UserAgent == "" || details.IpAddress == "" {
		return nil
	}

	pageViewParams := repository.StorePageviewParams{
		Pathname:  details.Route,
		IpAddress: details.IpAddress,
		UserAgent: details.UserAgent,
		Referrer:  utils.PgText(details.Referrer),
	}

	if err := s.repo.StorePageview(ctx, pageViewParams); err != nil {
		return fmt.Errorf("failed to record pageview: %w", err)
	}

	return nil
}
