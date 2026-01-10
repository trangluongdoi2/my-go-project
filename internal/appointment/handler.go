package appointment

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/appointments", h.getAppointments)
	r.POST("/appointments", h.create)
}

// GetAppointments godoc
// @Summary Get all appointments
// @Description Get list of all appointments
// @Tags appointments
// @Accept json
// @Produce json
// @Success 200 {array} Appointment
// @Failure 500 {object} map[string]string
// @Router /appointments [get]
func (h *Handler) getAppointments(c *gin.Context) {
	appointments, err := h.service.GetAppointments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appointments)
}

// CreateAppointment godoc
// @Summary Create a new appointment
// @Description Create a new appointment
// @Tags appointments
// @Accept json
// @Produce json
// @Param body body Appointment true "Appointment"
// @Success 201 {object} Appointment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /appointments [post]
func (h *Handler) create(c *gin.Context) {
	var body Appointment
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appointment, err := h.service.Create(c.Request.Context(), body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, appointment)
}
