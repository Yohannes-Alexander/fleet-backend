package repository

import "context"
import "fleet-backend/internal/entity"

type LocationRepository interface {
	Save(ctx context.Context, loc entity.Location) error
	GetLatest(ctx context.Context, vehicleID string) (*entity.Location, error)
	GetHistory(ctx context.Context, vehicleID string, start, end int64) ([]entity.Location, error)
}
