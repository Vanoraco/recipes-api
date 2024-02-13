// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin.
//
// Schemes: http

// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
// Contact: Seare Hagos
// <seareha369@lgmail.com>
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	handlers "github.com/vanoraco/recipes-api/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/go-redis/redis/v8"
)

var ctx context.Context
var err error
var client *mongo.Client
var recipesHandler *handlers.RecipesHandler

func init() {

	/* recipes = make([]Recipe, 0)

	file, _ := os.ReadFile("recipes.json")

	_ = json.Unmarshal([]byte(file), &recipes) */

    ctx = context.Background()


	client, err = mongo.Connect(ctx, 
	options.Client().ApplyURI(os.Getenv("MONGO_URI")))
      if err != nil {
       log.Fatal(err)
       }

	 if err = client.Ping(context.TODO(),
          readpref.Primary()); err != nil {
			log.Fatal(err)
		  }

	log.Println("Connected to MongoDB")

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	
	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)

	

	status := redisClient.Ping(context.Background())

	fmt.Println(status)
	/* var listOfRecipes []interface{}

	for _, recipe := range recipes {
		listOfRecipes = append(listOfRecipes, recipe)
	}

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserted Recipes: ", len(insertManyResult.InsertedIDs)) */

}



func main() {
	
	 router := gin.Default()
     router.POST("/recipes", recipesHandler.NewRecipeHandler)
	 router.GET("/recipes", recipesHandler.ListRecipesHandler)
	 router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	 router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	 router.GET("/recipes/search", recipesHandler.SearchRecipesHandler)
	 router.Run()
}