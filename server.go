package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
	"github.com/dgrijalva/jwt-go"
//	"github.com/logpacker/PayPal-Go-SDK"
	"gopkg.in/gomail.v2"
)

type User struct {
	UserID                string
	UserType              string
	Password              []byte
	FirstName             string
	LastName              string
	RentalHomeAddress     Address //Maybe make a struct for address
	OwnedPropertyAddressList     []Address //Maybe make a struct for address
	BillingAddress        Address
	PaymentHistory        []Payment
	ServiceRequestHistory []ServiceRequest
	NotificationList	[]Notification
	Landlord              string
	RentalPayment         string
	LegalDocuments        []Document //need to figure this out
	Email                 string
	PhoneNumber           string
}

type Address struct {
	Street  string
	Zipcode string
	City    string
	State   string
	PropertyType string
}

type ServiceRequest struct {
	Processing bool
	Completed bool
	SrvRequestID string
	SrvRequestTime string //epoch time - time.Unix(secs, 0) to print date
	SrvRequestBody string
	RentalAddress Address
	Tenant      User
}

//now := time.Now()
//secs := now.Unix()
//fmt.Println(secs)
//fmt.Println(time.Unix(secs, 0))
type Notification struct {
	NotificationID string
	CreatedOn string //epoch time
	Message string
	From string
}

type Document struct {
	DocumentID string
	DocumentType  string //receipt, contract, contact, personal
	DocumentBytes []byte
}

// could this be tied to Document
type Payment struct {
	PaymentID string
	PaymentType string
	PaymentMethod string
	Amount      string
	Paid bool
	DueDate string //epoch
}

type JwtToken struct {
	Token string `json:"token"`
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
		if !strings.EqualFold(req.URL.Path, "/users/authenticate") && !strings.EqualFold(req.URL.Path, "/users/register") && !strings.EqualFold(req.URL.Path, "/stateList") {
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

func DailyTaskExec() {
	for t := range time.NewTicker(86400 * time.Second).C {
		if t.Day() == 1 {
			emailRentDueNotification()
		}
	}
}

func main() {

	go DailyTaskExec()

	var err error
	tokenPass, err = ioutil.ReadFile("/home/jkwon/Git/project/etc/tokenpass")
	if err != nil {
		altLogger(err.Error())
	}

	router := mux.NewRouter()

	router.Handle("/users/authenticate", appHandler(authenticateUser))
	router.Handle("/users/register", appHandler(registerUser))
	router.Handle("/users/all", appHandler(getAllUsers))
	router.Handle("/users/notification/all", appHandler(getAllNotifications))
	router.Handle("/stateList", appHandler(getStateList))
	router.Handle("/users/landlord/property/register", appHandler(registerLandlordProperty))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "http://localhost:8080", "http://192.168.1.125", "http://192.168.1.125:8080","http://rentalmgmt.co:8080","http://rentalmgmt.co"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Accept-Language", "Content-Type", "Authorization", "Bearer"},
		AllowCredentials: true,
		Debug:            false,
	})

	router.Use(AuthMiddleware)

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", handler))
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

	//kvs := map[string]string{"Alan Ayala": "kanonbolt128@gmail.com", "Laura Smith": "smithlaura9295@gmail.com"}
	kvs := map[string]string{"Alan Ayala": "jasonykwon91@gmail.com", "Laura Smith": "jasonykwon91@gmail.com"}

	t := time.Now()
	year := strconv.Itoa(t.Year())
	month := t.Month().String()
	day := strconv.Itoa(t.Day())
	theDate := month + " " + day + " " + year

	for name, tenantEmail := range kvs {
		m := gomail.NewMessage()
		m.SetHeader("From", string(email))
		m.SetHeader("To", tenantEmail)
		m.SetHeader("Subject", "Rent due for " + month + " " + year)
		m.SetBody("text/html",`<p><b>Hi `+name+`</b></p><br><p>This is a reminder that rent is due today(`+ theDate +`)</p>`)


		d := gomail.NewDialer("smtp.gmail.com", 587, string(email), string(emailPass))

		// Send the email to Bob, Cora and Dan.
		if err := d.DialAndSend(m); err != nil {
			altLogger(err.Error())
		}
	}
}

