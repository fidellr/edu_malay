package teacher

import (
	"context"

	"github.com/fidellr/edu_malay/model/teacher"
)

type ProfileRepository interface {
	Create(ctx context.Context, t *teacher.ProfileEntity) error
	FindAll(ctx context.Context, filter *teacher.Filter) ([]*teacher.ProfileEntity, string, error)
	GetByID(ctx context.Context, id string) (*teacher.ProfileEntity, error)
	Update(ctx context.Context, id string, t *teacher.ProfileEntity) error
}
