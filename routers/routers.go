package routers

import (
	"github.com/Raihanpoke/shop-clothing/controllers"
	"github.com/gin-gonic/gin"
)

func UserRouters(incomingRouters *gin.Engine) {
	incomingRouters.POST("/users/signup", controllers.SignUp())
	incomingRouters.POST("/users/login", controllers.Login())
	incomingRouters.POST("/admin/addproduct", controllers.ProductViewerAdmin())
	incomingRouters.GET("/users/product", controllers.SearchProduct())
	incomingRouters.GET("/users/search", controllers.SearchProductByQuery())
}
