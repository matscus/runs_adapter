package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func CheckSQLError(c *gin.Context, err error) {
	if err == sql.ErrNoRows {
		c.JSON(200, gin.H{"status": "error", "message": "No data found for the specified parameters"})
		return
	}
	c.JSON(500, gin.H{"status": "error", "message": err.Error()})
}
