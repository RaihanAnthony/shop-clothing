package controllers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/Raihanpoke/shop-clothing/models"
	"github.com/gin-gonic/gin"
)

func AddAddress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")

		if userID == "" {
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Search Index"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var AddressUsers models.Address

		// Assuming Address_id is an auto-increment primary key in mysql
		query := "INSERT INTO address (address_id) VALUES (NULL))"
		_, err := db.ExecContext(ctx, query)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "internal server error")
			return
		}

		AddressIDInt := AddressUsers.Address_ID

		var size int
		query = "SELECT COUNT(*) FROM address WHERE address_id = ?"
		err = db.QueryRowContext(ctx, query, AddressIDInt).Scan(&size)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "internal server error")
			return
		}

		if size < 2 {
			query := "UPDATE address SET House = CONCAT(House, ?) WHERE addres_id = ?"
			stmt, err := db.Prepare(query)
			if err != nil {
				c.JSON(http.StatusInternalServerError, "internal server error")
				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(AddressUsers.House, AddressIDInt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, "internal server error")
				return
			} else {
				c.JSON(http.StatusBadRequest, "Not Allowed")
			}
		}
	}
}

func EditHomeAddress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid seacrh index"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		query := "UPDATE address SET house_name = ?, street_name = ?, city_name = ?, pincode = ? WHERE id = ?"
		stmt, err := db.Prepare(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "internal server error")
			return
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, editaddress.House, editaddress.Street, editaddress.City, editaddress.Pincode, userID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "internal server error")
			return
		}

		c.IndentedJSON(http.StatusOK, "succesfully update the home address")
	}
}

func EditWorkAddress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid"})
			c.Abort()
			return
		}

		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		query := "UPDATE address SET address = JSON_SET (address, '$[1].house_name', ?,'$[1].Street_name', ?, '$[1].city_name', ?,'$[1].pin_code', ?) WHERE id = ?"
		stmt, err := db.Prepare(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, editaddress.House, editaddress.Street, editaddress.City, editaddress.Pincode, user_id)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		c.IndentedJSON(http.StatusOK, "succesfully update the data")
	}
}

func DeleteAddress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Seacrh Index"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		query := "DELETE FROM addresses WHERE id = ?"
		result, err := db.ExecContext(ctx, query, userID)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		if rowsAffected == 0 {
			c.IndentedJSON(http.StatusNotFound, "Addres not Found")
			return
		}

		c.IndentedJSON(http.StatusOK, "Succcessfully deleted")
	}
}
