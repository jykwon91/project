package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"context"

        "github.com/braintree-go/braintree-go"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

const (
	LANDLORD = "landlord"
	TENANT = "tenant"
	RENT = "rent"
	PAID = "paid"
	OPEN = "open"
	PROCESSING = "processing"
	LATE = "late"
	ERROR = "error"
)

type User struct {
	UserID                   string
	UserType                 string
	Password                 []byte
	FirstName                string
	LastName                 string
	RentalAddress            Address
	OwnedPropertyAddressList []Address
	BillingAddress           Address
	PaymentList           []Payment
	ServiceRequestList       []ServiceRequest
	NotificationList         []Notification
	LandLordID                 string
	RentalPaymentAmt            int64
	RentDueDate       string //epoch
	LateFeeRate string //in percentage maybe make this an int
	LegalDocuments           []Document
	Email                    string
	PhoneNumber              string
}

type Address struct {
	AddressID    string
	Street       string
	Zipcode      string
	City         string
	State        string
	PropertyType string
}

type ServiceRequest struct {
	Status          string
	RequestID       string
	RequestTime     string
	StartTime       string
	CompletedTime   string
	Message         string
	RentalAddress Address
	TenantName      string
}

//now := time.Now()
//secs := now.Unix()
//fmt.Println(secs)
//fmt.Println(time.Unix(secs, 0))
type Notification struct {
	NotificationID string
	CreatedOn      string //epoch time
	Message        string
	From           string
}

type Document struct {
	DocumentID    string
	DocumentType  string //receipt, contract, contact, personal
	DocumentBytes []byte
}

// could this be tied to Document
type Payment struct {
	PaymentID string
	LandLordID string
	TenantID string
	BTTransactionID string
	Category   string //rent, service, repairs, utility, etc
	PaymentMethod string
	Amount        int64 //cents
	LateFeeAmount	int64 //cents
	Status string //open, processing, paid, late
	PaidDate string
	DueDate string
	Description string
}

type JwtToken struct {
	Token    string `json:"token"`
	UserType string `json:"userType"`
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

type appError struct {
	Error            error
	ServerLogMessage string
	Message          string
	Code             uint64
}

func (fn appHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if e := fn(resp, req); e != nil {
		logger(e, resp, req)
		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(int(e.Code))
		fmt.Fprintf(resp, `{"type":"error","message":"%s", "code":"%d"}`, e.ServerLogMessage, e.Code)
	}
}

func logger(e *appError, resp http.ResponseWriter, req *http.Request) {
	var logMessage string
	origin := strings.Replace(req.Header.Get("Origin"), "http://", "", -1)
	if e != nil {
		log.Printf("[HTTP %d][%s][%s] - %s - ERROR: %s: %s", e.Code, req.Method, origin, req.RequestURI, e.Message, e.Error)
		logMessage = "[HTTP " + strconv.FormatUint(e.Code, 10) + "][" + req.Method + "][" + origin + "] - " + req.RequestURI + " - ERROR: " + e.Message + ": " + e.Error.Error() + "\n"
	} else {
		log.Printf("[HTTP 200][%s][%s] - %s ", req.Method, origin, req.RequestURI)
		logMessage = "[HTTP 200][" + req.Method + "][" + origin + "] - " + req.RequestURI + "\n"
	}

	f, err := os.OpenFile("/home/jkwon/Git/project/log/server.log", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		altLogger(err.Error())
	}

	_, err = f.WriteString(logMessage)
	if err != nil {
		altLogger(err.Error())
	}

	f.Close()
}

func altLogger(errStr string) {
	logMessage := "ERROR: " + errStr
	f, err := os.OpenFile("/home/jkwon/Git/project/log/server.log", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Can't print to server log")
	}

	_, err = f.WriteString(logMessage)
	if err != nil {
		fmt.Printf("Can't print to server log")
	}

	f.Close()
}

// Middleware function which will be called for each request
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if !strings.EqualFold(req.URL.Path, "/users/authenticate") && !strings.EqualFold(req.URL.Path, "/users/register") && !strings.EqualFold(req.URL.Path, "/stateList") && strings.EqualFold(req.URL.Path, "/tenant/pay") {
			var tokenList []string
			token := req.Header.Get("Authorization")

			bytes, _ := ioutil.ReadFile("/home/jkwon/Git/project/database/Tokens")
			json.Unmarshal(bytes, &tokenList)

			var authenticated bool
			for _, tmpToken := range tokenList {
				if strings.EqualFold(tmpToken, token) {
					authenticated = true
				}
			}
			if authenticated {
				next.ServeHTTP(resp, req)
			} else {
				http.Error(resp, "\"Forbidden. You do not have permission to view this content.\"", http.StatusForbidden)
			}
		} else {
			next.ServeHTTP(resp, req)
		}
	})
}

var tokenPass []byte

