package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	svc "temporal-proj/service"
)

type Server struct {
	coord *svc.Coordinator
}

func NewServer(coord *svc.Coordinator) *Server {
	return &Server{coord: coord}
}

func (s *Server) Run(addr string) error {
	router := gin.New()
	router.Use(gin.Recovery())
	router.POST("/scan", s.handleScan)
	return router.Run(addr)
}

type scanRequest struct {
	URL string `json:"url" binding:"required"`
}

func (s *Server) handleScan(c *gin.Context) {
	var req scanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	scan, err := s.coord.Start(c.Request.Context(), req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"id":      scan.ID,
		"url":     scan.URL,
		"status":  scan.Status,
		"message": "Scan started",
	})
}
