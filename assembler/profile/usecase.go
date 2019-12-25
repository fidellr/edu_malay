package assembler

import (
	"context"

	"github.com/fidellr/edu_malay/model/assembler"
)

type ProfileAssemblerUsecase interface {
	Create(ctx context.Context, clcID string, m *assembler.ProfileAssemblerParam) error
	FetchAll(ctx context.Context) ([]*assembler.ProfileAssemblerEntity, error)
	GetByID(ctx context.Context, id string) (*assembler.ProfileAssemblerEntity, error)
	Update(ctx context.Context, id string, teacherParam *assembler.ProfileAssemblerParam, isEditing bool) error
	Remove(ctx context.Context, assmblrID string) error
}
