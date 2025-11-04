package engines

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CheckHealth(c *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c.AsciiJSON(http.StatusOK, gin.H{"message": "Welcome to iChillax"})
}