func DailyCheckRentDue() {
	for t := range time.NewTicker(86400 * time.Second).C {
	//TEST: create test payment obj
	//for t := range time.NewTicker(8 * time.Second).C {
		fmt.Println(t)

		userList, err := getUserListFromDatabase()
		if err != nil {
			altLogger(err.Error())
		}

		for _, user := range userList {
			//if strings.EqualFold(user.UserType, TENANT) && (t.Day() == user.RentDueDate) {
			if strings.EqualFold(user.UserType, TENANT) {
				now := time.Now()
				secs := now.Unix()
				dueDate := strconv.FormatInt(secs, 10)
				description := "Pay period " + now.Month().String() + " " + strconv.Itoa(now.Year())

				//TODO: Need to redo how I add payments. currently adding to only landlord payment list
				// too confusing to keep track of
				//payment := Payment{ PaymentID: uuid.New().String(), LandLordID: user.LandLordID, TenantID: user.UserID, BTTransactionID: "", Category: RENT, PaymentMethod: "", Status: OPEN, Amount: user.RentalPaymentAmt, PaidDate: "", DueDate: dueDate, Description: description}
				//TEST: production api test, $1
				payment := Payment{ PaymentID: uuid.New().String(), LandLordID: user.LandLordID, TenantID: user.UserID, BTTransactionID: "", Category: RENT, PaymentMethod: "", Status: OPEN, Amount: 100, PaidDate: "", DueDate: dueDate, Description: description}

				for i, tmpUser := range userList {
					if strings.EqualFold(tmpUser.UserID, user.LandLordID) {
						userList[i].PaymentList = append([]Payment{payment}, userList[i].PaymentList...)
					}
				}

				//emailRentDueNotification(user.email or user)
			}
		}

		err = updateUserDatabase(userList)
		if err != nil {
			altLogger(err.Error())
		}
	}
}

func DailyCheckPendingPayments() {
	for t := range time.NewTicker(86400 * time.Second).C {

		var completedPaymentList []Payment

		//TODO: refactor altLogger() to include INFO log messages instead of only ERROR messages.
		fmt.Printf("%v\n", t)
		pendingPaymentList, err := getPendingPayments()
		if err != nil {
			altLogger(err.Error())
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
				altLogger(err.Error())
				pendingPaymentList[i].Status = ERROR
			} else {
				pendingPaymentList[i].Status = PAID
			}
		}

		userList, err := getUserListFromDatabase()
		if err != nil {
			altLogger(err.Error())
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

		err = updateUserDatabase(userList)
		if err != nil {
			altLogger(err.Error())
		}

		emailCompletedPaymentConfirmation(completedPaymentList)
	}
}

func getPendingPayments() ([]Payment, error) {

	var pendingPaymentList []Payment

	bytes, err := ioutil.ReadFile("/home/jkwon/Git/project/database/PendingPayments")
	if err != nil {
		return nil, err
	}
	json.Unmarshal(bytes, &pendingPaymentList)

	return pendingPaymentList, nil
}

func main() {

	go DailyCheckRentDue()
	go DailyCheckPendingPayments()

	var err error
	tokenPass, err = ioutil.ReadFile("/home/jkwon/Git/project/etc/tokenpass")
	if err != nil {
		altLogger(err.Error())
	}

	router := mux.NewRouter()

	router.Handle("/users/authenticate", appHandler(authenticateUser))
	router.Handle("/users/landlord/all", appHandler(getLandLordList))
	router.Handle("/users/register", appHandler(registerUser))
	router.Handle("/users/all", appHandler(getAllUsers))
	router.Handle("/users/notification/all", appHandler(getAllNotifications))
	router.Handle("/users/service/all", appHandler(getServiceRequestList))
	router.Handle("/users/currentUser", appHandler(getCurrentUser))
	router.Handle("/users/update", appHandler(updateUser))
	router.Handle("/users/payment/all", appHandler(getPaymentList))

	router.Handle("/stateList", appHandler(getStateList))
	router.Handle("/users/landlord/property/register", appHandler(registerLandlordProperty))
	router.Handle("/landlord/property/all", appHandler(getAllLandLordProperties))
	router.Handle("/landlord/notification", appHandler(sendNotification))
	router.Handle("/landlord/service/request/update", appHandler(updateServiceRequest))
	router.Handle("/landlord/tenant/all", appHandler(getTenantList))
	router.Handle("/tenant/service/request", appHandler(sendServiceRequest))
	router.Handle("/tenant/pay/{tokenKey}", appHandler(tenantPayment))
	router.Handle("/tenant/payment/overview/{landLordID}", appHandler(getPaymentOverview))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "http://localhost:8080", "http://192.168.1.125", "http://192.168.1.125:8080", "http://rentalmgmt.co:8080", "http://rentalmgmt.co", "http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Accept-Language", "Content-Type", "Authorization", "Bearer"},
		AllowCredentials: true,
		Debug:            false,
	})

	router.Use(AuthMiddleware)

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", handler))
}

