package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetFranchises(c *gin.Context) {
	ctx, cancel := getContext(10)
	defer cancel()

	franchises, err := h.franchises.List(ctx)
	if err != nil {
		log.Printf("GetFranchises error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch franchises"})
		return
	}
	c.JSON(http.StatusOK, franchises)
}

func (h *Handler) GetFranchise(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing franchise key"})
		return
	}
	ctx, cancel := getContext(10)
	defer cancel()

	detail, err := h.franchises.GetByKey(ctx, key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Franchise not found"})
		return
	}
	c.JSON(http.StatusOK, detail)
}
