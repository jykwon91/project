package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/jykwon91/project/db"
	"github.com/jykwon91/project/util/appError"
	"github.com/jykwon91/project/util/constant"
	"github.com/jykwon91/project/util/logger"
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

func RegisterUser(resp http.ResponseWriter, req *http.Request) *appError.AppError {

	type SelectObj struct {
		Value string
		Id    int64
	}
	type LandLordObj struct {
		LandLordID string
		Name       string
	}
	type ReqBody struct {
		FirstName        string
		LastName         string
		Password         string
		Email            string
		LandLord         LandLordObj
		BillingStreet    string
		BillingCity      string
		BillingZipcode   string
		BillingState     SelectObj
		PhoneNumber      string
		RentalPaymentAmt int64
	}

	var reqBody ReqBody
	err := ReadReqBody(req, &reqBody)
	if err != nil {
		return &appError.AppError{err, "Registration failed. Please contact customer support or try again later.", "Failed to read request", 500}
	}

	userList, err := db.User.GetUserList()
	if err != nil {
		return &appError.AppError{err, "Getting land lord list failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var user db.UserData
	user.UserID = uuid.New().String()
	user.FirstName = reqBody.FirstName
	user.LastName = reqBody.LastName
	user.UserType = constant.TENANT
	user.LandLordID = reqBody.LandLord.LandLordID
	user.BillingAddress.Street = reqBody.BillingStreet
	user.BillingAddress.Zipcode = reqBody.BillingZipcode
	user.BillingAddress.City = reqBody.BillingCity
	user.BillingAddress.State = reqBody.BillingState.Value
	user.Email = reqBody.Email
	user.PhoneNumber = reqBody.PhoneNumber
	user.RentalPaymentAmt = reqBody.RentalPaymentAmt

	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return &appError.AppError{err, "Registration failed. Please contact customer support or try again later.", "Failed to hash password", 500}
	}
	user.Password = hash

	bytes, err := json.Marshal(user)
	if err != nil {
		return &appError.AppError{err, "Registration failed. Please contact customer support or try again later.", "Failed to marshal json", 400}
	}

	for _, tmpUser := range userList {
		if compareUsers(tmpUser, user) {
			return &appError.AppError{errors.New("Email already registered"), "Email already registered. Please contact customer support for log in info.", "Duplicate user, did not register", 500}
		}
	}

	userList = append(userList, user)

	bytes, err = json.Marshal(userList)
	if err != nil {
		return &appError.AppError{err, "Registration failed. Please contact customer support or try again later.", "Failed to create json file", 500}
	}

	err = ioutil.WriteFile(constant.USERFILE, bytes, 0644)
	if err != nil {
		return &appError.AppError{err, "Registration failed. Please contact customer support or try again later.", "Failed to write new user to file", 500}
	}

	logger.Logger(nil, resp, req)
	return nil
}

func compareUsers(a db.UserData, b db.UserData) bool {
	if a.Email == b.Email {
		return true
	} else {
		return false
	}
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

func UpdateUser(resp http.ResponseWriter, req *http.Request) *appError.AppError {

	_, err := AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError.AppError{err, "Updating user failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	userList, err := db.User.GetUserList()
	if err != nil {
		return &appError.AppError{err, "Updating user failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var updateUserObj db.UserData
	err = ReadReqBody(req, &updateUserObj)
	if err != nil {
		return &appError.AppError{err, "Updating user failed. Please contact customer support or try again later.", "Failed to read request", 500}
	}

	for i, user := range userList {
		if strings.EqualFold(user.Email, updateUserObj.Email) {
			userList[i] = updateUserObj
			break
		}
	}

	err = db.User.UpdateUserList(userList)
	if err != nil {
		return &appError.AppError{err, "Sending service request failed. Server down. Please contact customer support or try again later.", "Failed to update database", 500}
	}

	logger.Logger(nil, resp, req)
	return nil
}

func GetCurrentUser(resp http.ResponseWriter, req *http.Request) *appError.AppError {

	claims, err := AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError.AppError{err, "Getting current user failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	email, err := GetEmailFromClaims(claims)
	if err != nil {
		return &appError.AppError{err, "Getting current user failed. Server down. Please contact customer support or try again later.", "Failed to get email from claims", 403}
	}

	currentUser, err := db.User.GetUser(email)
	currentUser.Password = nil
	if err != nil {
		return &appError.AppError{err, "Getting current user failed. Server down. Please contact customer support or try again later.", "Failed to get current user", 500}
	}

	bytes, err := json.Marshal(currentUser)
	if err != nil {
		return &appError.AppError{err, "Getting current user failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)
	resp.Write(bytes)

	logger.Logger(nil, resp, req)
	return nil
}

func GetAllUsers(resp http.ResponseWriter, req *http.Request) *appError.AppError {

	claims, err := AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError.AppError{err, "Getting all users failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	var userType string
	for k, v := range claims {
		if strings.EqualFold(k, "userType") {
			userType = v.(string)
		}
	}

	if !strings.EqualFold(userType, constant.LANDLORD) {
		return &appError.AppError{errors.New("Forbidden"), "Getting all users failed. You do not have permission to view this content.", "Forbidden", 403}
	}
	var userList []db.UserData

	type RespBody struct {
		FirstName string
		LastName  string
		Email     string
	}
	var respBody []RespBody

	bytes, err := ioutil.ReadFile(constant.USERFILE)
	if err != nil {
		return &appError.AppError{err, "Getting all users failed. Server down. Please contact customer support or try again later.", "Failed to read database file", 500}
	}
	json.Unmarshal(bytes, &userList)

	for _, user := range userList {
		tmpUser := RespBody{FirstName: user.FirstName, LastName: user.LastName, Email: user.Email}
		respBody = append(respBody, tmpUser)
	}

	bytes, err = json.Marshal(respBody)
	if err != nil {
		return &appError.AppError{err, "Getting all users failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)
	resp.Write(bytes)

	logger.Logger(nil, resp, req)
	return nil
}

func UploadDocument(resp http.ResponseWriter, req *http.Request) *appError.AppError {

	claims, err := AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError.AppError{err, "Uploading document failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	email, err := GetEmailFromClaims(claims)
	if err != nil {
		return &appError.AppError{err, "Uploading document failed. Server down. Please contact customer support or try again later.", "Failed to get email from claims", 403}
	}

	logger.Logger(nil, resp, req)
	return nil
}
