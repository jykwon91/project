package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/braintree-go/braintree-go"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	//"github.com/jykwon91/project/db/address"
	//"github.com/jykwon91/project/db/payment"
	"github.com/jykwon91/project/db"
	"github.com/jykwon91/project/util/email"
	"github.com/jykwon91/project/util/constant"
	"github.com/jykwon91/project/util/appError"
	"github.com/jykwon91/project/util/logger"
	"github.com/jykwon91/project/rest"
)

var tokenPass []byte

func DailyCheckRentDue() {
	for t := range time.NewTicker(constant.TWELVE_HOURS).C {
		//TEST: create test payment obj
		//for t := range time.NewTicker(8 * time.Second).C {
		fmt.Println(t)

		userList, err := db.User.GetUserList()
		if err != nil {
			logger.AltLogger(err.Error())
		}

		for _, user := range userList {
			if strings.EqualFold(user.UserType, constant.TENANT) && (t.Day() == 1) {
				//TEST:
				//if strings.EqualFold(user.UserType, TENANT) {
				tenantName := user.FirstName + " " + user.LastName
				rentalAddress := user.RentalAddress.Street + "," + user.RentalAddress.Zipcode + "," + user.RentalAddress.City + "," + user.RentalAddress.State
				now := time.Now()
				secs := now.Unix()
				dueDate := strconv.FormatInt(secs, 10)
				description := "Pay period " + now.Month().String() + " " + strconv.Itoa(now.Year())

				//TODO: Need to redo how I add payments. currently adding to only landlord payment list
				// too confusing to keep track of
				//payment := Payment{ TenantName: tenantName, RentalAddress: rentalAddress, PaymentID: uuid.New().String(), LandLordID: user.LandLordID, TenantID: user.UserID, BTTransactionID: "", Category: RENT, PaymentMethod: "", Status: OPEN, Amount: user.RentalPaymentAmt, PaidDate: "", DueDate: dueDate, Description: description}
				//TEST: production api test, $1
				payment := db.PaymentData{TenantName: tenantName, RentalAddress: rentalAddress, PaymentID: uuid.New().String(), LandLordID: user.LandLordID, TenantID: user.UserID, BTTransactionID: "", Category: constant.RENT, PaymentMethod: "", Status: constant.OPEN, Amount: 100, PaidDate: "", DueDate: dueDate, Description: description}

				for i, tmpUser := range userList {
					if strings.EqualFold(tmpUser.UserID, user.LandLordID) {
						userList[i].PaymentList = append([]db.PaymentData{payment}, userList[i].PaymentList...)
					}
				}

				//emailRentDueNotification(user.email or user)
			}
		}

		err = db.User.UpdateUserList(userList)
		if err != nil {
			logger.AltLogger(err.Error())
		}
	}
}

func DailyCheckPendingPayments() {
	for t := range time.NewTicker(constant.TWELVE_HOURS).C {

		var completedPaymentList []db.PaymentData

		//TODO: refactor altLogger() to include INFO log messages instead of only ERROR messages.
		fmt.Printf("%v\n", t)
		pendingPaymentList, err := getPendingPayments()
		if err != nil {
			logger.AltLogger(err.Error())
		}

		//Production: production api
		bt := braintree.New(
			braintree.Production,
			"kc2j6g7k7gnvz8nj",
			"fx65ws9kvkkqtp68",
			"73e89ee295205330104dca83df884b7a",
		)

		/* TEST: Sandbox API
		bt := braintree.New(
			braintree.Sandbox,
				"k5yn2w9sq696n7br",
				"x88xbrkyzq49h47b",
				"261c7177b5cb9228f1cf4e4a0ac13c91",
			)
		*/

		ctx := context.Background()

		for i, payment := range pendingPaymentList {
			t, err := bt.Transaction().SubmitForSettlement(ctx, payment.BTTransactionID)
			//TODO: refactor altLogger() to include INFO log messages instead of only ERROR messages.
			fmt.Printf("%v\n", t)
			if err != nil {
				logger.AltLogger(err.Error())
				pendingPaymentList[i].Status = constant.ERROR
			} else {
				pendingPaymentList[i].Status = constant.PAID
			}
		}

		userList, err := db.User.GetUserList()
		if err != nil {
			logger.AltLogger(err.Error())
		}

		for _, pendingPayment := range pendingPaymentList {
			for i, user := range userList {
				if strings.EqualFold(user.UserID, pendingPayment.LandLordID) {
					for j, payment := range user.PaymentList {
						if strings.EqualFold(pendingPayment.PaymentID, payment.PaymentID) {
							userList[i].PaymentList[j] = pendingPayment
							completedPaymentList = append(completedPaymentList, pendingPayment)
						}
					}
				}
			}
		}

		err = db.User.UpdateUserList(userList)
		if err != nil {
			logger.AltLogger(err.Error())
		}

		email.EmailCompletedPaymentConfirmation(completedPaymentList)
	}
}

