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

var validateOrder = validator.New()
var orderCollection *mongo.Collection = OpenCollection(Client, "Booking")
var bookingTicketCollection *mongo.Collection = OpenCollection(Client, "BookingTicket")

func CreateOrder(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var order models.Booking

	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.ID = primitive.NewObjectID()

	result, err := orderCollection.InsertOne(ctx, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var ticketInterfaces []interface{}
	for _, ticket := range order.Tickets {
    	ticketInterfaces = append(ticketInterfaces, ticket)
	}

	res, err := bookingTicketCollection.InsertMany(ctx, ticketInterfaces)

	defer cancel()
	c.JSON(http.StatusOK, gin.H{"data": result})
	c.JSON(http.StatusOK, gin.H{"data Ticket": res})
}

func CalculateTotalPrice(order models.Booking) float64 {
    totalPrice := 0.0
    for _, ticketBooking := range order.Tickets {
        ticket := GetTicketByID(ticketBooking.TicketID)
        if ticket.Price != nil && ticketBooking.Quantity != nil {
            totalPrice += *ticket.Price * float64(*ticketBooking.Quantity)
        }
    }
    return totalPrice
}

func GetTicketByID(id primitive.ObjectID) models.Ticket {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var ticket models.Ticket

	err := ticketCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&ticket)
	if err != nil {
		return models.Ticket{}
	}

	defer cancel()
	return ticket
}

func GetOrders(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var orders []models.Booking

	cursor, err := orderCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var order models.Booking
		cursor.Decode(&order)
		orders = append(orders, order)
	}
	for i := range orders{
		var totalPrice = CalculateTotalPrice(orders[i]) 
		c.JSON(http.StatusOK, gin.H{"Total Price": totalPrice})
	}

	var totalPrice = CalculateTotalPrice(orders[2])
	defer cancel()
	c.JSON(http.StatusOK, gin.H{"data": orders})
	c.JSON(http.StatusOK, gin.H{"len": len(orders)})
	c.JSON(http.StatusOK, gin.H{"testing": orders[2].Tickets})
	c.JSON(http.StatusOK, gin.H{"Total Price": totalPrice})
}

func DeleteOrder(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := orderCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cancel()
	c.JSON(http.StatusOK, gin.H{"data": result})
}