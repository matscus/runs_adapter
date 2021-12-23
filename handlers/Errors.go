package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func CheckSQLError(c *gin.Context, err error) {
	if err == sql.ErrNoRows {
		c.JSON(200, gin.H{"Status": "error", "Message": "No data found for the specified parameters"})
		return
	}
	c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
}
