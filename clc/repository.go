package clc

import (
	"context"

	"github.com/fidellr/edu_malay/model/assembler"

	"github.com/fidellr/edu_malay/model"
	"github.com/fidellr/edu_malay/model/clc"
)

type ProfileRepository interface {
	Create(ctx context.Context, t *clc.ProfileEntity) error
	FindAll(ctx context.Context, filter *model.Filter) ([]*clc.ProfileEntity, string, error)
	GetByID(ctx context.Context, id string) (*clc.ProfileEntity, error)
	Update(ctx context.Context, id string, t *clc.ProfileEntity) error
	AssembleProfile(ctx context.Context, clcID string, teacherID string, startDate string) error
	UpdateAssembledProfile(ctx context.Context, clcID string, m *assembler.TeacherIdentity, isEditing bool) error
	Remove(ctx context.Context, id string) error
}
