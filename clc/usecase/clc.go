package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/fidellr/edu_malay/model"

	"github.com/fidellr/edu_malay/utils"

	"github.com/fidellr/edu_malay/clc"
	clcModel "github.com/fidellr/edu_malay/model/clc"
)

type profileUsecase struct {
	profileRepo    clc.ProfileRepository
	contextTimeout time.Duration
}

func NewClcProfileUsecase(cr clc.ProfileRepository, timeout time.Duration) clc.ProfileUsecase {
	return &profileUsecase{
		profileRepo:    cr,
		contextTimeout: timeout,
	}
}

func (u *profileUsecase) Create(c context.Context, clc *clcModel.ProfileEntity) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	clc.CreatedAt = time.Now()
	clc.UpdatedAt = time.Now()

	if err := utils.Validate(clc); err != nil {
		return err
	}

	err := validateClcLevelDataSupport(clc.ClcLevel, len(clc.ClcLevelDataSupport.StudentPerClass))
	if err != nil {
		return err
	}

	return u.profileRepo.Create(ctx, clc)
}

func (u *profileUsecase) FindAll(c context.Context, filter *model.Filter) ([]*clcModel.ProfileEntity, string, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	clc, cursor, err := u.profileRepo.FindAll(ctx, filter)
	if err != nil {
		return make([]*clcModel.ProfileEntity, 0), "", err
	}

	return clc, cursor, nil
}

func (u *profileUsecase) GetByID(c context.Context, id string) (*clcModel.ProfileEntity, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	clc, err := u.profileRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return clc, nil
}

func (u *profileUsecase) Update(c context.Context, id string, clc *clcModel.ProfileEntity) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	clc.UpdatedAt = time.Now()

	if err := utils.Validate(clc); err != nil {
		return err
	}

	return u.profileRepo.Update(ctx, id, clc)
}

func validateClcLevelDataSupport(level string, dataSupportSize int) error {
	var err error
	switch level {
	case "clc_sd":
		if dataSupportSize < 6 {
			err = fmt.Errorf("clc_sd must required 6 data set")
			break
		}
	case "clc_smp":
		if dataSupportSize < 3 {
			err = fmt.Errorf("clc_smp must required 3 data set")
			break
		}

	default:
		err = fmt.Errorf("clc_level must be one of clc_sd or clc_smp, other than that is not supported")
		break
	}

	return err
}