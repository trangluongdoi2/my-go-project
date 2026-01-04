package booking

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
	r.GET("/bookings", h.getBookings)
	r.POST("/bookings", h.create)
}

// GetBookings godoc
// @Summary Get all bookings
// @Description Get list of all bookings
// @Tags bookings
// @Accept json
// @Produce json
// @Success 200 {array} Booking
// @Failure 500 {object} map[string]string
// @Router /bookings [get]
func (h *Handler) getBookings(c *gin.Context) {
	bookings, err := h.service.GetBookings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// CreateBooking godoc
// @Summary Create a new booking
// @Description Create a new booking
// @Tags bookings
// @Accept json
// @Produce json
// @Param body body Booking true "Booking"
// @Success 201 {object} Booking
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bookings [post]
func (h *Handler) create(c *gin.Context) {
	var body Booking
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := h.service.Create(c.Request.Context(), body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, booking)
}
