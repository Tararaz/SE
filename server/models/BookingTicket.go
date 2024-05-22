package models

import (

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingTicket struct{
	TicketID	primitive.ObjectID	`json:"ticket_id"`
	Quantity 	*int				`json:"quantity"`
}
