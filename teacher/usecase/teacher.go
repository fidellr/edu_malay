package usecase

import (
	"context"
	"time"

	"github.com/fidellr/edu_malay/model"
	"github.com/fidellr/edu_malay/utils"

	teacherModel "github.com/fidellr/edu_malay/model/teacher"
	"github.com/fidellr/edu_malay/teacher"
)

type profileUsecase struct {
	profileRepo    teacher.ProfileRepository
	contextTimeout time.Duration
}

func NewTeacherProfileUsecase(tr teacher.ProfileRepository, timeout time.Duration) teacher.ProfileUsecase {
	return &profileUsecase{
		profileRepo:    tr,
		contextTimeout: timeout,
	}
}

func (u *profileUsecase) Create(c context.Context, t *teacherModel.ProfileEntity) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	if err := utils.Validate(t); err != nil {
		return err
	}

	return u.profileRepo.Create(ctx, t)
}

func (u *profileUsecase) FindAll(c context.Context, filter *model.Filter) ([]*teacherModel.ProfileEntity, string, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	t, cursor, err := u.profileRepo.FindAll(ctx, filter)
	if err != nil {
		return make([]*teacherModel.ProfileEntity, 0), "", err
	}

	return t, cursor, nil
}

func (u *profileUsecase) GetByID(c context.Context, id string) (*teacherModel.ProfileEntity, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	t, err := u.profileRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (u *profileUsecase) Update(c context.Context, id string, t *teacherModel.ProfileEntity) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	t.UpdatedAt = time.Now()

	if err := utils.Validate(t); err != nil {
		return err
	}

	return u.profileRepo.Update(ctx, id, t)
}

func (u *profileUsecase) Remove(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	return u.profileRepo.Remove(ctx, id)
}
