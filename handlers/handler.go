package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vanoraco/recipes-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type RecipesHandler struct {
	collection *mongo.Collection
	ctx context.Context
}

func NewRecipesHandler(ctx context.Context, collection *mongo.Collection) *RecipesHandler {
	return &RecipesHandler{
		collection: collection,
		ctx: ctx,
	}
}

func (handler *RecipesHandler) NewRecipeHandler(c *gin.Context) {
   var recipe models.Recipe

   if err := c.ShouldBindJSON(&recipe); err != nil {
	  c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	  return
   }

   recipe.ID = primitive.NewObjectID()
   recipe.PublishedAt = time.Now()

   _, err := handler.collection.InsertOne(handler.ctx, recipe)


  if err != nil {
	fmt.Println(err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting new recipe"})
	return
  }

   c.JSON(http.StatusOK, recipe)
}

func (handler *RecipesHandler) ListRecipesHandler (c *gin.Context) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cur.Close(handler.ctx)

	recipes := make([]models.Recipe, 0)

	for cur.Next(handler.ctx) {
		var recipe models.Recipe

		cur.Decode(&recipe)
		recipes = append(recipes, recipe)

	}
	c.JSON(http.StatusOK, recipes)
}

func (handler *RecipesHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	 var recipe models.Recipe

	 if err := c.ShouldBindJSON(&recipe) ; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	objectId, _ := primitive.ObjectIDFromHex(id)

	_, err := handler.collection.UpdateOne(handler.ctx, bson.M{"_id": objectId}, bson.D{{
		Key: "$set", Value: bson.D{
			{Key: "name", Value: recipe.Name},
			{Key: "instructions", Value: recipe.Instructions},
			{Key: "ingredients", Value: recipe.Ingredients},
			{Key: "tags", Value: recipe.Tags},
		},
	}})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

     c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
	/* index := -1

	for i := 0 ; i < len(recipes) ; i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	recipes[index] = recipe */

	

}

func (handler *RecipesHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	objectId, _ :=  primitive.ObjectIDFromHex(id) 

	_, err :=handler.collection.DeleteOne(handler.ctx, bson.M{"_id": objectId})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}


    c.JSON(http.StatusOK, gin.H{"message":"Recipe has been deleted"})
	/* index := -1

	for i := 0 ; i < len(recipes) ; i++ {
        if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	recipes = append(recipes[0: index], recipes[:index + 1]...) */
}

func (handler *RecipesHandler) SearchRecipesHandler(c *gin.Context) {
    id := c.Query("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

   curr := handler.collection.FindOne(handler.ctx, bson.M{"_id": objectId})

   var recipe models.Recipe

   err := curr.Decode(&recipe)

   if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	return
   }
  
  /*  listOfRecipes := make([]Recipe, 0)

   for i := 0; i < len(recipes) ; i++ {
	found := false

	for _,t := range recipes[i].Tags {
		if strings.EqualFold(t, tag) {
			found = true
		}
	}

	if found {
		listOfRecipes = append(listOfRecipes, recipes[i])
	}
   } */

   c.JSON(http.StatusOK, recipe)
}




