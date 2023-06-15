package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"gitlab.com/api-gateway/pkg/logger"
)


type JWTHandler struct{
	Sub string 
	Iss string
	Exp string
	Iat string
	Aud	[]string
	Role string
	SigninKey string
	Log	logger.Logger
	Token string
}

// GenerateAuthJWT	
func (jwtHandler *JWTHandler) GenerateAuthJWT() ([]string, error){
	var (
		accessToken *jwt.Token
		refreshToken *jwt.Token
		claims jwt.MapClaims
	)
	accessToken = jwt.New(jwt.SigningMethodHS256)
	refreshToken = jwt.New(jwt.SigningMethodHS256)
	
	claims = accessToken.Claims.(jwt.MapClaims)
	claims["iss"] = jwtHandler.Iss
	claims["sub"] = jwtHandler.Sub
	claims["exp"] = time.Now().Add(time.Hour * 500).Unix()
	claims["iat"] = time.Now().Unix()
	claims["role"] = jwtHandler.Role
	claims["aud"] = jwtHandler.Aud

	access, err := accessToken.SignedString([]byte(jwtHandler.SigninKey))
	if err != nil{
		jwtHandler.Log.Error("Error generating access token", logger.Error(err))
		return []string{access, ""}, err
	}
	refresh, err := refreshToken.SignedString([]byte(jwtHandler.SigninKey))
	if err != nil{
		jwtHandler.Log.Error("Error generating refresh token", logger.Error(err))
		return []string{access, refresh}, err
	}

	return []string{access, refresh}, nil
}

func (jwtHandler *JWTHandler) ExtractClaims() (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err error
	)

	token, err = jwt.Parse(jwtHandler.Token, func(t *jwt.Token) (interface{}, error){
		return []byte(jwtHandler.SigninKey), nil
	})
		if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid){
		jwtHandler.Log.Error("Invalid jwt Token")
		return nil, err
	}

	return claims, nil
}