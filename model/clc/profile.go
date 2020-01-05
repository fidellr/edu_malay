package clc

import (
	"time"

	"github.com/fidellr/edu_malay/model/assembler"

	"github.com/globalsign/mgo/bson"
)

type ProfileEntity struct {
	ID                  bson.ObjectId               `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt           time.Time                   `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time                   `json:"updated_at" bson:"updated_at"`
	Name                string                      `json:"name" bson:"name" validate:"required"`
	ClcLevel            string                      `json:"clc_level" bson:"clc_level" validate:"oneof=clc_sd clc_smp"`
	ClcLevelDataSupport dataSupportDetail           `json:"clc_level_data_support" bson:"clc_level_data_support"`
	Status              string                      `json:"status" bson:"status" validate:"oneof=ladang non_ladang"`
	Gugus               string                      `json:"gugus" bson:"gugus" validate:"required,oneof=I II III IV VI VII VIII IX X XI XII XIII XIV sarawak"`
	Logo                string                      `json:"logo" bson:"logo"`
	Coordinate          coordinateDetail            `json:"coordinate" bson:"coordinate"`
	Teachers            []assembler.TeacherIdentity `json:"teachers" bson:"teachers"`
	Note                string                      `json:"note" bson:"note"`
	Vakum               bool                        `json:"vakum" bson:"vakum"`
	Permit              string                      `json:"permit" bson:"permit"`
}

type coordinateDetail struct {
	Long string `json:"long" bson:"long" validate:"required"`
	Lat  string `json:"lat" bson:"lat" validate:"required"`
}

type dataSupportDetail struct {
	TotalStudentPerClc int32         `json:"total_student_per_clc" bson:"total_student_per_clc" validate:"required"`
	StudentPerClass    []classDetail `json:"student_per_class" bson:"student_per_class"`
}

type classDetail struct {
	ClassLevel        int `json:"class_level" bson:"class_level" validate:"required"`
	TotalClassStudent int `json:"total_class_student" bson:"total_class_student" validate:"required"`
}
