package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookAuthor struct {
	Firstname string 	`json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname string 	`json:"lastname,omitempty" bson:"lastname,omitempty"`
}

type Book struct {
	Id    primitive.ObjectID 	 `json:"_id,omitempty" bson:"_id,omitempty"`
	Title string             	 `json:"title,omitempty" bson:"title,omitempty"`
	Description  string          `json:"description,omitempty" bson:"description,omitempty"`
	Author BookAuthor			 `json:"author,omitempty" bson:"author,omitempty"`
	ReleaseYear int32 			 `json:"release_year,omitempty" bson:"release_year,omitempty"`
}