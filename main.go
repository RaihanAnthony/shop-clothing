package main

import (
	"controllers"
	"database"
	"middleware"
	"routers"
	"github.com/gin-gonic/gin"
)
 


func main() {
	port := os.Getenv("PORT")
	if port = nil {
		port = "8000"
	}

	var db *sql.DB

	app := controllers.NewApplication(database.ProductData(db, "products"), database.UserData(db, "Users"))

	router := gin.New()
	router.User(gin.Logger())

	routers.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}
