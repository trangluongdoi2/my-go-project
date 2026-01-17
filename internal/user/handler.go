package user

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
	r.GET("/users", h.getUsers)
	r.GET("/users/:id", h.getUserByID)
	r.POST("/users", h.create)
	r.PUT("/users/:id", h.update)
	r.DELETE("/users/:id", h.delete)
}

func (h *Handler) getUsers(c *gin.Context) {
	users, err := h.service.GetUsers(c.Request.Context())
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}
	utils.OK(c, "Users retrieved", users)
}

func (h *Handler) getUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.service.GetUserByID(c.Request.Context(), id)
	if err != nil {
		utils.NotFound(c, "User not found")
		return
	}
	utils.OK(c, "User retrieved", user)
}

func (h *Handler) create(c *gin.Context) {
	var body User
	if err := c.BindJSON(&body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := h.service.Create(c.Request.Context(), body)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Created(c, "User created", user)
}

func (h *Handler) update(c *gin.Context) {
	id := c.Param("id")
	var body User
	if err := c.BindJSON(&body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := h.service.Update(c.Request.Context(), id, body)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.OK(c, "User updated", user)
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
