package repository

import (
	"context"
	"database/sql"
	"fleet-backend/internal/entity"
)

type LocationPg struct{ db *sql.DB }

func NewLocationPg(db *sql.DB) *LocationPg { return &LocationPg{db: db} }

func (r *LocationPg) Save(ctx context.Context, loc entity.Location) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO vehicle_locations (vehicle_id, latitude, longitude, ts_unix)
		 VALUES ($1,$2,$3,$4)`,
		loc.VehicleID, loc.Latitude, loc.Longitude, loc.TsUnix)
	return err
}

func (r *LocationPg) GetLatest(ctx context.Context, vehicleID string) (*entity.Location, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT vehicle_id, latitude, longitude, ts_unix
		   FROM vehicle_locations
		  WHERE vehicle_id=$1
		  ORDER BY ts_unix DESC
		  LIMIT 1`, vehicleID)
	var loc entity.Location
	if err := row.Scan(&loc.VehicleID, &loc.Latitude, &loc.Longitude, &loc.TsUnix); err != nil {
		return nil, err
	}
	return &loc, nil
}

func (r *LocationPg) GetHistory(ctx context.Context, vehicleID string, start, end int64) ([]entity.Location, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT vehicle_id, latitude, longitude, ts_unix
		   FROM vehicle_locations
		  WHERE vehicle_id=$1 AND ts_unix BETWEEN $2 AND $3
		  ORDER BY ts_unix ASC`, vehicleID, start, end)
	if err != nil { return nil, err }
	defer rows.Close()

	var res []entity.Location
	for rows.Next() {
		var loc entity.Location
		if err := rows.Scan(&loc.VehicleID, &loc.Latitude, &loc.Longitude, &loc.TsUnix); err != nil {
			return nil, err
		}
		res = append(res, loc)
	}
	return res, rows.Err()
}
