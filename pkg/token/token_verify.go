package token

import (
	"Excel-Props/pkg/config"
	"Excel-Props/pkg/constant"
	"Excel-Props/pkg/utils"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	UserID string
	jwt.RegisteredClaims
}

func BuildClaims(userID string, ttlDay int64) Claims {
	now := time.Now()
	before := now.Add(-time.Minute * 5)
	return Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(ttlDay*24) * time.Hour)), //Expiration time
			IssuedAt:  jwt.NewNumericDate(now),                                           //Issuing time
			NotBefore: jwt.NewNumericDate(before),                                        //Begin Effective time
		}}
}

func CreateToken(UserID string, ttlDay int64) (string, error) {
	claims := BuildClaims(UserID, ttlDay)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.TokenPolicy.AccessSecret))
	if err != nil {
		return "", utils.Wrap(err, "")
	}
	return tokenString, utils.Wrap(err, "")
}
func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.TokenPolicy.AccessSecret), nil
	}
}

func GetUserIDFromToken(tokensString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokensString, &Claims{}, secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return "", utils.Wrap(constant.ErrTokenMalformed, "")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return "", utils.Wrap(constant.ErrTokenExpired, "")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return "", utils.Wrap(constant.ErrTokenNotValidYet, "")
			} else {
				return "", utils.Wrap(constant.ErrTokenUnknown, "")
			}
		} else {
			return "", utils.Wrap(constant.ErrTokenNotValidYet, "")
		}
	} else {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims.UserID, nil
		}
		return "", utils.Wrap(constant.ErrTokenNotValidYet, "")
	}
}
