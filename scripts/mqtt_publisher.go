package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	broker := getenv("MQTT_BROKER_URL", "tcp://localhost:1883")
	vehicle := getenv("VEHICLE_ID", "B1234XYZ")
	topic := fmt.Sprintf("/fleet/vehicle/%s/location", vehicle)

	opts := paho.NewClientOptions().AddBroker(broker).SetClientID("mock-pub-" + vehicle)
	client := paho.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil { panic(token.Error()) }

	// sekitar Monas
	baseLat, baseLon := -6.1754, 106.8272

	t := time.NewTicker(2 * time.Second)
	defer t.Stop()
	for now := range t.C {
		// random small jitter +- 120m
		dLat := (rand.Float64()-0.5)*0.002
		dLon := (rand.Float64()-0.5)*0.002 / math.Cos(baseLat*math.Pi/180)
		payload := map[string]any{
			"vehicle_id": vehicle,
			"latitude":   baseLat + dLat,
			"longitude":  baseLon + dLon,
			"timestamp":  now.Unix(),
		}
		b, _ := json.Marshal(payload)
		token := client.Publish(topic, 1, false, b)
		token.Wait()
		fmt.Println("published:", string(b))
	}
}

func getenv(k, d string) string { if v := getenv0(k); v != "" { return v }; return d }
func getenv0(k string) string { return map[string]string{}[k] } // replaced by os.Getenv in real use
