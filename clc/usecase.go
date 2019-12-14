package clc

import (
	"context"

	"github.com/fidellr/edu_malay/model"
	"github.com/fidellr/edu_malay/model/clc"
)

type ProfileUsecase interface {
	Create(ctx context.Context, t *clc.ProfileEntity) error
	FindAll(ctx context.Context, filter *model.Filter) ([]*clc.ProfileEntity, string, error)
	GetByID(ctx context.Context, id string) (*clc.ProfileEntity, error)
	Update(ctx context.Context, id string, t *clc.ProfileEntity) error
}
