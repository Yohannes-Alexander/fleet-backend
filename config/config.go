package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppPort string

	DBHost string
	DBPort int
	DBUser string
	DBPass string
	DBName string
	DBSSL  string

	MqttURL     string
	MqttClient  string
	MqttTopic   string

	RabbitURL    string
	RabbitEx     string
	RabbitExType string
	RabbitKey    string
	RabbitQueue  string

	GeofenceLat float64
	GeofenceLon float64
	GeofenceRad float64
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" { return v }
	return def
}

func mustFloat(key string, def float64) float64 {
	s := getenv(key, "")
	if s == "" { return def }
	f, _ := strconv.ParseFloat(s, 64); return f
}

func mustInt(key string, def int) int {
	s := getenv(key, "")
	if s == "" { return def }
	i, _ := strconv.Atoi(s); return i
}

func Load() *Config {
	return &Config{
		AppPort: getenv("APP_PORT", "8080"),

		DBHost: getenv("DB_HOST", "localhost"),
		DBPort: mustInt("DB_PORT", 5432),
		DBUser: getenv("DB_USER", "postgres"),
		DBPass: getenv("DB_PASS", "mypassword"),
		DBName: getenv("DB_NAME", "fleet"),
		DBSSL:  getenv("DB_SSLMODE", "disable"),

		MqttURL:    getenv("MQTT_BROKER_URL", "tcp://localhost:1883"),
		MqttClient: getenv("MQTT_CLIENT_ID", "fleet-consumer"),
		MqttTopic:  getenv("MQTT_TOPIC", "/fleet/vehicle/+/location"),

		RabbitURL:    getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		RabbitEx:     getenv("RABBITMQ_EXCHANGE", "fleet.events"),
		RabbitExType: getenv("RABBITMQ_EXCHANGE_TYPE", "topic"),
		RabbitKey:    getenv("RABBITMQ_ROUTING_KEY", "geofence.entry"),
		RabbitQueue:  getenv("RABBITMQ_QUEUE", "geofence_alerts"),

		GeofenceLat: mustFloat("GEOFENCE_LAT", -6.2088),
		GeofenceLon: mustFloat("GEOFENCE_LON", 106.8456),
		GeofenceRad: mustFloat("GEOFENCE_RADIUS_M", 50),
	}
}
