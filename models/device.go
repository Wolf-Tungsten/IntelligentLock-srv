package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Device struct{
	Id primitive.ObjectID `bson:"_id,omitempty"`
	Uuid string `bson:"uuid" json:"uuid"`
	Admin primitive.ObjectID `bson:"admin"`
}

type Key struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	Uuid string `bson:"uuid" json:"uuid"`
	Key string `bson:"key" json:"uuid"`
}

type DeviceAccess struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	DeviceUuid string `bson:"device_uuid" json:"deviceUuid"`
	User primitive.ObjectID `bson:"user"`
	Allowed bool `bson:"allowed"`
}