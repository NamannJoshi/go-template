package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"template/data/dtos"
	"template/services"
	"template/utils"
)

type IAuthHandler interface {
	Login(http.ResponseWriter, *http.Request)
	Register(http.ResponseWriter, *http.Request)
	RefreshAccessToken(http.ResponseWriter, *http.Request)
	Logout(http.ResponseWriter, *http.Request)
}

type AuthHandlerStruct struct {
	authSvc services.IAuthService
	credSvc services.ICredentialService
}

func AuthHandlerInit(userSvc services.IAuthService, credSvc services.ICredentialService) *AuthHandlerStruct {
	return &AuthHandlerStruct{authSvc: userSvc, credSvc: credSvc}
}

func (a *AuthHandlerStruct) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("welcome inside login")
	if r.Method == http.MethodPost {
		fmt.Println("welcome inside login method")

		var loginDto dtos.LoginDto
		if err := json.NewDecoder(r.Body).Decode(&loginDto); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, err)
			return
		}

		creds, err := a.authSvc.Login(r.Context(), loginDto)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, err)
			return
		}

		createCookieHandler(w, creds.RefreshToken)

		credResponse := dtos.LoginResponseDto{
			AccessToken:  creds.AccessToken,
			RefreshToken: creds.RefreshToken,
			ExpiresIn:    creds.ExpiresIn,
		}

		fmt.Println(credResponse)

		utils.WriteJSON(w, http.StatusOK, credResponse)
	} else {
		utils.WriteJSON(w, http.StatusNotFound, nil)
	}
}

func (a *AuthHandlerStruct) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var registerDto dtos.RegisterDto
		json.NewDecoder(r.Body).Decode(&registerDto)

		if err := a.authSvc.Register(r.Context(), registerDto); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, err)
			return
		}
	} else {
		utils.WriteJSON(w, http.StatusNotFound, nil)
	}
}

func (a *AuthHandlerStruct) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var userId int
		json.NewDecoder(r.Body).Decode(&userId)

		// if err := a.credSvc.DeleteCredential(r.Context(), userId); err != nil {
		// 	utils.WriteJSON(w, http.StatusInternalServerError, err)
		// 	return
		// }

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/auth/refresh",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})

		utils.WriteJSON(w, http.StatusOK, nil)
	} else {
		utils.WriteJSON(w, http.StatusNotFound, nil)
	}
}

func (a *AuthHandlerStruct) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var refreshToken dtos.RefreshTokenDto
		if err := json.NewDecoder(r.Body).Decode(&refreshToken); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, nil)
			return
		}

		accessToken, err := a.authSvc.RefreshAccessToken(refreshToken)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, accessToken)
	} else {
		utils.WriteJSON(w, http.StatusNotFound, nil)
	}
}

func createCookieHandler(w http.ResponseWriter, refreshToken string) {
	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}
