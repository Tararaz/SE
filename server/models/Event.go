package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct{
	ID			primitive.ObjectID	`bson:"_id"`
	Title		*string				`json:"title"`
	StartDate  time.Time			`json:"start_date"`
	EndDate    time.Time			`json:"end_date"`
	ImageURL	*string				`json:"image_url"`
	Description	*string				`json:"description"`
	Location	*string				`json:"location"`
	AvailableTicket *int			`json:"available_ticket"`
	Tickets 	[]Ticket			`json:"tickets"`
	PublisherName *string			`json:"publisher_name"`
}
