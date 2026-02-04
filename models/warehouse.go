package models

import (
	"log"
	"rabbitmq/utils"
	"time"

	"github.com/jftuga/geodist"
)

type Warehouse struct {
	Id        int64     `json:"id"`
	Location  Location  `json:"location"`
	Inventory Inventory `json:"inventory"`
}

func NewWarehouse(location Location) *Warehouse {
	return &Warehouse{
		Id:       time.Now().UTC().Unix(),
		Location: location,
	}
}

func (w *Warehouse) EstimateTimeToDestination(destination Location) time.Duration {
	geoSource := geodist.Coord{Lat: w.Location.Latitude, Lon: w.Location.Longitude}
	geoDestination := geodist.Coord{Lat: destination.Latitude, Lon: destination.Longitude}
	miles, _, err := geodist.VincentyDistance(geoSource, geoDestination)

	latDiff := destination.Latitude - w.Location.Latitude
	lonDiff := destination.Longitude - w.Location.Longitude

	log.Printf(
		"[geo] src(%.6f,%.6f) dst(%.6f,%.6f) Δ(lat=%.6f, lon=%.6f), miles:%f",
		w.Location.Latitude,
		w.Location.Longitude,
		destination.Latitude,
		destination.Longitude,
		latDiff,
		lonDiff,
		miles,
	)

	if err != nil {
		utils.FailOnError(err, "unable to compute distance between source and destination")
	}

	estimateTime := time.Duration(miles * 3 * float64(time.Second))
	return estimateTime
}
