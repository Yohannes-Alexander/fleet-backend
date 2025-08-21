package route

import (
	"github.com/gin-gonic/gin"
	"fleet-backend/internal/handler"
	"fleet-backend/internal/middleware"
)

func Setup(r *gin.Engine, vh *handler.VehicleHandler) {
	r.Use(middleware.RequestLogger())
	api := r.Group("/vehicles")
	{
		api.GET("/:vehicle_id/location", vh.GetLatest)
		api.GET("/:vehicle_id/history", vh.GetHistory)
	}
}
