package usecase

import (
	"sosmed/src/modules/profile/model"
)

//ProfileUsecase ...
type ProfileUsecase interface {
	SaveProfile(*model.Profile) (*model.Profile, error)
	UpdateProfile(string, *model.Profile) (*model.Profile, error)
	GetByID(string) (*model.Profile, error)
	GetByEmail(string) (*model.Profile, error)
}
