package util

import (
	"errors"
	"os"
	"time"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte(os.Getenv("SecretKey"))

func GenerateJwt(userId string, role roles.Role) (string, error) {
	claims:= jwt.MapClaims{}

	claims["authorized"]="true"

	claims["user_id"]=userId
	claims["role"]=role
	claims["exp"]=time.Now().Add(time.Hour*24).Unix()
    
	// generate a token with the claims

	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	return token.SignedString(JwtSecret)
}

func VarifyJwt(tokenString string)(*jwt.Token,error){

	token,err:=jwt.Parse(tokenString,func(token *jwt.Token)(interface{},error){
		if _,ok:=token.Method.(*jwt.SigningMethodHMAC);!ok{
			return nil,errors.New("invalid signing method")
		}
		return JwtSecret,nil
	})
    
	if err!=nil || !token.Valid{
       return nil, errors.New("invalid token")
	}
    
	return token,nil
}