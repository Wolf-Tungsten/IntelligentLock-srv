package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Openid string `bson:"openid" json:"-"`
	Name string `bson:"name" json:"name"`
	PhoneNumber string `bson:"phoneNumber" json:"phoneNumber"`
	SessionToken string `bson:"sessionToken"`
}