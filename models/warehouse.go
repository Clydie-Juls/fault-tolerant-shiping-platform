package models

import (
	"rabbitmq/utils"
	"time"

	"github.com/jftuga/geodist"
)

type Warehouse struct {
	Id        int64
	location  Location
	Inventory Inventory
}

func NewWarehouse(location Location) *Warehouse {
	return &Warehouse{
		Id:       time.Now().UTC().Unix(),
		location: location,
	}
}

func (w *Warehouse) EstimateTimeToDestination(destination Location) time.Duration {
	geoSource := geodist.Coord{Lat: w.location.Latitude, Lon: w.location.Longitude}
	geoDestination := geodist.Coord{Lat: destination.Latitude, Lon: destination.Longitude}
	miles, _, err := geodist.VincentyDistance(geoSource, geoDestination)
	if err != nil {
		utils.FailOnError(err, "unable to compute distance between source and destination")
	}

	estimateTime := time.Duration(miles) * 2 * time.Second
	return estimateTime
}
