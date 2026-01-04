package serviceoffering

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
	r.GET("/service-offerings", h.getServiceOfferings)
	r.GET("/service-offerings/:id", h.getServiceOfferingByID)
	r.POST("/service-offerings", h.create)
	r.PUT("/service-offerings/:id", h.update)
	r.DELETE("/service-offerings/:id", h.delete)
}

// GetServiceOfferings godoc
// @Summary Get all service offerings
// @Description Get list of all service offerings
// @Tags service-offerings
// @Accept json
// @Produce json
// @Success 200 {array} ServiceOffering
// @Failure 500 {object} map[string]string
// @Router /service-offerings [get]
func (h *Handler) getServiceOfferings(c *gin.Context) {
	offerings, err := h.service.GetServiceOfferings(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, offerings)
}

// GetServiceOfferingByID godoc
// @Summary Get service offering by ID
// @Description Get a single service offering by ID
// @Tags service-offerings
// @Accept json
// @Produce json
// @Param id path string true "Service Offering ID"
// @Success 200 {object} ServiceOffering
// @Failure 404 {object} map[string]string
// @Router /service-offerings/{id} [get]
func (h *Handler) getServiceOfferingByID(c *gin.Context) {
	id := c.Param("id")
	offering, err := h.service.GetServiceOfferingByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service offering not found"})
		return
	}
	c.JSON(http.StatusOK, offering)
}

// CreateServiceOffering godoc
// @Summary Create a new service offering
// @Description Create a new service offering
// @Tags service-offerings
// @Accept json
// @Produce json
// @Param body body ServiceOffering true "Service Offering"
// @Success 201 {object} ServiceOffering
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-offerings [post]
func (h *Handler) create(c *gin.Context) {
	var body ServiceOffering
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offering, err := h.service.Create(c.Request.Context(), body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, offering)
}

// UpdateServiceOffering godoc
// @Summary Update a service offering
// @Description Update an existing service offering
// @Tags service-offerings
// @Accept json
// @Produce json
// @Param id path string true "Service Offering ID"
// @Param body body ServiceOffering true "Service Offering"
// @Success 200 {object} ServiceOffering
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-offerings/{id} [put]
func (h *Handler) update(c *gin.Context) {
	id := c.Param("id")
	var body ServiceOffering
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offering, err := h.service.Update(c.Request.Context(), id, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, offering)
}

// DeleteServiceOffering godoc
// @Summary Delete a service offering
// @Description Delete a service offering by ID
// @Tags service-offerings
// @Accept json
// @Produce json
// @Param id path string true "Service Offering ID"
// @Success 204
// @Failure 500 {object} map[string]string
// @Router /service-offerings/{id} [delete]
func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
