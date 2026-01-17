package appointment

import (
	"go-backend-project/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.GET("/appointments", h.getAppointments)
	r.POST("/appointments", h.create)
}

func (h *Handler) getAppointments(c *gin.Context) {
	appointments, err := h.service.GetAppointments(c.Request.Context())
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}
	utils.OK(c, "Appointments retrieved", appointments)
}

func (h *Handler) create(c *gin.Context) {
	var body Appointment
	if err := c.BindJSON(&body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	appointment, err := h.service.Create(c.Request.Context(), body)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Created(c, "Appointment created", appointment)
}
