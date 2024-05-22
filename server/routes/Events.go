package routes

import (
	"context"
	"go-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validateEvent = validator.New()
var eventCollection *mongo.Collection = OpenCollection(Client, "Event")
var ticketCollection *mongo.Collection = OpenCollection(Client, "Ticket")

// CreateEvent creates a new event
func CreateEvent(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var event models.Event

	if err := c.BindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.ID = primitive.NewObjectID()
    for i := range event.Tickets {
        event.Tickets[i].ID = primitive.NewObjectID()
    }
	// check if the event is valid
	if err := validateEvent.Struct(event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	result, err := eventCollection.InsertOne(ctx, event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var ticketInterfaces []interface{}
	for _, ticket := range event.Tickets {
    	ticketInterfaces = append(ticketInterfaces, ticket)
	}

	res, err := ticketCollection.InsertMany(ctx, ticketInterfaces)

	defer cancel()
	c.JSON(http.StatusOK, gin.H{"data": result})
	c.JSON(http.StatusOK, gin.H{"data Ticket": res})
}

func GetEvents(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var events []models.Event

	cursor, err := eventCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var event models.Event
		cursor.Decode(&event)
		events = append(events, event)
	}

	defer cancel()
	c.JSON(http.StatusOK, gin.H{"data": events})
}

func DeleteEvent(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = eventCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cancel()
	c.JSON(http.StatusOK, gin.H{"data": "Event deleted"})
}

