package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Usuario struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id"`
	Name  string             `json:"name" bson:"name"`
	Email string             `json:"email" bson:"email"`
	Senha string             `json:"senha" bson:"senha"`
}
