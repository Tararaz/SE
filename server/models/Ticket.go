package models

import (

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ticket struct{
	ID			primitive.ObjectID	`bson:"_id"`
	Category 	*string				`json:"category"`
	Price 		*float64			`json:"price"`
	Quantity 	*int				`json:"quantity"`
}
