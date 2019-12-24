package assembler

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type ProfileAssemblerParam struct {
	TeacherID     bson.ObjectId `json:"teacher_id" bson:"teacher_id" validate:"required"`
	StartWorkDate string        `json:"start_work_date" bson:"start_work_date" validate:"required"`
}

type TeacherIdentity struct {
	ID            bson.ObjectId `json:"id" bson:"_id"`
	FirstName     string        `json:"first_name" bson:"first_name"`
	LastName      string        `json:"last_name" bson:"last_name"`
	Gender        string        `json:"gender" bson:"gender" validate:"oneof=L P"`
	StartWorkDate string        `json:"start_work_date" bson:"start_work_date" validate:"required"`
}

type ProfileAssemblerEntity struct {
	ID        bson.ObjectId     `json:"id,omitempty" bson:"_id,omitempty" `
	CreatedAt time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" bson:"updated_at"`
	ClcID     bson.ObjectId     `json:"clc_id" bson:"clc_id" validate:"required"`
	Teachers  []TeacherIdentity `json:"teachers" bson:"teachers"`
}
