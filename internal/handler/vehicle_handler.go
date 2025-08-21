package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"fleet-backend/internal/usecase"
	"fleet-backend/internal/dto"
)

type VehicleHandler struct{ uc *usecase.FleetUsecase }

func NewVehicleHandler(uc *usecase.FleetUsecase) *VehicleHandler { return &VehicleHandler{uc: uc} }

func (h *VehicleHandler) GetLatest(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	loc, err := h.uc.GetLatest(c.Request.Context(), vehicleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, dto.LocationResponse{
		VehicleID: loc.VehicleID, Latitude: loc.Latitude, Longitude: loc.Longitude, Timestamp: loc.TsUnix,
	})
}

func (h *VehicleHandler) GetHistory(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	start, _ := strconv.ParseInt(c.Query("start"), 10, 64)
	end, _ := strconv.ParseInt(c.Query("end"), 10, 64)
	locs, err := h.uc.GetHistory(c.Request.Context(), vehicleID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	out := make([]dto.LocationResponse, 0, len(locs))
	for _, v := range locs {
		out = append(out, dto.LocationResponse{
			VehicleID: v.VehicleID, Latitude: v.Latitude, Longitude: v.Longitude, Timestamp: v.TsUnix,
		})
	}
	c.JSON(http.StatusOK, out)
}
