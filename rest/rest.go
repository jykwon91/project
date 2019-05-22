package rest

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/jykwon91/project/util/logger"
	"github.com/jykwon91/project/util/constant"
)

type JwtToken struct {
        Token    string `json:"token"`
        UserType string `json:"userType"`
}

func GetTokenPass() ([]byte, error) {
	tokenPass, err := ioutil.ReadFile(constant.TOKENPASSFILE)
	if err != nil {
		logger.AltLogger(err.Error())
		return nil, err
	}

	return tokenPass, nil
}

func ReadReqBody(requestObj *http.Request, resultObj interface{}) error {
        defer requestObj.Body.Close()

        bytes, err := ioutil.ReadAll(requestObj.Body)
        if err != nil {
                return err
        }

        err = json.Unmarshal(bytes, resultObj)
        if err != nil {
                return err
        }

        return nil
}

func ReadRespBody(responseObj *http.Response, resultObj interface{}) error {
        defer responseObj.Body.Close()

        bytes, err := ioutil.ReadAll(responseObj.Body)
        if err != nil {
                return err
        }

        err = json.Unmarshal(bytes, resultObj)
        if err != nil {
                return err
        }

        return nil
}

func AuthenticateTokenAndReturnClaims(tokenString string) (jwt.MapClaims, error) {

	tokenPass, err := GetTokenPass()
	if err != nil {
		logger.AltLogger("Error getting token pass")
		return nil, fmt.Errorf("Failed to authenticate user")
	}

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, fmt.Errorf("Failed to authenticate user")
                }
                return tokenPass, nil
        })
        if err != nil {
                return nil, err
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
                return claims, nil
        }

        return nil, fmt.Errorf("Failed to authenticate user")
}

func GetEmailFromClaims(claims jwt.MapClaims) (string, error) {
	var email string
	for k, v := range claims {
		if strings.EqualFold(k, "email") {
			email = v.(string)
		}
	}

	if len(email) <= 0 {
		return "", fmt.Errorf("Email was not found in claims")
	}
	return email, nil
}
