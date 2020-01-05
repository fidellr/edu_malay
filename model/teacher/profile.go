package teacher

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// ProfileEntity : teacher's entity model
type ProfileEntity struct {
	ID               bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt        time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at" bson:"updated_at"`
	FirstName        string        `json:"first_name" bson:"first_name" validate:"required"`
	LastName         string        `json:"last_name" bson:"last_name" validate:"required"`
	POB              string        `json:"place_of_birth" bson:"place_of_birth" validate:"required"`
	DOB              string        `json:"date_of_birth" bson:"date_of_birth" validate:"required"`
	Gender           string        `json:"gender" bson:"gender" validate:"required"`
	Religion         string        `json:"religion" bson:"religion" validate:"required"`
	University       string        `json:"university" bson:"university" validate:"required"`
	Major            string        `json:"major" bson:"major" validate:"required"`
	YearOfDedication string        `json:"year_of_dedication" bson:"year_of_dedication" validate:"required"`
	StartWorkDate    string        `json:"start_work_date" bson:"start_work_date"`
	MapelUtama       string        `json:"mapel_utama" bson:"mapel_utama"`
	IsAssembled      bool          `json:"is_assembled" bson:"is_assembled"`
}

// ProfileHardDeleteQueue : teacher's archived entity model
type ProfileHardDeleteQueue struct {
	ProfileEntity
	ApproveBy string `json:"approve_by" bson:"approve_by" validate:"required"`
	IsApprove bool   `json:"is_approve" bson:"is_approve" validate:"required"`
}
