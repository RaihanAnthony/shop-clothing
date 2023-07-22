package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/Raihanpoke/shop-clothing/models"
	"github.com/google/uuid"
)

var (
	ErrCantFindProduct    = errors.New("can't find the product")
	ErrCantDecodeProducts = errors.New("cant find the product")
	ErrUserIdIsNotValid   = errors.New("this user is not falid")
	ErrCantUpdateUser     = errors.New("cannot add the produt to the cart")
	ErrCantRemoveItemCart = errors.New("cannot delete item from the cart")
	ErrCantGetItem        = errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
)

func AddProductToCart(ctx context.Context, db *sql.DB, productID int, userID string) error {

	// Query the product from the mysql database
	query := "SELECT * FROM products WHERE id = ?"
	rows, err := db.QueryContext(ctx, query, productID)
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}
	defer rows.Close()

	var product models.ProductUser
	if rows.Next() {
		err := rows.Scan(&product.Product_ID, &product.Product_Name, &product.Price)
		if err != nil {
			log.Println(err)
			return ErrCantDecodeProducts
		}
	} else {
		return ErrCantFindProduct
	}

	// Update the user's cart in the MYSQL database
	query = "UPDATE user SET usercart = JSON_ARRAY_APPEND(usercart, '$, ?) WHERE id = ?"
	_, err = db.ExecContext(ctx, query, productID, userID)
	if err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}

	return nil
}

func RemoveCartItem(ctx context.Context, db *sql.DB, productID int, userID string) error {
	query := "UPDATE user SET usercart = JSON_REMOVE(usercart, JSON_UNQUOTE(JSON_SEARCH(usercart, 'one', ?))) WHERE id = ?"
	_, err := db.ExecContext(ctx, query, productID, userID)
	if err != nil {
		log.Println(err)
		return ErrCantRemoveItemCart
	}
	return nil
}

func BuyItemFromCart(ctx context.Context, db *sql.DB, userID string) error {
	//fetch the cart of the user
	//find the cart total
	//create an order with the items
	//added order to the user collection
	//added items in the cart to order list
	//empty up the cart
	var ordercart models.Order

	ordercart.Order_ID = uuid.New().String()
	ordercart.Ordered_At = time.Now()
	ordercart.Order_Cart = make([]models.ProductUser, 0)
	ordercart.Payment_Method = true

	// get the total price from the usercart
	var total_price int
	query := "SELECT SUM(usercart->'$.product_price') AS total_price FROM users WHERE user_id = ?"
	err := db.QueryRowContext(ctx, query, userID).Scan(&total_price)
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	ordercart.Price = total_price

	// Insert order into orders tabel
	query = "INSERT INTO orders(order_id, user_id, ordered_at, price, payment_method) VALUE(?, ?, ?, ?, ?)"
	_, err = db.ExecContext(ctx, query, ordercart.Order_ID, userID, ordercart.Ordered_At, ordercart.Price, ordercart.Payment_Method)
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	// Insert into from usercart into order_list table
	query = "INSERT INTO order_list (order_id, product_id, product_name, product_price) SELECT ?, JSON_UNQUOTE(JSON_EXTRACT(usercart, '$[*].product_id')), JSON_UNQUOTE(JSON_EXTRACT(usercart, '$[*].product_name')), JSON_UNQUOTE(JSON_EXTRACT(usercart, '$[*].product_price')) FROM users WHERE id = ?"
	_, err = db.ExecContext(ctx, query, ordercart.Order_ID, userID)
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	// Clear the usercart
	query = "UPDATE users SET usercart = '[]' WHERE user_id = ?"
	_, err = db.ExecContext(ctx, query, userID)
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	return nil
}

func InstandBuyer(ctx context.Context, db *sql.DB, productID int, userID string) error {
	var orders_detail models.Order

	orders_detail.Order_ID = uuid.New().String()
	orders_detail.Ordered_At = time.Now()
	orders_detail.Order_Cart = make([]models.ProductUser, 0)
	orders_detail.Payment_Method = true

	// get price from tabel product
	var price int
	query := "SELECT price FROM products WHERE id = ?"
	err := db.QueryRowContext(ctx, query, productID).Scan(&price)
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	orders_detail.Price = price

	// Insert order into tabel order
	query = "INSERT INTO orders(order_id, ordered_at, price, payment_method) VALUE(?, ?, ?, ?, ?)"
	_, err = db.ExecContext(ctx, query, orders_detail.Order_ID, orders_detail.Ordered_At, orders_detail.Price, orders_detail.Payment_Method)
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	// UPDATE COLOUM ORDERS IN TABEL USER
	query = "UPDATE users SET orders = JSON_ARRAY_APPEND(orders, '$', ?) WHERE id = ?"
	orderJSON, _ := json.Marshal(orders_detail)
	_, err = db.ExecContext(ctx, query, orderJSON, userID)
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	return nil
}
