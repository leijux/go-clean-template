package v1

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/leijux/go-clean-template/internal/controller/restapi/v1/request"
	"github.com/leijux/go-clean-template/internal/controller/restapi/v1/response"
	"github.com/leijux/go-clean-template/internal/entity"
)

// @Summary     Register
// @Description Register a new user
// @ID          register
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body     request.Register true "Registration data"
// @Success     201     {object} response.Envelope{Code=int,Message=string,Data=entity.User}
// @Failure     400     {object} response.Envelope
// @Failure     409     {object} response.Envelope
// @Failure     500     {object} response.Envelope
// @Router      /auth/register [post]
func (r *V1) register(ctx fiber.Ctx) error {
	var body request.Register

	if err := ctx.Bind().Body(&body); err != nil {
		r.l.Error(err, "restapi - v1 - register")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - register")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	user, err := r.u.Register(ctx, body.Username, body.Email, body.Password)
	if err != nil {
		r.l.Error(err, "restapi - v1 - register")

		if errors.Is(err, entity.ErrUserAlreadyExists) {
			return errorResponse(ctx, http.StatusConflict, "user already exists")
		}

		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return okResponse(ctx, http.StatusCreated, user)
}

// @Summary     Login
// @Description Authenticate user and get JWT token
// @ID          login
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body     request.Login true "Login credentials"
// @Success     200     {object} response.Envelope{Code=int,Message=string,Data=response.Token}
// @Failure     400     {object} response.Envelope
// @Failure     401     {object} response.Envelope
// @Failure     500     {object} response.Envelope
// @Router      /auth/login [post]
func (r *V1) login(ctx fiber.Ctx) error {
	var body request.Login

	if err := ctx.Bind().Body(&body); err != nil {
		r.l.Error(err, "restapi - v1 - login")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - login")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	token, err := r.u.Login(ctx, body.Email, body.Password)
	if err != nil {
		r.l.Error(err, "restapi - v1 - login")

		if errors.Is(err, entity.ErrInvalidCredentials) {
			return errorResponse(ctx, http.StatusUnauthorized, "invalid credentials")
		}

		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return okResponse(ctx, http.StatusOK, response.Token{Token: token})
}

// @Summary     Get profile
// @Description Get current user profile
// @ID          profile
// @Tags        user
// @Produce     json
// @Success     200 {object} response.Envelope{Code=int,Message=string,Data=entity.User}
// @Failure     401 {object} response.Envelope
// @Failure     404 {object} response.Envelope
// @Failure     500 {object} response.Envelope
// @Security    BearerAuth
// @Router      /user/profile [get]
func (r *V1) profile(ctx fiber.Ctx) error {
	userID, ok := ctx.Locals("userID").(string)
	if !ok {
		return errorResponse(ctx, http.StatusUnauthorized, "unauthorized")
	}

	user, err := r.u.GetUser(ctx, userID)
	if err != nil {
		r.l.Error(err, "restapi - v1 - profile")

		if errors.Is(err, entity.ErrUserNotFound) {
			return errorResponse(ctx, http.StatusNotFound, "user not found")
		}

		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return okResponse(ctx, http.StatusOK, user)
}
