package controllers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Raihanpoke/shop-clothing/database"
	"github.com/Raihanpoke/shop-clothing/models"
	"github.com/gin-gonic/gin"
)

func AddToCart(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New(400, "user id is empty"))
			return
		}

		productID, err := strconv.Atoi(productQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusBadRequest, errors.New(400, "product id is empty"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := database.AddProductToCart(ctx, db, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerEroror, errors.New("donts't add product to cart"))
		}
		c.IndentedJSON(200, "successfuly added to the cart")
	}
}

func RemoveItem(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("produc id is empty"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ := c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := strconv.Atoi(productQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Backgroud(), 5*time.Second)
		defer cancel()

		err := database.RemoveCartItem(ctx, db, productID, userQueryID)
		if err != nil {
			c.IntendedJSON(http.StatusInternalServer, err)
			return
		}
		c.IndentedJSON(200, "succesfuly remove item from cart")
	}
}

func GetItemFromCart(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")

		if userID == "" {
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Backgroud(), 100*time.Second)

		var filledCart models.Users
		query := "SELECT *FROM users WHERE id = ?"
		err := db.QueryRowContext(ctx, query, userID).Scan(&filledCart.ID, &filledCart.UserCart)
		if err != nil {
			c.IndentedJSON(hhtp.StatusInternalServerEroror, "Not Found")
			return
		}

		var amount float64
		query := "SELECT SUM(price) FROM user_cart WHERE user_id = ?"
		err := db.QueryRowContext(ctx, query, UserID).Scan(&amount)
		if err != nil {
			log.Println(err)
		}

		response := struct {
			Amount   float64              `json:"amount"`
			UserCart []models.ProductUser `json:"user_cart"`
		}{
			Amount:   amount,
			UserCart: filledCart.UserCart,
		}

		c.IndentedJSON(hhtp.StatusOK, response)
	}
}

func ByFromCart(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")
		if userQueryID == "" {
			log.panicln("user id is empty")
			_ := c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		}

		ctx, cancel := context.WithTimeout(context.Backgroud(), 100*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(ctx, db, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IntendedJSON(200, "successfully places to order")
	}
}

func InstantBuy(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			c.AbortWithStatus(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := db.Exec("SELECT UNHEX(?, '-', '')", productQueryID).LastInsertId()
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Backgroud(), 5*time.Second)
		defer cancel()

		err := database.InstandBuyer(ctx, app, productID, userQueryID)
		if err != nil {
			c.IntendedJSON(http.StatusInternalServeError, err)
			return
		}
		c.IntendedJSON(200, "successfully places to order")
	}
}
