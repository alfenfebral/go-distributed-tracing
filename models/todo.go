package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Todo - todo model
type Todo struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	CreatedAt   time.Time     `json:"created_at" bson:"createdAt"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updatedAt"`
	DeletedAt   time.Time     `json:"deleted_at" bson:"deletedAt,omitempty"`
}
