package middleware

import (
	token "github.com/Raihanpoke/shop-clothing/tokens"
	"net/http"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc{
	return func(c *gin.Context){
		ClientToken := c.Request.Header.Get("token")
		if ClientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"No authorization header provided"})
			c.Abort()
			return
		}
		claims, err := tokens.ValidateToken(ClientToken)
		if err != nil{
			c.JSON(hhtp.StatusInternalServer, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("uid", claims.uid)
		c.Next()
	}
}