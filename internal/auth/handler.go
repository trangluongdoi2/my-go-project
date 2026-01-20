package auth

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

func (h *Handler) RegisterPublicRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", h.login)
		auth.POST("/register", h.register)
		auth.POST("/refresh", h.refresh)
		auth.POST("/logout", h.logout)
	}
}

func (h *Handler) RegisterProtectedRoutes(r gin.IRouter) {
	auth := r.Group("/auth")
	{
		auth.GET("/me", h.getMe)
	}
}

func (h *Handler) login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	response, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.OK(c, "Login successful", response)
}

func (h *Handler) register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	response, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, "Registration successful", response)
}

func (h *Handler) refresh(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	response, err := h.service.RefreshToken(c.Request.Context(), req)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.OK(c, "Token refreshed", response)
}

func (h *Handler) getMe(c *gin.Context) {
	userID := c.GetString("user_id")

	response, err := h.service.GetMe(c.Request.Context(), userID)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.OK(c, "User retrieved", response)
}

func (h *Handler) logout(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	err := h.service.Logout(c.Request.Context(), req.RefreshToken)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.OK(c, "Logout successfully!", nil)
}