func emailReceipts() {

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

	//kvs := map[string]string{"Alan Ayala": "kanonbolt128@gmail.com", "Laura Smith": "smithlaura9295@gmail.com"}
	kvs := map[string]string{"Alan Ayala": "jasonykwon91@gmail.com", "Laura Smith": "jasonykwon91@gmail.com"}

	t := time.Now()
	year := strconv.Itoa(t.Year())
	month := t.Month().String()
	day := strconv.Itoa(t.Day())
	theDate := month + " " + day + " " + year

	for name, tenantEmail := range kvs {
		m := gomail.NewMessage()
		m.Embed("/home/jkwon/Git/project/signature.jpg")
		m.SetHeader("From", string(email))
		m.SetHeader("To", tenantEmail)
		m.SetHeader("Subject", "Rent Receipt for " + month + " " + year)
		m.SetBody("text/html",`
			<table width='600' style='border:1px solid #333'
			<tbody>
				<tr><td align='left'><b>Transaction number:</b>1234-5678-1234-5678</td></tr>
				<tr><td align='left'><b>Name:</b> ` + name + `</td></tr>
				<tr><td align='left'><b>Date:</b> ` + theDate + `</td></tr>
				<tr><td align='left'><b>Transaction Type:</b> Cash</td></tr>
				<tr>
					<td align='center'>
						<table align='center' width='300' border='0' cellspacing='0' cellpadding='0' style='border:1px solid #ccc; padding:10px 0px 10px 10px'>
							<tr>
								<td><b>Amount Due:</b></td>
								<td>$550</td>
							</tr>
							<tr>
								<td><b>Amount Paid:</b></td>
								<td>$550</td>
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

		// Send the email to Bob, Cora and Dan.
		if err := d.DialAndSend(m); err != nil {
			altLogger(err.Error())
		}
	}
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

func registerLandlordProperty(resp http.ResponseWriter, req *http.Request) *appError {

	token, err := jwt.Parse(req.Header.Get("Authorization"), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Failed to authenticate user")
		}
		return tokenPass, nil
	})
	if err != nil {
		return &appError{err, "Registering landlord property failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		var address Address
		err := readReqBody(req, &address)
		if err != nil {
			return &appError{err, "Registering landlord property failed. Please contact customer support or try again later.", "Failed to read request", 500}
		}

		userList, email, err := getUserListFromDatabaseAndUserEmail(claims)
		if err != nil {
			return &appError{err, "Registering landlord property failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
		}

		for i, user := range userList {
			if strings.EqualFold(user.Email, email) {
				userList[i].OwnedPropertyAddressList = append(userList[i].OwnedPropertyAddressList, address)
			}
		}

		err = updateUserDatabase(userList)
		if err != nil {
			return &appError{err, "Registering landlord property failed. Server down. Please contact customer support or try again later.", "Failed to update database", 500}
		}

	} else {
		return &appError{errors.New("Failed to authenticate user"), "Registering landlord property failed. Server down. Please contact customer support or try again later", "Failed to authenticate user", 403}
	}

	logger(nil, resp, req)
	return nil
}

func getStateList(resp http.ResponseWriter, req *http.Request) *appError {
	type State struct {
		id int
		value string
		name string
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
	token, err := jwt.Parse(req.Header.Get("Authorization"), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Failed to authenticate user")
		}
		return tokenPass, nil
	})
	if err != nil {
		return &appError{err, "Getting notifications failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		var email string
		for k, v := range claims {
			if strings.EqualFold(k, "email") {
				email = v.(string)
			}
		}

		var userList []User

		bytes, err := ioutil.ReadFile("/home/jkwon/Git/project/database/Users")
		if err != nil {
			return &appError{err, "Getting notifications failed. Server down. Please contact customer support or try again later.", "Failed to read database file", 500}
		}
		json.Unmarshal(bytes, &userList)

		var notificationList []Notification
		for _, user := range userList {
			if strings.EqualFold(user.Email, email) {
				notificationList = user.NotificationList
			}
		}

		bytes, err = json.Marshal(notificationList)
		if err != nil {
			return &appError{err, "Getting notifications failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
		}

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(200)
		resp.Write(bytes)
	} else {
		return &appError{errors.New("Failed to authenticate user"), "Get notifications failed. Server down. Please contact customer support or try again later", "Failed to authenticate user", 403}
	}

	logger(nil, resp, req)
	return nil
}

func registerUser(resp http.ResponseWriter, req *http.Request) *appError {

	var user User
	var userList []User
	var rentalAdd Address
	type SelectObj struct {
		Value string
		Id int64
	}

	type ReqBody struct {
		FirstName string
		LastName  string
		Password  string
		Email string
		RentalAddress SelectObj
		BillingStreet string
		BillingCity string
		BillingZipcode string
		BillingState SelectObj
		PhoneNumber string
	}
	var reqBody ReqBody

	err := readReqBody(req, &reqBody)
	if err != nil {
		return &appError{err, "Registration failed. Please contact customer support or try again later.", "Failed to read request", 500}
	}

	//need to use rental address by id, not string, to relate to landlord
	tmpString := strings.Split(reqBody.RentalAddress.Value, ",")
	rentalAdd.Street = tmpString[0]
	rentalAdd.City = tmpString[2]
	rentalAdd.State = tmpString[1]
	rentalAdd.Zipcode = tmpString[3]

	user.UserID = uuid.New().String()
	user.FirstName = reqBody.FirstName
	user.LastName = reqBody.LastName
	user.UserType = "tenant"
	user.RentalHomeAddress = rentalAdd
	user.BillingAddress.Street = reqBody.BillingStreet
	user.BillingAddress.Zipcode = reqBody.BillingZipcode
	user.BillingAddress.City = reqBody.BillingCity
	user.BillingAddress.State = reqBody.BillingState.Value

	//search landlord by rental address with API
	user.Landlord = "Jason"
	user.Email = reqBody.Email
	user.PhoneNumber = reqBody.PhoneNumber

	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return &appError{err, "Registration failed. Please contact customer support or try again later.", "Failed to hash password", 500}
	}
	user.Password = hash

	bytes, err := json.Marshal(user)
	//bytes, err := json.Marshal(make(chan int))
	if err != nil {
		return &appError{err, "Registration failed. Please contact customer support or try again later.", "Failed to marshal json", 400}
	}

	bytes, err = ioutil.ReadFile("/home/jkwon/Git/project/database/Users")
	if err != nil {
		return &appError{err, "Registration failed. Please contact customer support or try again later.", "Failed to read database file", 500}
	}
	json.Unmarshal(bytes, &userList)

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

func compareUsers (a User, b User) bool {
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
		Email string
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
			if err := bcrypt.CompareHashAndPassword(user.Password, []byte(reqBody.Password)); err!= nil {
				return &appError{err, "Login failed. Wrong password was entered. Contact customer support for forgotten password", "Wrong password", 403}
			}
			found = true
			userObj = user
		}
	}

	if !found {
		return &appError{errors.New("Wrong email was entered"), "Login failed. Wrong email was entered. Contact customer support for forgotten email", "Wrong email", 403}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": userObj.Email,
		"password": userObj.Password,
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

	tokenList = append(tokenList, tokenString)
	bytes, err = json.Marshal(tokenList)
	if err != nil {
		return &appError{err, "Login failed. Please contact customer support or try again later.", "Failed to marshal json", 500}
	}

	err = ioutil.WriteFile("/home/jkwon/Git/project/database/Tokens", bytes, 0644)
	if err != nil {
		return &appError{err, "Registration failed. Please contact customer support or try again later.", "Failed to write new user to file", 500}
	}

	json.NewEncoder(resp).Encode(JwtToken{Token: tokenString, UserType: userObj.UserType})

	logger(nil, resp, req)
	return nil
}

func getAllUsers(resp http.ResponseWriter, req *http.Request) *appError {

	token, err := jwt.Parse(req.Header.Get("Authorization"), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Failed to authenticate user")
		}
		return tokenPass, nil
	})
	if err != nil {
		return &appError{err, "Getting all users failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		var userType string
		for k, v := range claims {
			if strings.EqualFold(k, "userType") {
				userType = v.(string)
			}
		}

		if !strings.EqualFold(userType, "landlord") {
			return &appError{errors.New("Forbidden"), "Getting all users failed. You do not have permission to view this content.", "Forbidden", 403}
		}
		var userList []User

		type RespBody struct {
			FirstName string
			LastName string
			Email string
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
	} else {
		return &appError{errors.New("Failed to authenticate user"), "Get all users failed. Server down. Please contact customer support or try again later", "Failed to authenticate user", 403}
	}

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
