package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	ID			primitive.ObjectID	`bson:"_id"`
	Username	*string				`json:"username"`
	Password	*string				`json:"password"`
	Phone		*string				`json:"phone"`
	Email		*string				`json:"email"`
	Age			*int				`json:"age"`
	CreatedAt   time.Time			`json:"created_at"`
}