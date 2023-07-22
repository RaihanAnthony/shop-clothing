package models

import (
	"time"
)

type User struct {
	ID             []byte        `gorm:"column:id;type:BLOB"`
	First_Name     *string       `gorm:"column:first_name" 	validate:"required,min=2,max=30"`
	Last_Name      *string       `gorm:"column:last_name" 	validate:"required,min=2,max=20"`
	Password       *string       `gorm:"column:password" 	validate:"required,min=6`
	Email          *string       `gorm:"column:email" 		validate:"email, required"`
	Phone          *string       `gorm:"column:phone" 		validate:"required"`
	Token          *string       `gorm:"column:token"`
	Refresh_Token  *string       `gorm:"column:refresh_token"`
	Created_Token  time.Time     `gorm:"column:created_token"`
	Updated_At     time.Time     `gorm:"column:update_at"`
	User_ID        string        `gorm:"column:user_id"`
	UserCart       []ProductUser `gorm:"column:foreignKey:UserID"`
	Addres_Details []Address     `gorm:"column:foreignKey:UserID"`
	Order_Status   []Order       `gorm:"column:foreignkey:UserID"`
}

type Product struct {
	product_ID   []byte  `gorm:"column:id;type:BLOB"`
	Product_Name *string `gorm:"column:product_name"`
	Price        int     `gorm:"column:price"`
	Rating       *uint   `gorm:"column:rating"`
	Image        *string `gormn"column:image"`
}

type ProductUser struct {
	Product_ID   []byte  `gorm:"column:id;type:BLOB"`
	Product_Name *string `gorm:"column:product_name"`
	Price        *string `gorm:"column:price"`
	Rating       *string `gorm:"column:rating"`
	Image        *string `gorm:"column:image"`
}

type Address struct {
	Address_ID []byte  `gorm:"column:id;type:BLOB"`
	House      *string `gormn:"column:house_name"`
	Street     *string `gorm:"column:street_name"`
	City       *string `gorm:"column:city_name"`
	Pincode    *string `gorm:"column:pin_code"`
}

type Order struct {
	Order_ID       string        `gorm:"column:id;type:BLOB"`
	Order_Cart     []ProductUser `gorm:"column:foreignKey:UserID"`
	Ordered_At     time.Time     `gorm:"column:ordered_at"`
	Price          int           `gorm:"column:price"`
	Discount       *int          `gorm:"column:discount"`
	Payment_Method bool          `gorm:"column:payment_method"`
}

type Payment struct {
	Digital       bool
	Created_Token bool
}
