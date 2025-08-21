CREATE TABLE IF NOT EXISTS vehicle_locations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  vehicle_id TEXT NOT NULL,
  latitude DOUBLE PRECISION NOT NULL,
  longitude DOUBLE PRECISION NOT NULL,
  ts_unix BIGINT NOT NULL,            -- epoch seconds
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_vehicle_ts ON vehicle_locations (vehicle_id, ts_unix DESC);