func getPendingPayments() ([]db.PaymentData, error) {

	var pendingPaymentList []db.PaymentData

	bytes, err := ioutil.ReadFile(constant.PENDINGPAYMENTSFILE)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(bytes, &pendingPaymentList)

	return pendingPaymentList, nil
}

func main() {

	go DailyCheckRentDue()
	go DailyCheckPendingPayments()

	rest.InitRestClient()
}

/*
// change this to update payment
func createPayment(tenantName string, rentalAddress string, landLordID string, tenantID string, amount int64, btTransactionID string, paymentMethod string, category string, status string) error {

	date := ""
	if strings.EqualFold(constant.PAID, status) {
		now := time.Now()
		secs := now.Unix()
		date = strconv.FormatInt(secs, 10)
	}

	thePayment := payment.Payment{TenantName: tenantName, RentalAddress: rentalAddress, LandLordID: landLordID, TenantID: tenantID, BTTransactionID: btTransactionID, Category: category, PaymentMethod: paymentMethod, Status: status, Amount: amount, PaidDate: date}

	userList, err := db.User.GetUserList()
	if err != nil {
		return err
	}

	for i, user := range userList {
		if strings.EqualFold(user.UserID, landLordID) {
			userList[i].PaymentList = append(userList[i].PaymentList, thePayment)
		}
	}

	err = db.User.UpdateUserList(userList)
	if err != nil {
		return err
	}

	return nil
}
*/


