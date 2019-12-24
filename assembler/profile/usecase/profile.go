package usecase

import (
	"context"
	"time"

	assembler "github.com/fidellr/edu_malay/assembler/profile"
	assemblerM "github.com/fidellr/edu_malay/model/assembler"
)

type profileAssemblerUsecase struct {
	profileAssemblerRepo assembler.ProfileAssemblerRepository
	contextTimeout       time.Duration
}

func NewProfileAssemblerUsecase(r assembler.ProfileAssemblerRepository, timeout time.Duration) assembler.ProfileAssemblerRepository {
	return &profileAssemblerUsecase{
		profileAssemblerRepo: r,
		contextTimeout:       timeout,
	}
}

func (u *profileAssemblerUsecase) Create(c context.Context, clcID string, m *assemblerM.ProfileAssemblerParam) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	return u.profileAssemblerRepo.Create(ctx, clcID, m)
}

func (u *profileAssemblerUsecase) FetchAll(c context.Context) ([]*assemblerM.ProfileAssemblerEntity, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.profileAssemblerRepo.FetchAll(ctx)
}

func (u *profileAssemblerUsecase) Remove(c context.Context, assmblrProfileID string) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.profileAssemblerRepo.Remove(ctx, assmblrProfileID)
}
