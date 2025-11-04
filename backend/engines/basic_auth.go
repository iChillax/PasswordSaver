package engines

import (
	"backend/models"
	"backend/settings"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateBasicAuth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	basicAuthCollection := settings.MongoDatabase.Collection("BasicAuth")

	var new_basic_auth models.BasicAuth
	err := c.ShouldBindJSON(&new_basic_auth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = basicAuthCollection.InsertOne(ctx, new_basic_auth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.AsciiJSON(http.StatusOK, gin.H{
		"message": "Completed create new BasicAuth",
	})
}

func ListAllBasicAuth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	basicAuthCollection := settings.MongoDatabase.Collection("BasicAuth")
	cursor, err := basicAuthCollection.Find(ctx, bson.M{})
	if err != nil {
		c.AsciiJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer cursor.Close(ctx)
	var listOfBasicAuth []models.BasicAuth
	for cursor.Next(ctx) {
		var account models.BasicAuth
		if err := cursor.Decode(&account); err != nil {
			c.AsciiJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		listOfBasicAuth = append(listOfBasicAuth, account)
	}
	c.AsciiJSON(http.StatusOK, gin.H{
		"count": len(listOfBasicAuth),
		"data":  listOfBasicAuth})
}