func authenticateTokenAndReturnClaims(tokenString string) (jwt.MapClaims, error) {

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

func getPaymentOverview(resp http.ResponseWriter, req *http.Request) *appError {

	vars := mux.Vars(req)
	landLordID := vars["landLordID"]
	type PaymentOverview struct {
		CurrentPayPeriod string
		CurrentAmountDue int64
		TotalLateAmount int64
		TotalLateFees int64
		TotalDue int64
	}

	claims, err := authenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError{err, "Getting payment overview failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	userList, email, err := getUserListFromDatabaseAndUserEmail(claims)
	if err != nil {
		return &appError{err, "Getting payment overview failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var tenantID string
	for _, user := range userList {
		if strings.EqualFold(user.Email, email) {
			tenantID = user.UserID
			break
		}
	}

	var paymentOverview PaymentOverview
	now := time.Now()
	for _, user := range userList {
		if strings.EqualFold(user.UserID, landLordID) {
			for _, payment := range user.PaymentList {
				if !strings.EqualFold(payment.Status, PAID) && strings.EqualFold(payment.TenantID, tenantID) {
					epoch, err := strconv.ParseInt(payment.DueDate, 10, 64)
					if err != nil {
						return &appError{err, "Getting payment overview failed. Server down. Please contact customer support or try again later.", "Failed to convert epoch", 500}
					}
					date := time.Unix(epoch, 0)
					if date.Month() == now.Month() {
						paymentOverview.CurrentPayPeriod = now.Month().String() + " " + strconv.Itoa(now.Year())
						paymentOverview.CurrentAmountDue = paymentOverview.CurrentAmountDue + payment.Amount
						paymentOverview.TotalDue = paymentOverview.TotalDue + payment.Amount
					} else if strings.EqualFold(payment.Status, LATE) {
						paymentOverview.TotalLateAmount = paymentOverview.TotalLateAmount + payment.Amount
						paymentOverview.TotalLateFees = paymentOverview.TotalLateFees + payment.LateFeeAmount
						paymentOverview.TotalDue = paymentOverview.TotalDue + payment.Amount + payment.LateFeeAmount
					}
				}
			}
		}
	}

	bytes, err := json.Marshal(paymentOverview)
	if err != nil {
		return &appError{err, "Getting payment overview failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)
	resp.Write(bytes)

	logger(nil, resp, req)
	return nil
}
func getPaymentList(resp http.ResponseWriter, req *http.Request) *appError {

	claims, err := authenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError{err, "Getting payment list failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	var userType string
	for k, v := range claims {
		if strings.EqualFold(k, "userType") {
			userType = v.(string)
		}
	}

	userList, email, err := getUserListFromDatabaseAndUserEmail(claims)
	if err != nil {
		return &appError{err, "Getting payment list failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var paymentList []Payment
	var tenantID string
	var landLordID string

	for _, user := range userList {
		if strings.EqualFold(userType, TENANT) && strings.EqualFold(user.Email, email) {
			tenantID = user.UserID
			landLordID = user.LandLordID
			break
		} else if strings.EqualFold(userType, LANDLORD) && strings.EqualFold(user.Email, email) {
			paymentList = user.PaymentList
			break
		}
	}

	if strings.EqualFold(userType, TENANT) {
		for _, user := range userList {
			if strings.EqualFold(user.UserID, landLordID) {
				for _, payment := range user.PaymentList {
					if strings.EqualFold(payment.TenantID, tenantID) {
						paymentList = append(paymentList, payment)
					}
				}
				break
			}
		}
	}

	bytes, err := json.Marshal(paymentList)
	if err != nil {
		return &appError{err, "Getting service request list failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)
	resp.Write(bytes)

	logger(nil, resp, req)
	return nil
}

func getLandLordList(resp http.ResponseWriter, req *http.Request) *appError {

	type LandLord struct {
		LandLordID string
		Name string
	}

	userList, err := getUserListFromDatabase()
	if err != nil {
		return &appError{err, "Getting land lord list failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var landLordList []LandLord

	for _, user := range userList {
		if strings.EqualFold(user.UserType, LANDLORD) {
			var landLord LandLord
			name := user.FirstName + " " + user.LastName
			landLord.Name = name
			landLord.LandLordID = user.UserID
			landLordList = append(landLordList, landLord)
		}
	}

	bytes, err := json.Marshal(landLordList)
	if err != nil {
		return &appError{err, "Getting service request list failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)
	resp.Write(bytes)
	logger(nil, resp, req)
	return nil
}

func tenantPayment(resp http.ResponseWriter, req *http.Request) *appError {
	vars := mux.Vars(req)
	tokenKey := vars["tokenKey"]
	type ReqBody struct {
		PaymentID string
		LandLordID string
		TenantID string
		Amount int64
	}

	var reqBody ReqBody
	err := readReqBody(req, &reqBody)
	if err != nil {
		return &appError{err, "Payment failed. Please contact customer support or try again later.", "Failed to read request", 500}
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
	t, err := bt.Transaction().Create(ctx, &braintree.TransactionRequest{
		Type: "sale",
		Amount: braintree.NewDecimal(reqBody.Amount, 2),
		PaymentMethodNonce: tokenKey,
	})
	if err != nil {
		return &appError{err, "Payment failed. Server down. Please contact customer support.", "Payment failed", 403}
	}


	pendingPayment, err := updatePayment(reqBody.PaymentID, reqBody.LandLordID, reqBody.TenantID, reqBody.Amount, t.Id, string(t.PaymentInstrumentType))
	if err != nil {
		return &appError{err, "Creating payment history failed. Server down. Please contact customer support.", "Creating Payment failed", 500}
	}

	err = addToPendingPaymentList(pendingPayment)
	if err != nil {
		return &appError{err, "Updating payment failed. Please contact customer support.", "Adding payment to pending payment list failed", 500}
	}

	err = emailPaymentConfirmation(pendingPayment)
	if err != nil {
		return &appError{err, "Updating payment failed. Please contact customer support.", "Adding payment to pending payment list failed", 500}
	}

	logger(nil, resp, req)
	return nil
}

func updatePayment(paymentID string, landLordID string, tenantID string, amount int64, btTransactionID string, paymentMethod string) (Payment, error) {
	userList, err := getUserListFromDatabase()
	if err != nil {
		return Payment{}, err
	}

	var pendingPayment Payment
	for i, user := range userList {
		if strings.EqualFold(user.UserID, landLordID) {
			for j, payment := range user.PaymentList {
				if strings.EqualFold(paymentID, payment.PaymentID) {
					if amount == payment.Amount {
						now := time.Now()
						secs := now.Unix()
						date := strconv.FormatInt(secs, 10)
						userList[i].PaymentList[j].Status = PROCESSING
						userList[i].PaymentList[j].PaidDate = date
					} else {
						userList[i].PaymentList[j].Amount = userList[i].PaymentList[j].Amount - amount
					}
					userList[i].PaymentList[j].BTTransactionID = btTransactionID
					userList[i].PaymentList[j].PaymentMethod = paymentMethod
					pendingPayment = userList[i].PaymentList[j]
				}
			}
		}
	}

	err = updateUserDatabase(userList)
	if err != nil {
		return Payment{}, err
	}

	return pendingPayment, nil
}

func addToPendingPaymentList(pendingPayment Payment) error {

	var pendingPaymentList []Payment

	bytes, err := ioutil.ReadFile("/home/jkwon/Git/project/database/PendingPayments")
	if err != nil {
		return err
	}
	json.Unmarshal(bytes, &pendingPaymentList)

	pendingPaymentList = append(pendingPaymentList, pendingPayment)

	bytes, err = json.Marshal(pendingPaymentList)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("/home/jkwon/Git/project/database/PendingPayments", bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// change this to update payment
func createPayment( landLordID string, tenantID string, amount int64, btTransactionID string, paymentMethod string, category string, status string) error {

	date := ""
	if strings.EqualFold(PAID, status) {
		now := time.Now()
		secs := now.Unix()
		date = strconv.FormatInt(secs, 10)
	}

	payment := Payment{ LandLordID: landLordID, TenantID: tenantID, BTTransactionID: btTransactionID, Category: category, PaymentMethod: paymentMethod, Status: status, Amount: amount, PaidDate: date}

	userList, err := getUserListFromDatabase()
	if err != nil {
		return err
	}

	for i, user := range userList {
		if strings.EqualFold(user.UserID, landLordID) {
			userList[i].PaymentList = append(userList[i].PaymentList, payment)
		}
	}

	err = updateUserDatabase(userList)
	if err != nil {
		return err
	}

	return nil
}

func getServiceRequestList(resp http.ResponseWriter, req *http.Request) *appError {
	claims, err := authenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError{err, "Getting service request list failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	var userType string
	for k, v := range claims {
		if strings.EqualFold(k, "userType") {
			userType = v.(string)
		}
	}

	userList, email, err := getUserListFromDatabaseAndUserEmail(claims)
	if err != nil {
		return &appError{err, "Getting service request list failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var serviceRequestList []ServiceRequest

	if strings.EqualFold(userType, "tenant") {
		var tenantRentalAddress Address
		for _, user := range userList {
			if strings.EqualFold(user.Email, email) {
				tenantRentalAddress = user.RentalAddress
				break
			}
		}

		var found bool
		for _, user := range userList {
			found = false
			if len(user.OwnedPropertyAddressList) > 0 {
				for _, address := range user.OwnedPropertyAddressList {
					if strings.EqualFold(address.AddressID, tenantRentalAddress.AddressID) {
						found = true
						break
					}
				}
			}
			if found {
				serviceRequestList = user.ServiceRequestList
			}
		}
	} else if strings.EqualFold(userType, LANDLORD) {
		for i, user := range userList {
			if strings.EqualFold(user.Email, email) {
				serviceRequestList = userList[i].ServiceRequestList
			}
		}
	}

	bytes, err := json.Marshal(serviceRequestList)
	if err != nil {
		return &appError{err, "Getting service request list failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)
	resp.Write(bytes)

	logger(nil, resp, req)
	return nil
}

func getTenantList(resp http.ResponseWriter, req *http.Request) *appError {

	claims, err := authenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError{err, "Getting tenant list failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	userList, email, err := getUserListFromDatabaseAndUserEmail(claims)
	if err != nil {
		return &appError{err, "Getting tenant list failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	landLordID := ""
	for _, user := range userList {
		if strings.EqualFold(user.Email, email) {
			landLordID = user.UserID
			break
		}
	}

	if strings.EqualFold(landLordID, "") {
		return &appError{err, "Getting tenant list failed. Server down. Please contact customer support or try again later.", "Could not find land lord", 500}
	}

	var tenantList []User
	for _, user := range userList {
		if strings.EqualFold(user.LandLordID, landLordID) {
			tenantList = append(tenantList, user)
		}
	}

	bytes, err := json.Marshal(tenantList)
	if err != nil {
		return &appError{err, "Getting all landlord properties failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)
	resp.Write(bytes)

	logger(nil, resp, req)
	return nil
}

func updateUser(resp http.ResponseWriter, req *http.Request) *appError {

	_, err := authenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError{err, "Updating user failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	userList, err := getUserListFromDatabase()
	if err != nil {
		return &appError{err, "Updating user failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var updateUserObj User
	err = readReqBody(req, &updateUserObj)
	if err != nil {
		return &appError{err, "Updating user failed. Please contact customer support or try again later.", "Failed to read request", 500}
	}

	for i, user := range userList {
		if strings.EqualFold(user.Email, updateUserObj.Email) {
			userList[i] = updateUserObj
			break
		}
	}

	err = updateUserDatabase(userList)
	if err != nil {
		return &appError{err, "Sending service request failed. Server down. Please contact customer support or try again later.", "Failed to update database", 500}
	}

	logger(nil, resp, req)
	return nil
}

func updateServiceRequest(resp http.ResponseWriter, req *http.Request) *appError {

	claims, err := authenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError{err, "Sending service request failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	userList, email, err := getUserListFromDatabaseAndUserEmail(claims)
	if err != nil {
		return &appError{err, "Sending service request failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var updateServiceReqObj ServiceRequest
	err = readReqBody(req, &updateServiceReqObj)
	if err != nil {
		return &appError{err, "Sending service request failed. Please contact customer support or try again later.", "Failed to read request", 500}
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

	err = updateUserDatabase(userList)
	if err != nil {
		return &appError{err, "Sending service request failed. Server down. Please contact customer support or try again later.", "Failed to update database", 500}
	}

	logger(nil, resp, req)
	return nil
}

func createServiceRequest(message string, address Address, tenantName string) ServiceRequest {
	//need to add function to create Service request
	var serviceReq ServiceRequest
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

func sendServiceRequest(resp http.ResponseWriter, req *http.Request) *appError {
	type ReqService struct {
		Message       string
		TenantName    string
		RentalAddress Address
	}

	claims, err := authenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError{err, "Sending service request failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	userList, _, err := getUserListFromDatabaseAndUserEmail(claims)
	if err != nil {
		return &appError{err, "Sending service request failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var reqService ReqService
	err = readReqBody(req, &reqService)
	if err != nil {
		return &appError{err, "Sending service request failed. Please contact customer support or try again later.", "Failed to read request", 500}
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
			userList[i].ServiceRequestList = append([]ServiceRequest{serviceReq}, userList[i].ServiceRequestList...)
			break
		}
	}

	err = updateUserDatabase(userList)
	if err != nil {
		return &appError{err, "Sending service request failed. Server down. Please contact customer support or try again later.", "Failed to update database", 500}
	}

	logger(nil, resp, req)
	return nil
}

func sendNotification(resp http.ResponseWriter, req *http.Request) *appError {

	type ReqNotification struct {
		Message     string
		AddressList []string
	}

	claims, err := authenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError{err, "Sending notification failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	userList, _, err := getUserListFromDatabaseAndUserEmail(claims)
	if err != nil {
		return &appError{err, "Sending notification failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var reqNotification ReqNotification
	err = readReqBody(req, &reqNotification)
	if err != nil {
		return &appError{err, "Sending notification property failed. Please contact customer support or try again later.", "Failed to read request", 500}
	}

	//need to add function to create notification
	var notification Notification
	notification.Message = reqNotification.Message
	notification.From = "Jason"
	now := time.Now()
	secs := now.Unix()
	notification.CreatedOn = strconv.FormatInt(secs, 10)
	notification.NotificationID = uuid.New().String()

	for _, reqAddress := range reqNotification.AddressList {
		for i, user := range userList {
			if strings.EqualFold(reqAddress, user.RentalAddress.AddressID) {
				userList[i].NotificationList = append([]Notification{notification}, userList[i].NotificationList...)
				break
			}
		}
	}

	err = updateUserDatabase(userList)
	if err != nil {
		return &appError{err, "Sending notification failed. Server down. Please contact customer support or try again later.", "Failed to update database", 500}
	}

	logger(nil, resp, req)
	return nil
}

func getUserFromDataBase(claims jwt.MapClaims) (User, error) {

	var email string
	for k, v := range claims {
		if strings.EqualFold(k, "email") {
			email = v.(string)
		}
	}

	var userList []User
	var theUser User

	bytes, err := ioutil.ReadFile("/home/jkwon/Git/project/database/Users")
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

func getCurrentUser(resp http.ResponseWriter, req *http.Request) *appError {

	claims, err := authenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError{err, "Getting current user failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	currentUser, err := getUserFromDataBase(claims)
	currentUser.Password = nil
	if err != nil {
		return &appError{err, "Getting current user failed. Server down. Please contact customer support or try again later.", "Failed to get current user", 500}
	}

	bytes, err := json.Marshal(currentUser)
	if err != nil {
		return &appError{err, "Getting current user failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)
	resp.Write(bytes)

	logger(nil, resp, req)
	return nil
}

func getUserListFromDatabase() ([]User, error) {

	var userList []User

	bytes, err := ioutil.ReadFile("/home/jkwon/Git/project/database/Users")
	if err != nil {
		return nil, err
	}
	json.Unmarshal(bytes, &userList)

	return userList, nil
}

func getUserListFromDatabaseAndUserEmail(claims jwt.MapClaims) ([]User, string, error) {

	var email string
	for k, v := range claims {
		if strings.EqualFold(k, "email") {
			email = v.(string)
		}
	}

	var userList []User

	bytes, err := ioutil.ReadFile("/home/jkwon/Git/project/database/Users")
	if err != nil {
		return nil, "", err
	}
	json.Unmarshal(bytes, &userList)

	return userList, email, nil
}

func updateUserDatabase(data []User) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("/home/jkwon/Git/project/database/Users", bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func getAllLandLordProperties(resp http.ResponseWriter, req *http.Request) *appError {

	claims, err := authenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError{err, "Getting all landlord properties failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	userList, email, err := getUserListFromDatabaseAndUserEmail(claims)
	if err != nil {
		return &appError{err, "Getting all landlord properties failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var addressList []Address
	for _, user := range userList {
		if strings.EqualFold(user.Email, email) {
			addressList = user.OwnedPropertyAddressList
			break
		}
	}

	bytes, err := json.Marshal(addressList)
	if err != nil {
		return &appError{err, "Getting all landlord properties failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)
	resp.Write(bytes)

	logger(nil, resp, req)
	return nil
}

func registerLandlordProperty(resp http.ResponseWriter, req *http.Request) *appError {

	claims, err := authenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError{err, "Registering landlord property failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	var address Address
	err = readReqBody(req, &address)
	if err != nil {
		return &appError{err, "Registering landlord property failed. Please contact customer support or try again later.", "Failed to read request", 500}
	}

	address.AddressID = uuid.New().String()
	address.PropertyType = "rental"

	userList, email, err := getUserListFromDatabaseAndUserEmail(claims)
	if err != nil {
		return &appError{err, "Registering landlord property failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	for i, user := range userList {
		if strings.EqualFold(user.Email, email) {
			userList[i].OwnedPropertyAddressList = append(userList[i].OwnedPropertyAddressList, address)
			break
		}
	}

	err = updateUserDatabase(userList)
	if err != nil {
		return &appError{err, "Registering landlord property failed. Server down. Please contact customer support or try again later.", "Failed to update database", 500}
	}

	logger(nil, resp, req)
	return nil
}

func getStateList(resp http.ResponseWriter, req *http.Request) *appError {
	type State struct {
		id    int
		value string
		name  string
	}
	var stateList []State

	bytes, err := ioutil.ReadFile("/home/jkwon/Git/project/database/StateListWithName")
	if err != nil {
		return &appError{err, "Getting stateList failed. Server down. Please contact customer support or try again later.", "Failed to read database file", 500}
	}
	json.Unmarshal(bytes, &stateList)

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)
	resp.Write(bytes)

	logger(nil, resp, req)
	return nil
}

func getAllNotifications(resp http.ResponseWriter, req *http.Request) *appError {
	claims, err := authenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError{err, "Getting notifications failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	userList, email, err := getUserListFromDatabaseAndUserEmail(claims)
	if err != nil {
		return &appError{err, "Getting notifications failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var notificationList []Notification
	for _, user := range userList {
		if strings.EqualFold(user.Email, email) {
			notificationList = user.NotificationList
			break
		}
	}

	bytes, err := json.Marshal(notificationList)
	if err != nil {
		return &appError{err, "Getting notifications failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)
	resp.Write(bytes)

	logger(nil, resp, req)
	return nil
}

func registerUser(resp http.ResponseWriter, req *http.Request) *appError {

	type SelectObj struct {
		Value string
		Id    int64
	}
	type LandLordObj struct {
		LandLordID string
		Name string
	}
	type ReqBody struct {
		FirstName      string
		LastName       string
		Password       string
		Email          string
		LandLord  LandLordObj
		BillingStreet  string
		BillingCity    string
		BillingZipcode string
		BillingState   SelectObj
		PhoneNumber    string
		RentalPaymentAmt int64
	}

	var reqBody ReqBody
	err := readReqBody(req, &reqBody)
	if err != nil {
		return &appError{err, "Registration failed. Please contact customer support or try again later.", "Failed to read request", 500}
	}

	userList, err := getUserListFromDatabase()
	if err != nil {
		return &appError{err, "Getting land lord list failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
	}

	var user User
	user.UserID = uuid.New().String()
	user.FirstName = reqBody.FirstName
	user.LastName = reqBody.LastName
	user.UserType = TENANT
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
		return &appError{err, "Registration failed. Please contact customer support or try again later.", "Failed to hash password", 500}
	}
	user.Password = hash

	bytes, err := json.Marshal(user)
	if err != nil {
		return &appError{err, "Registration failed. Please contact customer support or try again later.", "Failed to marshal json", 400}
	}

	for _, tmpUser := range userList {
		if compareUsers(tmpUser, user) {
			return &appError{errors.New("Email already registered"), "Email already registered. Please contact customer support for log in info.", "Duplicate user, did not register", 500}
		}
	}

	userList = append(userList, user)

	bytes, err = json.Marshal(userList)
	if err != nil {
		return &appError{err, "Registration failed. Please contact customer support or try again later.", "Failed to create json file", 500}
	}

	err = ioutil.WriteFile("/home/jkwon/Git/project/database/Users", bytes, 0644)
	if err != nil {
		return &appError{err, "Registration failed. Please contact customer support or try again later.", "Failed to write new user to file", 500}
	}

	logger(nil, resp, req)
	return nil
}

func compareUsers(a User, b User) bool {
	if a.Email == b.Email {
		return true
	} else {
		return false
	}
}

func authenticateUser(resp http.ResponseWriter, req *http.Request) *appError {

	var userObj User
	var userList []User
	type ReqBody struct {
		Password string
		Email    string
	}

	var reqBody ReqBody
	err := readReqBody(req, &reqBody)
	if err != nil {
		return &appError{err, "Login failed. Server down. Please contact customer support or try again later.", "Failed to read request body", 500}
	}

	bytes, err := ioutil.ReadFile("/home/jkwon/Git/project/database/Users")
	if err != nil {
		return &appError{err, "Login failed. Server down. Please contact customer support or try again later.", "Failed to read database file", 500}
	}
	json.Unmarshal(bytes, &userList)

	found := false
	for _, user := range userList {
		if strings.EqualFold(user.Email, reqBody.Email) {
			if err := bcrypt.CompareHashAndPassword(user.Password, []byte(reqBody.Password)); err != nil {
				return &appError{err, "Login failed. Wrong password was entered. Contact customer support for forgotten password", "Wrong password", 403}
			}
			found = true
			userObj = user
			break
		}
	}

	if !found {
		return &appError{errors.New("Wrong email was entered"), "Login failed. Wrong email was entered. Contact customer support for forgotten email", "Wrong email", 403}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    userObj.Email,
		"userType": userObj.UserType,
	})

	tokenString, err := token.SignedString(tokenPass)
	if err != nil {
		return &appError{err, "Login failed. Server down. Please contact customer support or try again later", "Failed to create token", 500}
	}

	var tokenList []string
	bytes, err = ioutil.ReadFile("/home/jkwon/Git/project/database/Tokens")
	if err != nil {
		return &appError{err, "Login failed. Server down. Please contact customer support or try again later.", "Failed to get tokens", 500}
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
			return &appError{err, "Login failed. Please contact customer support or try again later.", "Failed to marshal json", 500}
		}

		err = ioutil.WriteFile("/home/jkwon/Git/project/database/Tokens", bytes, 0644)
		if err != nil {
			return &appError{err, "Login failed. Please contact customer support or try again later.", "Failed to write token to file", 500}
		}
	}

	json.NewEncoder(resp).Encode(JwtToken{Token: tokenString, UserType: userObj.UserType})

	logger(nil, resp, req)
	return nil
}

func getAllUsers(resp http.ResponseWriter, req *http.Request) *appError {

	claims, err := authenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError{err, "Getting all users failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	var userType string
	for k, v := range claims {
		if strings.EqualFold(k, "userType") {
			userType = v.(string)
		}
	}

	if !strings.EqualFold(userType, LANDLORD) {
		return &appError{errors.New("Forbidden"), "Getting all users failed. You do not have permission to view this content.", "Forbidden", 403}
	}
	var userList []User

	type RespBody struct {
		FirstName string
		LastName  string
		Email     string
	}
	var respBody []RespBody

	bytes, err := ioutil.ReadFile("/home/jkwon/Git/project/database/Users")
	if err != nil {
		return &appError{err, "Getting all users failed. Server down. Please contact customer support or try again later.", "Failed to read database file", 500}
	}
	json.Unmarshal(bytes, &userList)

	for _, user := range userList {
		tmpUser := RespBody{FirstName: user.FirstName, LastName: user.LastName, Email: user.Email}
		respBody = append(respBody, tmpUser)
	}

	bytes, err = json.Marshal(respBody)
	if err != nil {
		return &appError{err, "Getting all users failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)
	resp.Write(bytes)

	logger(nil, resp, req)
	return nil
}

func readReqBody(requestObj *http.Request, resultObj interface{}) error {
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

func readRespBody(responseObj *http.Response, resultObj interface{}) error {
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

func emailRentDueNotification() {
	email, err := ioutil.ReadFile("/home/jkwon/Git/project/businessEmail")
	if err != nil {
		altLogger(err.Error())
	}

	emailPass, err := ioutil.ReadFile("/home/jkwon/Git/project/businessEmailPass")
	if err != nil {
		altLogger(err.Error())
	}
	var userList []User

	bytes, err := ioutil.ReadFile("/home/jkwon/Git/project/database/Users")
	if err != nil {
		altLogger(err.Error())
	}
	json.Unmarshal(bytes, &userList)

	t := time.Now()
	year := strconv.Itoa(t.Year())
	month := t.Month().String()
	day := strconv.Itoa(t.Day())
	theDate := month + " " + day + " " + year

	for _, user := range userList {
		if strings.EqualFold(user.UserType, TENANT) {
			m := gomail.NewMessage()
			m.SetHeader("From", string(email))
			m.SetHeader("To", user.Email)
			m.SetHeader("Subject", "Rent due for "+month+" "+year)
			m.SetBody("text/html", `
				<p><b>Hi `+user.FirstName+` `+user.LastName+`</b></p><br><p>This is a reminder that rent is due today(`+theDate+`)</p><br>
				<p>Log into www.rentalmgmt.co to pay.</p>
			`)

			d := gomail.NewDialer("smtp.gmail.com", 587, string(email), string(emailPass))

			if err := d.DialAndSend(m); err != nil {
				altLogger(err.Error())
			}
		}
	}

	//TEST
	//kvs := map[string]string{"Alan Ayala": "jasonykwon91@gmail.com", "Laura Smith": "jasonykwon91@gmail.com"}

	/* TEST
	for name, tenantEmail := range kvs {
		m := gomail.NewMessage()
		m.SetHeader("From", string(email))
		m.SetHeader("To", tenantEmail)
		m.SetHeader("Subject", "Rent due for "+month+" "+year)
		m.SetBody("text/html", `<p><b>Hi `+name+`</b></p><br><p>This is a reminder that rent is due today(`+theDate+`)</p>`)

		d := gomail.NewDialer("smtp.gmail.com", 587, string(email), string(emailPass))

		// Send the email to Bob, Cora and Dan.
		if err := d.DialAndSend(m); err != nil {
			altLogger(err.Error())
		}
	}*/
}

func emailPaymentConfirmation(paymentInfo Payment) error {

	email, err := ioutil.ReadFile("/home/jkwon/Git/project/businessEmail")
	if err != nil {
		return err
	}

	emailPass, err := ioutil.ReadFile("/home/jkwon/Git/project/businessEmailPass")
	if err != nil {
		return err
	}
	var userList []User

	bytes, err := ioutil.ReadFile("/home/jkwon/Git/project/database/Users")
	if err != nil {
		return err
	}
	json.Unmarshal(bytes, &userList)

	t := time.Now()
	year := strconv.Itoa(t.Year())
	month := t.Month().String()
	day := strconv.Itoa(t.Day())
	theDate := month + " " + day + " " + year

	for _, user := range userList {
		if strings.EqualFold(user.UserID, paymentInfo.TenantID) {
			name := user.FirstName + " "+ user.LastName
			m := gomail.NewMessage()
			m.Embed("/home/jkwon/Git/project/signature.jpg")
			m.SetHeader("From", string(email))
			m.SetHeader("To", user.Email)
			m.SetHeader("Subject", "Rent Receipt for "+month+" "+year)
			m.SetBody("text/html", `
				<p>Thank you for your payment! This is a confirmation email. We will process your payment within 24 hours.</p><br>
				<table width='600' style='border:1px solid #333'>
				<tbody>
					<tr><td align='left'><b>Transaction number:</b>`+paymentInfo.PaymentID+`</td></tr>
					<tr><td align='left'><b>Brain Tree Transaction number:</b>`+paymentInfo.BTTransactionID+`</td></tr>
					<tr><td align='left'><b>Name:</b> `+name+`</td></tr>
					<tr><td align='left'><b>Date:</b> `+theDate+`</td></tr>
					<tr><td align='left'><b>Transaction Type:</b>Card</td></tr>
					<tr>
						<td align='center'>
							<table align='center' width='300' border='0' cellspacing='0' cellpadding='0' style='border:1px solid #ccc; padding:10px 0px 10px 10px'>
								<tr>
									<td><b>Category:</b></td>
									<td>$`+paymentInfo.Category+`</td>
								</tr>
								<tr>
									<td><b>Amount Due:</b></td>
									<td>$`+strconv.FormatInt(paymentInfo.Amount,10)+`</td>
								</tr>
								<tr>
									<td><b>Amount Paid:</b></td>
									<td>$`+strconv.FormatInt(paymentInfo.Amount,10)+`</td>
								</tr>
								<tr>
									<td><b>Received by:</b></td>
									<td>Jason Kwon</td>
								</tr>
								<tr>
									<td><b>Signature:</b></td>
									<td><img src='cid:signature.jpg' alt='My Image' style='max-width: 100px; max-height: 100px' /></td>
								</tr>
							</table>
					</tr>
					<br>
				</tbody>`)

			d := gomail.NewDialer("smtp.gmail.com", 587, string(email), string(emailPass))

			if err := d.DialAndSend(m); err != nil {
				return err
			}
		}
	}

	return nil
}

func emailCompletedPaymentConfirmation(completedPaymentList []Payment) {

	email, err := ioutil.ReadFile("/home/jkwon/Git/project/businessEmail")
	if err != nil {
		altLogger(err.Error())
	}

	emailPass, err := ioutil.ReadFile("/home/jkwon/Git/project/businessEmailPass")
	if err != nil {
		altLogger(err.Error())
	}
	var userList []User

	bytes, err := ioutil.ReadFile("/home/jkwon/Git/project/database/Users")
	if err != nil {
		altLogger(err.Error())
	}
	json.Unmarshal(bytes, &userList)

	t := time.Now()
	year := strconv.Itoa(t.Year())
	month := t.Month().String()
	//day := strconv.Itoa(t.Day())
	//theDate := month + " " + day + " " + year

	for _, user := range userList {
		for _, payment := range completedPaymentList {
			if strings.EqualFold(payment.TenantID, user.UserID) {
				m := gomail.NewMessage()
				m.Embed("/home/jkwon/Git/project/signature.jpg")
				m.SetHeader("From", string(email))
				m.SetHeader("To", user.Email)
				m.SetHeader("Subject", "Payment confirmation for "+month+" "+year)
				m.SetBody("text/html", `
						<p> This is to confirm your payment has finished processing.</p>		
				`)

				d := gomail.NewDialer("smtp.gmail.com", 587, string(email), string(emailPass))

				// Send the email to Bob, Cora and Dan.
				if err := d.DialAndSend(m); err != nil {
					altLogger(err.Error())
				}
			}
		}
	}
}
