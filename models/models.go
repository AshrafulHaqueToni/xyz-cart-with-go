package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID             primitive.ObjectID    `json:"_id" bson:"_id"`
	FirstName      *string               `json:"firstName" validate:"required, min=2, max=30"`
	LastName       *string               `json:"lastName" validate:"required, min=2, max=30"`
	Password       *string               `json:"password" validate:"required, min=6"`
	Email          *string               `json:"email" validate:"email, required"`
	Phone          *string               `json:"phone" validate:"required"`
	Token          *string               `json:"token"`
	RefreshToken   *string               `json:"refresh_token"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	UserID         string                `json:"user_id"`
	UserCart       []ProductSelectByUser `json:"user_cart" bson:"user_cart"`
	AddressDetails []Address             `json:"address_details" bson:"address_details"`
	OrderStatus    []Order               `json:"order_status" bson:"order_status"`
}

type Product struct {
	ProductID   primitive.ObjectID `json:"_id" bson:"_id"`
	ProductName *string            `json:"product_name"`
	Price       *uint64            `json:"price"`
	Rating      *uint8             `json:"rating"`
	// image url
	Image *string `json:"image"`
}

type ProductSelectByUser struct {
	ProductID   primitive.ObjectID `json:"_id" bson:"_id"`
	ProductName *string            `json:"product_name" bson:"product_name"`
	Price       *uint64            `json:"price" bson:"price"`
	Rating      *uint8             `json:"rating" bson:"rating"`
	Image       *string            `json:"image" bson:"image"`
}

type Address struct {
	AddressID primitive.ObjectID `bson:"_id"`
	House     *string            `json:"house" bson:"house"`
	Street    *string            `json:"street" bson:"street"`
	City      *string            `json:"city" bson:"city"`
	PinCode   *string            `json:"pin_code" bson:"pin_code"`
}

type Order struct {
	OrderID       primitive.ObjectID    `bson:"_id"`
	OrderCart     []ProductSelectByUser `json:"order_cart" bson:"order_cart"`
	OrderedAt     time.Time             `json:"ordered_at" bson:"ordered_at"`
	Price         int                   `json:"price" bson:"price"`
	Discount      *int                  `json:"discount" bson:"discount"`
	PaymentMethod Payment               `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	Digital bool `json:"digital" bson:"digital"`
	COD     bool `json:"cod" bson:"cod"`
}
