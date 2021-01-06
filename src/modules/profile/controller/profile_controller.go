package controller

import (
	"fmt"
	"io"
	"os"
	"sosmed/src/modules/profile/model"
	"sosmed/src/modules/profile/usecase"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	uuid "github.com/satori/go.uuid"
	"github.com/wuriyanto48/replacer"
)

//ProfileIDKey ...
const ProfileIDKey = "ProfileID"

//ProfileController ...
type ProfileController struct {
	Ctx iris.Context

	Session *sessions.Session

	ProfileUsecase usecase.ProfileUsecase
}

func (c *ProfileController) getCurrentProfileID() string {
	return c.Session.GetString(ProfileIDKey)
}

func (c *ProfileController) isProfileLoggedIn() bool {
	return c.getCurrentProfileID() != ""
}

func (c *ProfileController) logout() {
	c.Session.Destroy()
}

//GetRegister ...
func (c *ProfileController) GetRegister() mvc.Result {
	if c.isProfileLoggedIn() {
		c.logout()
	}
	return mvc.View{
		Name: "profile/register.html",
		Data: iris.Map{"Title": "Profile Registration"},
	}
}

//PostRegister ...
func (c *ProfileController) PostRegister() mvc.Result {
	firstName := c.Ctx.FormValue("first_name")
	lastName := c.Ctx.FormValue("last_name")
	email := c.Ctx.FormValue("email")
	password := c.Ctx.FormValue("password")

	if firstName == "" || lastName == "" || email == "" || password == "" {
		return mvc.Response{
			Path: "/profile/register",
		}
	}
	id := uuid.NewV4()

	profileImage, err := c.uploadImage(c.Ctx, id.String())

	if err != nil {
		return mvc.Response{
			Path: "/profile/register",
		}
	}

	var profile model.Profile
	profile.ID = id.String()
	profile.FirstName = firstName
	profile.LastName = lastName
	profile.Email = email
	profile.Password = password
	profile.ImageProfile = profileImage
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	_, err = c.ProfileUsecase.SaveProfile(&profile)
	if err != err {
		return mvc.Response{
			Path: "/profile/register",
		}
	}

	c.Session.Set(ProfileIDKey, profile.ID)
	return mvc.Response{
		Path: "/profile/me",
	}
}

//GetLogin ...
func (c *ProfileController) GetLogin() mvc.Result {
	if c.isProfileLoggedIn() {
		c.logout()
	}
	return mvc.View{
		Name: "/profile/login.html",
		Data: iris.Map{"Title": "Login"},
	}
}

//PostLogin ...
func (c *ProfileController) PostLogin() mvc.Result {
	email := c.Ctx.FormValue("email")
	password := c.Ctx.FormValue("password")

	if email == "" || password == "" {
		return mvc.Response{
			Path: "/profile/login",
		}
	}

	profile, err := c.ProfileUsecase.GetByEmail(email)
	// profile, err := c.ProfileUsecase.GetByEmail(email)

	if err != nil {
		return mvc.Response{
			Path: "/profile/login",
		}
	}

	if !profile.IsValidPassword(password) {
		return mvc.Response{
			Path: "/profile/login",
		}
	}

	c.Session.Set(ProfileIDKey, profile.ID)

	return mvc.Response{
		Path: "/profile/me",
	}
}

//GetMe ...
func (c *ProfileController) GetMe() mvc.Result {
	if !c.isProfileLoggedIn() {
		return mvc.Response{
			Path: "/profile/login",
		}
	}

	profile, err := c.ProfileUsecase.GetByID(c.getCurrentProfileID())

	if err != nil {
		c.logout()
		c.GetMe()
	}

	return mvc.View{
		Name: "profile/me.html",
		Data: iris.Map{"Title": "My profile", "Profile": profile},
	}
}

//AnyLogout ...
func (c *ProfileController) AnyLogout() {
	if c.isProfileLoggedIn() {
		c.logout()
	}

	c.Ctx.Redirect("/profile/login")
}

func (c *ProfileController) uploadImage(ctx iris.Context, id string) (string, error) {
	file, info, err := ctx.FormFile("image_profile")

	if err != nil {
		return "", err
	}

	defer file.Close()

	fileName := fmt.Sprintf("%s%s%s", id, "_", replacer.Replace(info.Filename, "_"))
	out, err := os.OpenFile("./web/public/images/profile/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		return "", err
	}

	defer out.Close()

	io.Copy(out, file)
	return fileName, nil
}
