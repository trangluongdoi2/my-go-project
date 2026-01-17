package staff

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
	r.GET("/staffs", h.list)
}

func (h *Handler) list(c *gin.Context) {
	staffs, err := h.service.List(c.Request.Context())
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}
	utils.OK(c, "Staff list retrieved", staffs)
}