func getTenantList(resp http.ResponseWriter, req *http.Request) *appError.AppError {

	claims, err := rest.AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError.AppError{err, "Getting tenant list failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	email, err := rest.GetEmailFromClaims(claims)
	if err != nil {
		return &appError.AppError{err, "Getting tenant list failed. Server down. Please contact customer support or try again later.", "Failed to get email", 500}
	}

	userList, err := db.User.GetUserList()
	if err != nil {
		return &appError.AppError{err, "Getting tenant list failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	landLordID := ""
	for _, user := range userList {
		if strings.EqualFold(user.Email, email) {
			landLordID = user.UserID
			break
		}
	}

	if strings.EqualFold(landLordID, "") {
		return &appError.AppError{err, "Getting tenant list failed. Server down. Please contact customer support or try again later.", "Could not find land lord", 500}
	}

	var tenantList []db.UserData
	for _, user := range userList {
		if strings.EqualFold(user.LandLordID, landLordID) {
			tenantList = append(tenantList, user)
		}
	}

	bytes, err := json.Marshal(tenantList)
	if err != nil {
		return &appError.AppError{err, "Getting all landlord properties failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)
	resp.Write(bytes)

	logger.Logger(nil, resp, req)
	return nil
}

func updateUser(resp http.ResponseWriter, req *http.Request) *appError.AppError {

	_, err := rest.AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError.AppError{err, "Updating user failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	userList, err := db.User.GetUserList()
	if err != nil {
		return &appError.AppError{err, "Updating user failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var updateUserObj db.UserData
	err = rest.ReadReqBody(req, &updateUserObj)
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

func updateServiceRequest(resp http.ResponseWriter, req *http.Request) *appError.AppError {

	claims, err := rest.AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError.AppError{err, "Sending service request failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	email, err := rest.GetEmailFromClaims(claims)
	if err != nil {
		return &appError.AppError{err, "Sending service request failed. Server down. Please contact customer support or try again later.", "Failed to get email", 403}
	}

	userList, err := db.User.GetUserList()
	if err != nil {
		return &appError.AppError{err, "Sending service request failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var updateServiceReqObj db.ServiceRequestData
	err = rest.ReadReqBody(req, &updateServiceReqObj)
	if err != nil {
		return &appError.AppError{err, "Sending service request failed. Please contact customer support or try again later.", "Failed to read request", 500}
	}

	for i, user := range userList {
		if strings.EqualFold(user.Email, email) {
			for j, req := range user.ServiceRequestList {
				if strings.EqualFold(req.RequestID, updateServiceReqObj.RequestID) {
					userList[i].ServiceRequestList[j] = updateServiceReqObj
					break
				}
			}
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

func createServiceRequest(message string, address db.AddressData, tenantName string) db.ServiceRequestData {
	//need to add function to create Service request
	var serviceReq db.ServiceRequestData
	serviceReq.Status = "open"
	serviceReq.RequestID = uuid.New().String()
	now := time.Now()
	secs := now.Unix()
	serviceReq.RequestTime = strconv.FormatInt(secs, 10)
	serviceReq.StartTime = ""
	serviceReq.CompletedTime = ""
	serviceReq.Message = message
	serviceReq.RentalAddress = address
	serviceReq.TenantName = tenantName

	return serviceReq
}

func sendServiceRequest(resp http.ResponseWriter, req *http.Request) *appError.AppError {
	type ReqService struct {
		Message       string
		TenantName    string
		RentalAddress db.AddressData
	}

	_, err := rest.AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError.AppError{err, "Sending service request failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	userList, err := db.User.GetUserList()
	if err != nil {
		return &appError.AppError{err, "Sending service request failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var reqService ReqService
	err = rest.ReadReqBody(req, &reqService)
	if err != nil {
		return &appError.AppError{err, "Sending service request failed. Please contact customer support or try again later.", "Failed to read request", 500}
	}

	serviceReq := createServiceRequest(reqService.Message, reqService.RentalAddress, reqService.TenantName)

	var foundLandLord bool
	for i, user := range userList {
		foundLandLord = false
		if len(user.OwnedPropertyAddressList) > 0 {
			for _, address := range user.OwnedPropertyAddressList {
				if strings.EqualFold(address.AddressID, serviceReq.RentalAddress.AddressID) {
					foundLandLord = true
					break
				}
			}
		}
		if foundLandLord {
			userList[i].ServiceRequestList = append([]db.ServiceRequestData{serviceReq}, userList[i].ServiceRequestList...)
			email.EmailServiceReq(user, serviceReq)
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

func sendNotification(resp http.ResponseWriter, req *http.Request) *appError.AppError {

	type ReqNotification struct {
		Message     string
		AddressList []string
	}

	_, err := rest.AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError.AppError{err, "Sending notification failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	userList, err := db.User.GetUserList()
	if err != nil {
		return &appError.AppError{err, "Sending notification failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var reqNotification ReqNotification
	err = rest.ReadReqBody(req, &reqNotification)
	if err != nil {
		return &appError.AppError{err, "Sending notification property failed. Please contact customer support or try again later.", "Failed to read request", 500}
	}

	//need to add function to create notification
	var notification db.NotificationData
	notification.Message = reqNotification.Message
	notification.From = "Jason"
	now := time.Now()
	secs := now.Unix()
	notification.CreatedOn = strconv.FormatInt(secs, 10)
	notification.NotificationID = uuid.New().String()

	for _, reqAddress := range reqNotification.AddressList {
		for i, user := range userList {
			if strings.EqualFold(reqAddress, user.RentalAddress.AddressID) {
				userList[i].NotificationList = append([]db.NotificationData{notification}, userList[i].NotificationList...)
				break
			}
		}
	}

	err = db.User.UpdateUserList(userList)
	if err != nil {
		return &appError.AppError{err, "Sending notification failed. Server down. Please contact customer support or try again later.", "Failed to update database", 500}
	}

	logger.Logger(nil, resp, req)
	return nil
}

func getUserFromDataBase(claims jwt.MapClaims) (db.UserData, error) {

	var email string
	for k, v := range claims {
		if strings.EqualFold(k, "email") {
			email = v.(string)
		}
	}

	var userList []db.UserData
	var theUser db.UserData

	bytes, err := ioutil.ReadFile(constant.USERFILE)
	if err != nil {
		return theUser, err
	}

	json.Unmarshal(bytes, &userList)
	for _, user := range userList {
		if strings.EqualFold(email, user.Email) {
			theUser = user
			break
		}
	}

	return theUser, nil
}

func getCurrentUser(resp http.ResponseWriter, req *http.Request) *appError.AppError {

	claims, err := rest.AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError.AppError{err, "Getting current user failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	currentUser, err := getUserFromDataBase(claims)
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

func getAllLandLordProperties(resp http.ResponseWriter, req *http.Request) *appError.AppError {

	claims, err := rest.AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError.AppError{err, "Getting all landlord properties failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	email, err := rest.GetEmailFromClaims(claims)
	if err != nil {
		return &appError.AppError{err, "Getting all landlord properties failed. Server down. Please contact customer support or try again later.", "Failed to get email", 403}
	}

	userList, err := db.User.GetUserList()
	if err != nil {
		return &appError.AppError{err, "Getting all landlord properties failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var addressList []db.AddressData
	for _, user := range userList {
		if strings.EqualFold(user.Email, email) {
			addressList = user.OwnedPropertyAddressList
			break
		}
	}

	bytes, err := json.Marshal(addressList)
	if err != nil {
		return &appError.AppError{err, "Getting all landlord properties failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)
	resp.Write(bytes)

	logger.Logger(nil, resp, req)
	return nil
}

func registerLandlordProperty(resp http.ResponseWriter, req *http.Request) *appError.AppError {

	claims, err := rest.AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError.AppError{err, "Registering landlord property failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	var address db.AddressData
	err = rest.ReadReqBody(req, &address)
	if err != nil {
		return &appError.AppError{err, "Registering landlord property failed. Please contact customer support or try again later.", "Failed to read request", 500}
	}

	address.AddressID = uuid.New().String()
	address.PropertyType = "rental"

	email, err := rest.GetEmailFromClaims(claims)
	if err != nil {
		return &appError.AppError{err, "Registering landlord property failed. Server down. Please contact customer support or try again later.", "Failed to get email", 500}
	}

	userList, err := db.User.GetUserList()
	if err != nil {
		return &appError.AppError{err, "Registering landlord property failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	for i, user := range userList {
		if strings.EqualFold(user.Email, email) {
			userList[i].OwnedPropertyAddressList = append(userList[i].OwnedPropertyAddressList, address)
			break
		}
	}

	err = db.User.UpdateUserList(userList)
	if err != nil {
		return &appError.AppError{err, "Registering landlord property failed. Server down. Please contact customer support or try again later.", "Failed to update database", 500}
	}

	logger.Logger(nil, resp, req)
	return nil
}

func getStateList(resp http.ResponseWriter, req *http.Request) *appError.AppError {
	type State struct {
		id    int
		value string
		name  string
	}
	var stateList []State

	bytes, err := ioutil.ReadFile("/home/jkwon/Git/project/database/StateListWithName")
	if err != nil {
		return &appError.AppError{err, "Getting stateList failed. Server down. Please contact customer support or try again later.", "Failed to read database file", 500}
	}
	json.Unmarshal(bytes, &stateList)

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)
	resp.Write(bytes)

	logger.Logger(nil, resp, req)
	return nil
}

func getAllNotifications(resp http.ResponseWriter, req *http.Request) *appError.AppError {
	claims, err := rest.AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError.AppError{err, "Getting notifications failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	email, err := rest.GetEmailFromClaims(claims)
	if err != nil {
		return &appError.AppError{err, "Getting notifications failed. Server down. Please contact customer support or try again later.", "Failed to get email", 403}
	}

	userList, err := db.User.GetUserList()
	if err != nil {
		return &appError.AppError{err, "Getting notifications failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var notificationList []db.NotificationData
	for _, user := range userList {
		if strings.EqualFold(user.Email, email) {
			notificationList = user.NotificationList
			break
		}
	}

	bytes, err := json.Marshal(notificationList)
	if err != nil {
		return &appError.AppError{err, "Getting notifications failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)
	resp.Write(bytes)

	logger.Logger(nil, resp, req)
	return nil
}

func registerUser(resp http.ResponseWriter, req *http.Request) *appError.AppError {

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
	err := rest.ReadReqBody(req, &reqBody)
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

func getAllUsers(resp http.ResponseWriter, req *http.Request) *appError.AppError {

	claims, err := rest.AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
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
