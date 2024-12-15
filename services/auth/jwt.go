package auth

import (
	"context"
	"strconv"
	"time"
	"net/http"
	"log"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/charucjoshi/go_jwt/configs"
	"github.com/charucjoshi/go_jwt/utils"
	"github.com/charucjoshi/go_jwt/types"
)

type contextKey string
var UserKey contextKey

func WithJWTAuth(handleFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.GetTokenFromRequest(r)

		token, err := validateJWT(tokenString)

		if err != nil {
			log.Println("error in validating the token")
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("token invalid!")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, err := strconv.Atoi(str)
		
		if err != nil {
			log.Println("error in converting userID string to int")
			permissionDenied(w)
			return
		}

		u, err := store.GetUserByID(userID)
		
		if err != nil {
			log.Println("error in getting userStore")
			permissionDenied(w)
			return
		}
		
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		handleFunc(w, r)
	} 
}

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(configs.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(userID)),
		"expiresAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		tokenString = ""
	}
	return tokenString, err
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(configs.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}	

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}
	return userID
}
