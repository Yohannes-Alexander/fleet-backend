# Fleet Management System

Sistem manajemen armada bus yang dibangun menggunakan **Golang**, **PostgreSQL**, **RabbitMQ**, **MQTT**, dan **Docker Compose**.  
Fitur utama:
- Menerima data lokasi kendaraan melalui **MQTT**.
- Menyimpan data lokasi ke **PostgreSQL**.
- Menyediakan **REST API** untuk mendapatkan lokasi terakhir dan riwayat perjalanan kendaraan.
- Menggunakan **RabbitMQ** untuk mem-publish event ketika kendaraan memasuki area **geofence**.

---

## 📂 Struktur Project
<br>├── cmd/ 
<br>│   ├── api/                # Main service (REST API)
<br>│   │   └── main.go
<br>│   ├── worker/             # Worker service (consume event dari RabbitMQ)
<br>│   │   └── main.go
<br>│   └── publisher/          # Publisher service (publish event ke RabbitMQ)
<br>│       └── main.go
<br>│
<br>├── config/                 # Konfigurasi (database, env, dll)
<br>│   └── config.go
<br>│
<br>├── internal/
<br>│   ├── handler/            # HTTP handler (controller)
<br>│   ├── usecase/            # Business logic (application layer)
<br>│   ├── dto/                # Data transfer object
<br>│   ├── entity/             # Entity/domain model
<br>│   ├── repository/         # Repository (akses database)
<br>│   ├── utils/              # Helper utilities
<br>│   ├── route/              # Routing untuk API
<br>│   └── middleware/         # Middleware (auth, logging, dll)
<br>│
<br>├── pkg/                    # Paket tambahan (infra/docker/driver)
<br>│   ├── db/                 # PostgreSQL setup
<br>   ├── mqtt/               # MQTT broker setup
<br>│   └── rabbit/             # RabbitMQ setup
<br>│
<br>├── scripts/                # Script tambahan (migrate, seed, dll)
<br>│
<br>├── docker-compose.yml      # Docker Compose untuk semua service
<br>├── Dockerfile.api          # Dockerfile untuk service API
<br>├── Dockerfile.worker       # Dockerfile untuk service Worker
<br>├── Dockerfile.publisher    # Dockerfile untuk service Publisher
<br>├── .env.example            # Contoh file environment variables
<br>├── go.mod                  # Modul Go
<br>├── go.sum                  # Dependency checksum
<br>└── README.md               # Dokumentasi project



---

## ⚙️ Set Up Project
Pastikan Anda sudah menginstall:
- [Docker](https://docs.docker.com/get-docker/)
- [Golang](https://go.dev/dl/)


---

## 🚀 Cara Menjalankan

1. **Clone repository**
   ```bash
   git clone https://github.com/Yohannes-Alexander/fleet-backend.git

   cd fleet-backend

   go mod tidy
   ```
2. **Buat file .env** <br>
    Contoh
   ```env
    APP_PORT=8080

    # PostgreSQL
    DB_HOST=postgres
    DB_PORT=5432
    DB_USER=postgres
    DB_PASS=mypassword
    DB_NAME=fleet
    DB_SSLMODE=disable

    # MQTT
    MQTT_BROKER_URL=tcp://mosquitto:1883
    MQTT_CLIENT_ID=fleet-consumer
    MQTT_TOPIC=/fleet/vehicle/+/location

    # RabbitMQ
    RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
    RABBITMQ_EXCHANGE=fleet.events
    RABBITMQ_EXCHANGE_TYPE=topic
    RABBITMQ_ROUTING_KEY=geofence.entry
    RABBITMQ_QUEUE=geofence_alerts

    # Geofence: titik pusat & radius meter
    GEOFENCE_LAT=-6.2088
    GEOFENCE_LON=106.8456
    GEOFENCE_RADIUS_M=50
    ```
3. **Jalankan Scripts SQL di pkg/db/init.sql di Database PostgreSQL**
   ```sql
    CREATE TABLE IF NOT EXISTS vehicle_locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vehicle_id TEXT NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    ts_unix BIGINT NOT NULL,          
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

    CREATE INDEX IF NOT EXISTS idx_vehicle_ts ON vehicle_locations (vehicle_id, ts_unix DESC);
   ```
4. **Jalankan aplikasi dengan Docker Compose**
   ```bash
    docker compose up --build
   ```
   
   Ini akan menjalankan service:

    -  app → Golang backend

    -  postgres → Database PostgreSQL

    -  rabbitmq → Message broker

    -  mosquitto → MQTT broker

    -  publisher → Simulasi publisher data kendaraan
        
## 📡 API Endpoint
1. **Get Lokasi Terakhir Kendaraan**
    - METHOD : GET
    - END POINT : /vehicles/{license_plate}/latest
    - Response : <br>
    ```json
    {
        "vehicle_id": "B1234XYZ",
        "latitude": -6.176029135733974,
        "longitude": 106.82746624076846,
        "timestamp": 1755785081
    }
    ```


2. **Get Lokasi History Kendaraan**
    - METHOD : GET
    - END POINT : /vehicles/{license_plate}/history?start={start}&end={end}latest    
    - Response : <br>
    ```json
        [
            {
                "vehicle_id": "B1234XYZ",
                "latitude": -6.175876803503382,
                "longitude": 106.82785858019305,
                "timestamp": 1755765230
            },
            {
                "vehicle_id": "B1234XYZ",
                "latitude": -6.174701682185257,
                "longitude": 106.82657124081923,
                "timestamp": 1755765232
            },
            {
                "vehicle_id": "B1234XYZ",
                "latitude": -6.174674383470626,
                "longitude": 106.82716223124062,
                "timestamp": 1755765234
            },
            {
                "vehicle_id": "B1234XYZ",
                "latitude": -6.1752759058510565,
                "longitude": 106.82657209412086,
                "timestamp": 1755765236
            }
        ]
    ```

## Postman Collection <br>
Untuk mempermudah pengujian API, silakan gunakan file Postman berikut:  

[📥 Download Postman Collection](https://github.com/Yohannes-Alexander/fleet-backend/raw/main/Fleet.postman_collection.json)
