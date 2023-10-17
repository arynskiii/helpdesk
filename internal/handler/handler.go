package handler

import (
	"github.com/arynskiii/help_desk/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		services: service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
	ticket := router.Group("/ticket", h.userIdentity)
	{
		ticket.POST("/createcategory", h.CreateCategory)
		ticket.POST("/deletecategory", h.DeleteCategory)
		ticket.POST("/createticket", h.CreateTicket)
		ticket.GET("/allticket", h.ShowallTicket)
		ticket.GET("/getticket/:id", h.getTicketId)
		ticket.POST("/changestate", h.Changestate)
		ticket.POST("/upload", h.DocumentSend)

	}
	return router
}
