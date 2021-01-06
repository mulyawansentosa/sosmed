package repository

import (
	"sosmed/src/modules/profile/model"
)

//ProfileRepository ...
type ProfileRepository interface {
	Save(*model.Profile) error
	Update(string, *model.Profile) error
	Delete(string) error
	FindByID(string) (*model.Profile, error)
	FindByEmail(string) (*model.Profile, error)
	FindAll() (model.Profiles, error)
}
