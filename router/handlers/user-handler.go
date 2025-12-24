package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"template/data/dtos"
	"template/data/entities"
	"template/services"
	"template/utils"
	"template/utils/auth"

	"github.com/golang-jwt/jwt/v5"
)

type IUserHandler interface {
	UserHandler(http.ResponseWriter, *http.Request)
	UserHandlerById(http.ResponseWriter, *http.Request)
}

type UserHandlerStruct struct {
	userSvc services.IUserService
	credSvc services.ICredentialService
}

func UserHandlerInit(userService services.IUserService, credentialService services.ICredentialService) *UserHandlerStruct {
	return &UserHandlerStruct{userSvc: userService, credSvc: credentialService}
}

func (u *UserHandlerStruct) UserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method == http.MethodGet {
		users, err := u.userSvc.GetAllUsers(r.Context())
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, fmt.Errorf("error in get %v", err))
			return
		}

		utils.WriteJSON(w, http.StatusOK, users)
	} else if r.Method == http.MethodPost {
		var user entities.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, fmt.Errorf("error in parsing %v", err))
			return
		}

		err := u.userSvc.CreateUser(r.Context(), user)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("error in creation %v", err))
			return
		}

		utils.WriteJSON(w, http.StatusCreated, nil)
	} else {
		utils.WriteJSON(w, http.StatusNotFound, nil)
	}
}

func (u *UserHandlerStruct) UserHandlerById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdFromRequest(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, "path parameter is invalid.")
		return
	}
	token := r.Header.Get("access_token")
	if _, err := auth.VerfiyAccessToken(token); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err)
	}

	if r.Method == http.MethodGet {
		user, err := u.userSvc.GetUserById(r.Context(), id)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, fmt.Errorf("Error while completing request, %v.", err))
			return
		}

		utils.WriteJSON(w, http.StatusOK, user)
	} else if r.Method == http.MethodPut {
		fmt.Println("entered user service")
		var userRequest dtos.UpdatedUserDto
		err = json.NewDecoder(r.Body).Decode(&userRequest)
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, fmt.Errorf("error while parsing body: %v", err))
		}

		updatedUser, err := u.userSvc.UpdateUser(r.Context(), id, userRequest)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, updatedUser)
	} else if r.Method == http.MethodDelete {
		err = u.userSvc.DeleteUser(r.Context(), id)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, fmt.Errorf("Error while completing request, %v.", err))
			return
		}

		utils.WriteJSON(w, http.StatusOK, nil)
	} else {
		utils.WriteJSON(w, http.StatusNotFound, nil)
	}
}

func (u *UserHandlerStruct) IsUserAuthorized(id int, claims jwt.Claims) error {
	mapClaims, _ := claims.(jwt.MapClaims)
	userId := mapClaims["user_id"].(int)
	if id != userId {
		return fmt.Errorf("unauthenticated")
	}

	return nil
}
