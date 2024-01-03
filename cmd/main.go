package cmd

import (
	//"github.com/AshrafulHaqueToni/xyz-cart-with-go/controllers"
	//"github.com/AshrafulHaqueToni/xyz-cart-with-go/database"
	//"github.com/AshrafulHaqueToni/xyz-cart-with-go/middleware"
	//"github.com/AshrafulHaqueToni/xyz-cart-with-go/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.ByFromCart())
	router.GET("/instanbuy", app.InstantBuy())
	log.Fatal(router.Run(":" + port))model
}
