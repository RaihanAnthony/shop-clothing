package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Raihanpoke/shop-clothing/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var UserData = database.UserData("users")
var ProductData = database.ProductData("product")

func HashPassword(password string) string {
	bytes, err := bcrypt.GeneateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""

	if err != nil {
		msg = "login or Password is incorrect"
		valid = false
	}
	return valid, msg
}

func SignUp(DB *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Backgroud(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.statusBadRequest, gin.H{"error": err.error()})
			return
		}
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.statusBadRequest, gin.H{"error": validationErr})
			return
		}

		// Check if user exists
		var count int
		queryEmail := "SELECT COUNT(*) FROM user where email = ?"
		err := db.QueryRowContext(ctx, queryEmail, user.Email).Scan(&count)
		if err != nil {
			log.panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
			return
		}

		queryPhone := "SELECT COUNT(*) FROM user WHERE phone = ?"
		err := db.QueryRowContext(ctx, queryPhone, user.Phone).Scan(&count)
		if err != nil {
			log.panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this phone no. is already in use"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &Password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.User_ID = primitive.NewObjectID()
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, *user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		user.Addres_Details = make([]models.Addres, 0)
		user.Order_Status = make([]models.Order, 0)

		// Insert the user into the database
		query := "INSERT INTO user(email, password, created_at, update_at, user_id, token, refresh_token) VALUES(? , ?, ?, ?, ?, ?, ?)"
		_, err = db.ExecContext(ctx, query, user.Email, user.Password, user.Created_At, user.Updated_At, user.User_ID, user.Token, user.Refresh_Token)
		if err != nil {
			log.panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "successfuly signed in!"})
	}
}

func Login(DB *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel = context.WithTimeout(context.Backgroud(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		var founduser models.User
		query := "SELECT * FROM user WHERE email = ?"
		err := db.QueryRowContext(ctx, query, user.Email).scan(&founduser)
		if err != nil {
			if err == sql.ErrNoRows {
				// user not found, countinue with sign-up process
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		PasswordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password)

		if !PasswordIsValid {
			c.JSON(http.StatusInternalServeError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

		token, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.First_name, *founduser.Last_Name, *founduser.User_ID)
		defer cancel()

		generate.UpdateAllTokens(token, refreshToken, founduser.user_ID)

		c.JSON(http.StatusFound, founduser)
	}

}

func ProductViewerAdmin() gin.HandlerFunc {

}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product

		rows, err := DB.Query("SELECT * FROM products")
		if err != nil {
			c.IndentedJSON(hhtp.StatusInternalServerError, gin.H{"error": "something went wrong, plase try after some time"})
			time
		}
		defer rows.Close()

		var product models.Product
		err := rows.Scan(&product.ID, &product.Last_Name, &product.Price)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		productList = append(productList, product)

		if err := rows.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(hhtp.StatusBadRequest, gin.H{"error": "invalid"})
			return
		}

		c.IndentedJSON(hhtp.StatusOK, productList)

	}
}

func SearchProductByQuery(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryParam := c.Query("name")

		// cheack if it's empty
		if queryParam == "" {
			log.Println("query is empty")
			c.Header("content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid seacrh index"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Backgroud(), 100*time.Second)
		defer cancel()

		query := "SELECT * FROM products WHERE product_name LIKE ?"
		rows, err := db.QueryContext(ctx, query, "%"+queryParam+"%")
		if err != nil {
			log.Println(err)
			c.IndentedJSON(hhtp.StatusInternalServerError, "something went wrong while fetching the data")
			return
		}
		defer rows.Close()

		var searchProducts []models.produt

		for rows.Next() {
			var product models.Product
			err := rows.Scan(&product.ID, &product.Name, &product.Price) // adjust the column names accordingly
			if err != nil {
				log.Println(err)
				c.IndentedJSON(http.StatusInternalServerError, "Invalid data")
				return
			}
			SearchProducts = append(searchProducts, product)
		}

		if err := rows.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, "Invalid request")
			return
		}

		c.IndentedJSON(http.StatusOk, searchProducts)
	}
}
