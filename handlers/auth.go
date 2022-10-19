package handlers

import (
	"AlexSarva/media/admin"
	"AlexSarva/media/models"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const timeLayout = "2006-01-02 15:04:05"

// UserRegistration - user registration method
//
// Handler POST /api/user/register
//
// Registration is performed by a pair of login/password.
// Each login must be set.
// After successful registration, automatic user authentication is required.
// post message should contain such body:
//
//	"login": "<login>",
//	"password": "<password>"
//
// Possible response codes:
// 200 - user successfully registered and authenticated;
// 400 - invalid request format;
// 409 - login is already taken;
// 500 - an internal server error.
func UserRegistration(database *admin.PostgresDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Type")
		if !strings.Contains("application/json, application/x-gzip", headerContentType) {
			messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
			return
		}

		var user models.User
		var unmarshalErr *json.UnmarshalTypeError

		b, err := readBodyBytes(r)
		if err != nil {
			messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(b)
		decoder.DisallowUnknownFields()
		errDecode := decoder.Decode(&user)

		if errDecode != nil {
			if errors.As(errDecode, &unmarshalErr) {
				messageResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
			} else {
				messageResponse(w, "Bad Request. "+errDecode.Error(), "application/json", http.StatusBadRequest)
			}
			return
		}

		userID := uuid.New()
		userToken, userTokenExp := GenerateToken(userID)
		hashedPassword, bcrypteErr := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
		if bcrypteErr != nil {
			log.Println(bcrypteErr)
		}

		user.ID, user.Password, user.Token, user.TokenExp = userID, string(hashedPassword), userToken, userTokenExp

		newUserErr := database.RegisterUser(&user)
		if newUserErr != nil {
			if newUserErr == admin.ErrDuplicatePK {
				messageResponse(w, "login is already busy", "application/json", http.StatusConflict)
				return
			}
			messageResponse(w, "Internal Server Error "+newUserErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		tokenDetails := models.Token{
			Username: user.Username,
			Email:    user.Email,
			Type:     "Bearer",
			Token:    userToken,
			TokenExp: userTokenExp,
		}
		jsonResp, _ := json.Marshal(tokenDetails)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Authorization", tokenDetails.Type+" "+userToken)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}
}

// UserAuthentication - user authentication method
//
// Handler POST /api/user/login
//
// Authentication is performed by a login/password pair.
// Request format:
//
//	{"login": "<login>",
//	"password": "<password>"}
//
// Possible response codes:
// 200 - user successfully authenticated;
// 400 - invalid request format;
// 401 - invalid login/password pair;
// 500 - an internal server error.
func UserAuthentication(database *admin.PostgresDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Type")
		if !strings.Contains("application/json, application/x-gzip", headerContentType) {
			messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
			return
		}

		var user models.UserLogin
		var unmarshalErr *json.UnmarshalTypeError

		b, err := readBodyBytes(r)
		if err != nil {
			messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(b)
		decoder.DisallowUnknownFields()
		errDecode := decoder.Decode(&user)
		log.Printf("%+v\n", user)

		if errDecode != nil {
			if errors.As(errDecode, &unmarshalErr) {
				messageResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
			} else {
				messageResponse(w, "Bad Request. "+errDecode.Error(), "application/json", http.StatusBadRequest)
			}
			return
		}

		userDB, userDBErr := database.LoginUser(user.Email)
		if userDBErr != nil {
			if errors.Is(userDBErr, sql.ErrNoRows) {
				messageResponse(w, "email doesnt exist", "application/json", http.StatusUnauthorized)
				return
			}
			messageResponse(w, "Internal Server Error: "+userDBErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		cryptErr := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password))
		if cryptErr != nil {
			messageResponse(w, "password doesnt match", "application/json", http.StatusUnauthorized)
			return
		}
		// TODO Предусмотреть обновление куки
		if userDB.TokenExp.Before(time.Now()) {
			log.Println("cookie expired")
		}

		// Авторизация по токену
		token, tokenExp := userDB.Token, userDB.TokenExp

		tokenDetails := models.Token{
			Username: userDB.Username,
			Email:    userDB.Email,
			Type:     "Bearer",
			Token:    token,
			TokenExp: tokenExp,
		}
		jsonResp, _ := json.Marshal(tokenDetails)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Authorization", tokenDetails.Type+" "+token)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}
}

func GetUserInfo(database *admin.PostgresDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Length")
		if len(headerContentType) != 0 {
			messageResponse(w, "Content-Length is not equal 0", "application/json", http.StatusBadRequest)
			return
		}

		// Проверка авторизации по токену
		userID, tokenErr := GetToken(r)
		if tokenErr != nil {
			messageResponse(w, "User unauthorized: "+tokenErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		userInfo, userInfoErr := database.GetUserInfo(userID)
		if userInfoErr != nil {
			if errors.Is(userInfoErr, sql.ErrNoRows) {
				messageResponse(w, "user doesnt exist", "application/json", http.StatusUnauthorized)
				return
			}
			messageResponse(w, "Internal Server Error: "+userInfoErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		userInfo.Type = "Bearer"

		jsonResp, _ := json.Marshal(userInfo)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}
}
