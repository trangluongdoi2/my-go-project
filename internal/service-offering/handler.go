package serviceoffering

import (
	"go-backend-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.GET("/service-offerings", h.getServiceOfferings)
	r.GET("/service-offerings/:id", h.getServiceOfferingByID)
	r.POST("/service-offerings", h.create)
	r.PUT("/service-offerings/:id", h.update)
	r.DELETE("/service-offerings/:id", h.delete)
}

func (h *Handler) getServiceOfferings(c *gin.Context) {
	offerings, err := h.service.GetServiceOfferings(c.Request.Context(), nil)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}
	utils.OK(c, "Service offerings retrieved", offerings)
}

func (h *Handler) getServiceOfferingByID(c *gin.Context) {
	id := c.Param("id")
	offering, err := h.service.GetServiceOfferingByID(c.Request.Context(), id)
	if err != nil {
		utils.NotFound(c, "Service offering not found")
		return
	}
	utils.OK(c, "Service offering retrieved", offering)
}

func (h *Handler) create(c *gin.Context) {
	var body ServiceOffering
	if err := c.BindJSON(&body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	offering, err := h.service.Create(c.Request.Context(), body)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Created(c, "Service offering created", offering)
}

func (h *Handler) update(c *gin.Context) {
	id := c.Param("id")
	var body ServiceOffering
	if err := c.BindJSON(&body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	offering, err := h.service.Update(c.Request.Context(), id, body)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.OK(c, "Service offering updated", offering)
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
