package portals

import "github.com/gin-gonic/gin"

func wrapError(err error) gin.H {
	return gin.H{"error": err.Error()}
}
