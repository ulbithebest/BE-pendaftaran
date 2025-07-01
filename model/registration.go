package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Registration struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	Division   string             `bson:"division" json:"division"`
	Motivation string             `bson:"motivation" json:"motivation"`
	CVPath     string             `bson:"cv_path" json:"cv_path"`
	Status     string             `bson:"status" json:"status"` // lulus, tidak_lulus, menunggu
	Note       string             `bson:"note" json:"note"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}
