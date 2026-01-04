package staff

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
	r.GET("/staffs", h.list)
}

// GetStaffs godoc
// @Summary Get all staffs
// @Description Get list of all staffs
// @Tags staffs
// @Accept json
// @Produce json
// @Success 200 {array} Staff
// @Failure 500 {object} map[string]string
// @Router /staffs [get]
func (h *Handler) list(c *gin.Context) {
	staffs, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, staffs)
}
