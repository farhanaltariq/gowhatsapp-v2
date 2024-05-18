package controllers

import (
	"encoding/json"

	"github.com/farhanaltariq/fiberplate/common/codes"
	"github.com/farhanaltariq/fiberplate/common/status"
	"github.com/farhanaltariq/fiberplate/common/usertype"
	"github.com/farhanaltariq/fiberplate/database/models"
	"github.com/farhanaltariq/fiberplate/middleware"
	"github.com/farhanaltariq/fiberplate/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthenticationController interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

func NewAuthController(service middleware.Services) AuthenticationController {
	return &controller{service}
}

// @Summary Register
// @Description Register a new user
// @Tags Authentication
// @Accept json
// @Param data body models.Register true "Register data"
// @Produce json
// @Success 200 {object} common.ResponseMessage
// @Failure 400 {object} common.ResponseMessage
// @Router /auth/register [post]
func (s *controller) Register(c *fiber.Ctx) error {
	auth := models.Register{}
	if err := json.Unmarshal(c.Body(), &auth); err != nil {
		return err
	}

	if auth.Password != auth.ConfirmPassword {
		return status.Errorf(c, codes.BadRequest, "Password and confirm password does not match")
	}

	authOrm := &models.Authentications{}
	if err := json.Unmarshal(c.Body(), &authOrm); err != nil {
		logrus.Errorln("Failed to unmarshal auth ORM", err)
		return status.Errorf(c, codes.BadRequest, err.Error())
	}
	userOrm := models.User{}
	if err := json.Unmarshal(c.Body(), &userOrm); err != nil {
		logrus.Errorln("Failed to unmarshal user ORM", err)
		return status.Errorf(c, codes.BadRequest, err.Error())
	}

	data, err := s.UserService.GetDataByUsernameOrEmail(userOrm)
	if err != nil || data.ID != 0 {
		return status.Errorf(c, codes.BadRequest, "Username or email already used")
	}

	// insert to user tables
	userData := models.User{}
	userOrm.UserType = usertype.ADMIN
	if userData, err = s.UserService.InsertOrUpdate(userOrm); err != nil {
		logrus.Errorln("Failed to register user", err)
		return status.Errorf(c, codes.BadRequest, err.Error())
	}

	pass, salt := utils.Encrypt(auth.Password)
	authOrm.Password = pass
	authOrm.Salt = salt
	authOrm.UserId = userData.ID

	if err := s.AuthService.InsertOrUpdate(*authOrm); err != nil {
		logrus.Errorln("Failed to register user", err)
		return status.Errorf(c, codes.BadRequest, err.Error())
	}

	logrus.Infoln("User", userOrm)

	return status.Successf(c, codes.OK, "Success")
}

// @Summary Login
// @Description Login
// @Tags Authentication
// @Security Authorization
// @Accept json
// @Param data body models.Login true "Login Data"
// @Produce json
// @Success 200 {object} models.AuthenticationResponse
// @Failure 400 {object} common.ResponseMessage
// @Router /auth/login [post]
func (s *controller) Login(c *fiber.Ctx) error {
	cred := models.Login{}

	if err := json.Unmarshal(c.Body(), &cred); err != nil {
		return status.Errorf(c, codes.BadRequest, err.Error())
	}

	if cred.UsernameOrEmail == "" {
		return status.Errorf(c, codes.BadRequest, "Username or email is required")
	}

	data, err := s.UserService.GetDataByUsernameOrEmail(models.User{Username: cred.UsernameOrEmail, Email: cred.UsernameOrEmail})
	if err != nil || data == (models.User{}) {
		return status.Errorf(c, codes.BadRequest, "Username or email not found")
	}

	authData, err := s.AuthService.GetDataByUserId(data.ID)
	if err != nil {
		return status.Errorf(c, codes.BadRequest, err.Error())
	}

	pass, err := utils.Decrypt(authData.Password, authData.Salt)
	if err != nil {
		return status.Errorf(c, codes.BadRequest, err.Error())
	}
	if pass != cred.Password {
		return status.Errorf(c, codes.BadRequest, "Invalid credentials")
	}

	token, err := utils.GenerateToken(&data)
	if err != nil {
		return status.Errorf(c, codes.BadRequest, err.Error())
	}

	res := models.AuthenticationResponse{
		Status:      "Success",
		AccessToken: token,
		ExpiredAt:   utils.GetExpirationTime().Format("2006-01-02 15:04:05"),
	}

	return c.Status(codes.OK).JSON(res)
}
