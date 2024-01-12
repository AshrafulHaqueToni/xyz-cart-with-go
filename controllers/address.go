package controllers

import (
	"context"
	"fmt"
	"github.com/AshrafulHaqueToni/xyz-cart-with-go/database"
	"github.com/AshrafulHaqueToni/xyz-cart-with-go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		var addresses models.Address
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		address, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
		}
		addresses.AddressID = primitive.NewObjectID()
		if err := c.BindJSON(addresses); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()
		matchFilter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}
		pointCursor, err := database.UserData(database.DBSet(), Users).Aggregate(ctx, mongo.Pipeline{matchFilter, unwind, grouping})
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
		}
		var addressInfo []bson.M
		if err = pointCursor.All(ctx, &addressInfo); err != nil {
			panic(err)
		}
		var size int32
		for _, noOfAddress := range addressInfo {
			count := noOfAddress["count"]
			size = count.(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
			_, err := database.UserData(database.DBSet(), Users).UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			c.IndentedJSON(http.StatusBadRequest, "Not Allowed ")
		}
		defer cancel()
		ctx.Done()

	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		var editAddress models.Address
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		userTransID, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
		}
		if err := c.BindJSON(&editAddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: userTransID}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editAddress.House}, {Key: "address.0.street_name", Value: editAddress.Street}, {Key: "address.0.city_name", Value: editAddress.City}, {Key: "address.0.pin_code", Value: editAddress.PinCode}}}}
		_, err = database.UserData(database.DBSet(), Users).UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "something Went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "Successfully updated the Work Address")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		var editAddress models.Address
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		userTransID, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
		}
		if err := c.BindJSON(&editAddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: userTransID}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editAddress.House}, {Key: "address.1.street_name", Value: editAddress.Street}, {Key: "address.1.city_name", Value: editAddress.City}, {Key: "address.1.pin_code", Value: editAddress.PinCode}}}}
		_, err = database.UserData(database.DBSet(), Users).UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "something Went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "Successfully updated the Work Address")
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		addresses := make([]models.Address, 0)
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		userTransId, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: userTransId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{
			Key: "address", Value: addresses,
		}}}}
		_, err = database.UserData(database.DBSet(), Users).UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, "Wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "Successfully Deleted!")

	}
}
