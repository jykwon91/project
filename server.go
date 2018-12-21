package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID                string
	UserName              string
	UserType              string
	Password              []byte
	FirstName             string
	LastName              string
	RentalHomeAddress     Address //Maybe make a struct for address
	BillingAddress        Address
	PaymentHistory        []Payment
	ServiceRequestHistory []ServiceRequest
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
}

type ServiceRequest struct {
	RequestTime int64 //epoch time - time.Unix(secs, 0) to print date
	RequestBody string
	Tenant      User
}

type Document struct {
	DocumentType  string //receipt, contract, contact, personal
	DocumentBytes []byte
}

// could this be tied to Document
type Payment struct {
	PaymentType string
	Amount      string
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

type appError struct {
	Error            error
	ServerLogMessage string
	Message          string
	Code             int
}

func (fn appHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if e := fn(resp, req); e != nil {
		logger(e, resp, req)
		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(resp, `{"type":"error","message":"%s", "code":"%d"}`, e.ServerLogMessage, e.Code)
	}
}

func logger(e *appError, resp http.ResponseWriter, req *http.Request) {
	var logMessage string
	origin := strings.Replace(req.Header.Get("Origin"), "http://", "", -1)
	if e != nil {
		log.Printf("[HTTP %d][%s][%s] - %s - ERROR: %s: %s", e.Code, req.Method, origin, req.RequestURI, e.Message, e.Error)
		logMessage = "[HTTP " + strconv.FormatInt(int64(e.Code), 16) + "][" + req.Method + "][" + origin + "] - " + req.RequestURI + " - ERROR: " + e.Message + ": " + e.Error.Error() + "\n"
	} else {
		log.Printf("[HTTP 200][%s][%s] - %s ", req.Method, origin, req.RequestURI)
		logMessage = "[HTTP 200][" + req.Method + "][" + origin + "] - " + req.RequestURI + "\n"
	}

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

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/users/authenticate", authenticateUser).Methods("POST", "OPTIONS")
	router.Handle("/users/register", appHandler(registerUser))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "http://localhost:8080", "http://192.168.1.125", "http://192.168.1.125:8080"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Accept-Language", "Content-Type"},
		AllowCredentials: true,
		Debug:            false,
	})

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", handler))
}

func registerUser(resp http.ResponseWriter, req *http.Request) *appError {

	var user User
	type ReqBody struct {
		FirstName string
		LastName  string
		Username  string
		Password  string
	}
	var reqBody ReqBody

	err := readReqBody(req, &reqBody)

	user.FirstName = reqBody.FirstName
	user.LastName = reqBody.LastName
	user.UserName = reqBody.Username

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

	err = ioutil.WriteFile("/home/jkwon/Git/project/database/Users", bytes, 0644)
	if err != nil {
		return &appError{err, "Registration failed. Please contact customer support or try again later.", "Failed to write new user to file", 500}
	}

	logger(nil, resp, req)
	return nil
}

func authenticateUser(resp http.ResponseWriter, req *http.Request) {

	type ReqBody struct {
		Username string
		Password string
	}

	var reqBody ReqBody
	err := readReqBody(req, &reqBody)
	fmt.Println(reqBody.Username)
	fmt.Println(reqBody.Password)
	/*
		now := time.Now()
		secs := now.Unix()
		fmt.Println(secs)
		fmt.Println(time.Unix(secs, 0))
	*/
	userPassword1 := "hello my name is jason"
	hash, err := bcrypt.GenerateFromPassword([]byte(userPassword1), bcrypt.DefaultCost)
	if err != nil {
		// TODO: Properly handle error
		log.Fatal(err)
	}
	fmt.Println("Hash to store:", string(hash))
	fmt.Print(hash)
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
