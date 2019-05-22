package rest

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strings"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"

	"github.com/jykwon91/project/db"
	"github.com/jykwon91/project/util/appError"
	"github.com/jykwon91/project/util/logger"
	"github.com/jykwon91/project/util/constant"
)

func AuthenticateUser(resp http.ResponseWriter, req *http.Request) *appError.AppError {

        var userObj db.UserData
        var userList []db.UserData
        type ReqBody struct {
                Password string
                Email    string
        }

        var reqBody ReqBody
        err := ReadReqBody(req, &reqBody)
        if err != nil {
                return &appError.AppError{err, "Login failed. Server down. Please contact customer support or try again later.", "Failed to read request body", 500}
        }

        bytes, err := ioutil.ReadFile(constant.USERFILE)
        if err != nil {
                return &appError.AppError{err, "Login failed. Server down. Please contact customer support or try again later.", "Failed to read database file", 500}
        }
        json.Unmarshal(bytes, &userList)

        found := false
        for _, user := range userList {
                if strings.EqualFold(user.Email, reqBody.Email) {
                        if err := bcrypt.CompareHashAndPassword(user.Password, []byte(reqBody.Password)); err != nil {
                                return &appError.AppError{err, "Login failed. Wrong password was entered. Contact customer support for forgotten password", "Wrong password", 403}
                        }
                        found = true
                        userObj = user
                        break
                }
        }

        if !found {
                return &appError.AppError{errors.New("Wrong email was entered"), "Login failed. Wrong email was entered. Contact customer support for forgotten email", "Wrong email", 403}
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
                "email":    userObj.Email,
                "userType": userObj.UserType,
        })

	tokenPass, err := GetTokenPass()
        if err != nil {
                return &appError.AppError{err, "Login failed. Server down. Please contact customer support or try again later", "Failed to create token", 500}
        }

        tokenString, err := token.SignedString(tokenPass)
        if err != nil {
                return &appError.AppError{err, "Login failed. Server down. Please contact customer support or try again later", "Failed to create token", 500}
        }

        var tokenList []string
        bytes, err = ioutil.ReadFile(constant.TOKENFILE)
        if err != nil {
                return &appError.AppError{err, "Login failed. Server down. Please contact customer support or try again later.", "Failed to get tokens", 500}
        }
        json.Unmarshal(bytes, &tokenList)

        found = false
        for _, tokenObj := range tokenList {
                if strings.EqualFold(tokenObj, tokenString) {
                        found = true
                        break
                }
        }

        if !found {
                tokenList = append(tokenList, tokenString)
                bytes, err = json.Marshal(tokenList)
                if err != nil {
                        return &appError.AppError{err, "Login failed. Please contact customer support or try again later.", "Failed to marshal json", 500}
                }

                err = ioutil.WriteFile(constant.TOKENFILE, bytes, 0644)
                if err != nil {
                        return &appError.AppError{err, "Login failed. Please contact customer support or try again later.", "Failed to write token to file", 500}
                }
        }

        json.NewEncoder(resp).Encode(JwtToken{Token: tokenString, UserType: userObj.UserType})

        logger.Logger(nil, resp, req)
        return nil
}

func GetLandLordList(resp http.ResponseWriter, req *http.Request) *appError.AppError {

        type LandLord struct {
                LandLordID string
                Name       string
        }

        userList, err := db.User.GetUserList()
        if err != nil {
                return &appError.AppError{err, "Getting land lord list failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
        }

        var landLordList []LandLord

        for _, user := range userList {
                if strings.EqualFold(user.UserType, constant.LANDLORD) {
                        var landLord LandLord
                        name := user.FirstName + " " + user.LastName
                        landLord.Name = name
                        landLord.LandLordID = user.UserID
                        landLordList = append(landLordList, landLord)
                }
        }

        bytes, err := json.Marshal(landLordList)
        if err != nil {
                return &appError.AppError{err, "Getting service request list failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
        }

        resp.Header().Set("Content-Type", "application/json")
        resp.WriteHeader(200)
        resp.Write(bytes)
        logger.Logger(nil, resp, req)
        return nil
}

