package database

import (
	"context"
	"errors"
	"github.com/AshrafulHaqueToni/xyz-cart-with-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func AddProductToCart(ctx context.Context, productCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchFromMongoDB, err := productCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return errors.New("can't find the product")
	}
	var productCard []models.ProductSelectByUser
	err = searchFromMongoDB.All(ctx, &productCard)
	if err != nil {
		log.Println(err)
		return errors.New("can't decode the product")
	}
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return errors.New("userId is invalid")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCard}}}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return errors.New("can't update the product by user")
	}
	return nil
}

func RemoveCartItem(ctx context.Context, productCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return errors.New("userId is invalid")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return errors.New("can't remove from cart")
	}
	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return errors.New("userId is invalid")
	}
	var cartItems models.User
	var orderCart models.Order
	orderCart.OrderID = primitive.NewObjectID()
	orderCart.OrderedAt = time.Now()
	orderCart.OrderCart = make([]models.ProductSelectByUser, 0)
	orderCart.PaymentMethod.COD = true
	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
	currentResults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}
	var getUserCart []bson.M
	if err := currentResults.All(ctx, &getUserCart); err != nil {
		panic(err)
	}
	var totalPrice int32
	for _, userItem := range getUserCart {
		price := userItem["total"]
		totalPrice = price.(int32)
	}
	orderCart.Price = int(totalPrice)
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderCart}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&cartItems)
	if err != nil {
		log.Println(err)
	}
	filter = bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": cartItems.UserCart}}}
	_, err = userCollection.UpdateOne(ctx, filter, update2)
	if err != nil {
		log.Println(err)
	}
	emptyUsercart := make([]models.ProductSelectByUser, 0)
	filtered := bson.D{primitive.E{Key: "_id", Value: id}}
	updated := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: emptyUsercart}}}}
	_, err = userCollection.UpdateOne(ctx, filtered, updated)
	if err != nil {
		return errors.New("can't update the purchase")
	}
	return nil

}

func InstantBuyer(ctx context.Context, productCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return errors.New("userId is invalid")
	}
	var productDetails models.ProductSelectByUser
	var orderDetails models.Order
	orderDetails.OrderID = primitive.NewObjectID()
	orderDetails.OrderedAt = time.Now()
	orderDetails.OrderCart = make([]models.ProductSelectByUser, 0)
	orderDetails.PaymentMethod.COD = true
	err = productCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&productDetails)
	if err != nil {
		log.Println(err)
	}
	orderDetails.Price = productDetails.Price
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderDetails}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	filter = bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": productDetails}}
	_, err = userCollection.UpdateOne(ctx, filter, update2)
	if err != nil {
		log.Println(err)
	}
	return nil
}
