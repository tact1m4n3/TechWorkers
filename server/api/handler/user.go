package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"prohiking-server/api/auth"
	my_middleware "prohiking-server/api/handler/middleware"
	"prohiking-server/api/response"
	"prohiking-server/internal/database"
	"prohiking-server/internal/model"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

const TokenExpirationTime = 2 * time.Hour

func GetUserById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "user id not valid")
		return
	}

	user, err := database.GetUserById(id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			response.Error(w, http.StatusNotFound, fmt.Sprintf(
				"no user found with id %v", id,
			))
		} else {
			response.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.JSON(w, http.StatusOK, user)
}

type loginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (lp *loginPayload) validate() bool {
	if len(lp.Username) <= 0 || len(lp.Password) <= 0 {
		return false
	}
	return true
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	payload := &loginPayload{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	if !payload.validate() {
		response.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	user, err := database.GetUserByName(payload.Username)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			response.Error(w, http.StatusNotFound, fmt.Sprintf(
				"no user found with name '%v'", payload.Username,
			))
		} else {
			response.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if !auth.VerifyPassword(payload.Password, user.Password) {
		response.Error(w, http.StatusBadRequest, "wrong password")
		return
	}

	token, err := auth.GenerateJWT(user.Id, TokenExpirationTime)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Value:   token,
		MaxAge:  0,
		Path:    "/",
		Expires: time.Now().Add(TokenExpirationTime),
	})

	response.JSON(w, http.StatusAccepted, user)
}

type registerPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (rp *registerPayload) validate() bool {
	if len(rp.Username) <= 0 || len(rp.Email) <= 0 || len(rp.Password) <= 8 {
		return false
	}
	if _, err := mail.ParseAddress(rp.Email); err != nil {
		return false
	}
	return true
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	payload := &registerPayload{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	if !payload.validate() {
		response.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	user := &model.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashedPassword,
	}

	if err := database.CreateUser(user); err != nil {
		if errors.Is(err, database.ErrDuplicatedKey) {
			response.Error(w, http.StatusBadRequest, "username or email already exist")
		} else {
			response.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.JSON(w, http.StatusCreated, user)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Value:   "",
		MaxAge:  0,
		Path:    "/",
		Expires: time.Now().Add(-TokenExpirationTime),
	})

	w.Write([]byte("logged out"))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(my_middleware.AuthUserKey).(*model.User)
	if err := database.DeleteUser(user.Id); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Write([]byte("deleted user"))
}
