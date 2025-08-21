package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fleet-backend/internal/dto"
	"fleet-backend/internal/entity"
	"fleet-backend/internal/utils"
	"fleet-backend/pkg/rabbit"
	"fleet-backend/internal/repository"
)

type FleetUsecase struct {
	repo          repository.LocationRepository
	pub           *rabbit.Publisher
	geofenceLat   float64
	geofenceLon   float64
	geofenceMeter float64
	exchange      string
	routingKey    string
}

func NewFleetUsecase(repo repository.LocationRepository, pub *rabbit.Publisher,
	glat, glon, radius float64, exchange, routingKey string) *FleetUsecase {
	return &FleetUsecase{repo: repo, pub: pub, geofenceLat: glat, geofenceLon: glon, geofenceMeter: radius, exchange: exchange, routingKey: routingKey}
}

// Called by MQTT subscriber
func (uc *FleetUsecase) IngestLocation(ctx context.Context, m dto.MqttLocationMessage) error {
	// Basic validation already by binding; extra guard:
	if m.VehicleID == "" { return errors.New("vehicle_id required") }
	loc := entity.Location{VehicleID: m.VehicleID, Latitude: m.Latitude, Longitude: m.Longitude, TsUnix: m.Timestamp}
	if err := uc.repo.Save(ctx, loc); err != nil { return err }

	// Geofence check
	dist := utils.HaversineMeters(m.Latitude, m.Longitude, uc.geofenceLat, uc.geofenceLon)
	if dist <= uc.geofenceMeter {
		ev := map[string]any{
			"vehicle_id": m.VehicleID,
			"event":      "geofence_entry",
			"location": map[string]any{
				"latitude":  m.Latitude,
				"longitude": m.Longitude,
			},
			"timestamp": m.Timestamp,
		}
		body, _ := json.Marshal(ev)
		return uc.pub.Publish(uc.exchange, uc.routingKey, body)
	}
	return nil
}

func (uc *FleetUsecase) GetLatest(ctx context.Context, vehicleID string) (*entity.Location, error) {
	return uc.repo.GetLatest(ctx, vehicleID)
}

func (uc *FleetUsecase) GetHistory(ctx context.Context, vehicleID string, start, end int64) ([]entity.Location, error) {
	return uc.repo.GetHistory(ctx, vehicleID, start, end)
}
