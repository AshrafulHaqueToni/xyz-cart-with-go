package main

import (
	"github.com/AshrafulHaqueToni/xyz-cart-with-go/controllers"
	"github.com/AshrafulHaqueToni/xyz-cart-with-go/database"
	"github.com/AshrafulHaqueToni/xyz-cart-with-go/middleware"
	"github.com/AshrafulHaqueToni/xyz-cart-with-go/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app := controllers.NewApplication(database.ProductData(database.DBSet(), "Products"), database.UserData(database.DBSet(), "Users"))
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.ByFromCart())
	router.GET("/instanbuy", app.InstantBuy())
	router.GET("/listcart", controllers.GetItemFromCart())
	router.POST("/addaddress", controllers.AddAddress())
	router.PUT("/edithomeaddress", controllers.EditHomeAddress())
	router.PUT("/editworkaddress", controllers.EditWorkAddress())
	router.GET("/deleteaddresses", controllers.DeleteAddress())
	log.Fatal(router.Run(":" + port))
}
