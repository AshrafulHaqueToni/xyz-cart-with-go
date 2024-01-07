package controllers

import (
	"context"
	"fmt"
	"github.com/AshrafulHaqueToni/xyz-cart-with-go/database"
	"github.com/AshrafulHaqueToni/xyz-cart-with-go/models"
	"github.com/AshrafulHaqueToni/xyz-cart-with-go/tokens"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

const (
	Users    = "Users"
	Products = "Products"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func VerifyPassword(usersPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(usersPassword), []byte(givenPassword))
	if err != nil {
		return false, "Password is Incorrect"
	}
	return true, ""
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validator.New().Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		count, err := database.UserData(database.DBSet(), Users).CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Panic(err)
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		}
		count, err = database.UserData(database.DBSet(), Users).CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Panic(err)
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is already in use"})
		}
		password := HashPassword(*user.Password)
		user.Password = &password
		user.CreatedAt, _ = time.Parse(time.RFC822, time.Now().Format(time.RFC822))
		user.UpdatedAt, _ = time.Parse(time.RFC822, time.Now().Format(time.RFC822))
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()
		token, refreshToken, _ := tokens.GenerateToken(*user.Email, *user.FirstName, *user.LastName, user.UserID)
		user.Token = &token
		user.RefreshToken = &refreshToken
		user.UserCart = make([]models.ProductSelectByUser, 0)
		user.AddressDetails = make([]models.Address, 0)
		user.OrderStatus = make([]models.Order, 0)
		_, err = database.UserData(database.DBSet(), Users).InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, "Successfully Signed Up!!")
	}
}

func LogIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()
		var user models.User
		var foundUser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := database.UserData(database.DBSet(), Users).FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
			return
		}
		PasswordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}
		token, refreshToken, _ := tokens.GenerateToken(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.UserID)
		defer cancel()
		tokens.UpdateAllTokens(token, refreshToken, foundUser.UserID)
		c.JSON(http.StatusFound, foundUser)
	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		product.ProductID = primitive.NewObjectID()
		_, err := database.ProductData(database.DBSet(), Products).InsertOne(ctx, product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Successfully added Product by Admin!!")
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()
		cursor, err := database.ProductData(database.DBSet(), Products).Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Something Went Wrong Please Try After Some Time")
			return
		}
		err = cursor.All(ctx, &productList)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, "invalid request")
		}
		defer cancel()
		c.IndentedJSON(http.StatusOK, productList)
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchProduct []models.Product
		queryParameter := c.Query("name")
		if queryParameter == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()
		res, err := database.ProductData(database.DBSet(), Products).Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParameter}})
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, "something went wrong in fetching the db query")
			return
		}

		err = res.All(ctx, &searchProduct)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer res.Close(ctx)
		if err := res.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, "invalid request")
		}
		defer cancel()
		c.IndentedJSON(http.StatusOK, searchProduct)

	}
}
