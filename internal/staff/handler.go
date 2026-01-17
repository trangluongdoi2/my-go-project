package staff

import (
	"go-backend-project/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.GET("/staffs", h.list)
	r.POST("/staffs", h.create)
	r.PUT("/staffs/:id", h.update)
	r.DELETE("/staffs/:id", h.delete)
}

func (h *Handler) list(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	staffs, err := h.service.List(c.Request.Context(), page, limit)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}
	utils.OK(c, "Staff list retrieved", staffs)
}

func (h *Handler) create(c *gin.Context) {
	var body Staff
	if err := c.BindJSON(&body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	staff, err := h.service.Create(c.Request.Context(), body)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Created(c, "Staff created", staff)
}

func (h *Handler) update(c *gin.Context) {
	id := c.Param("id")
	var body Staff
	if err := c.BindJSON(&body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	staff, err := h.service.Update(c.Request.Context(), id, body)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.OK(c, "Staff updated", staff)
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
