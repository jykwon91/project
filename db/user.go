package db

import (
	"encoding/json"
	"github.com/jykwon91/project/util/constant"
	"io/ioutil"
	"strings"
)

type UserData struct {
	UserID                   string
	UserType                 string
	Password                 []byte
	FirstName                string
	LastName                 string
	RentalAddress            AddressData
	OwnedPropertyAddressList []AddressData
	BillingAddress           AddressData
	PaymentList              []PaymentData
	ServiceRequestList       []ServiceRequestData
	NotificationList         []NotificationData
	LandLordID               string
	RentalPaymentAmt         int64
	RentDueDate              string //epoch
	LateFeeRate              string //in percentage maybe make this an int
	LegalDocuments           []DocumentData
	Email                    string
	PhoneNumber              string
}

var User UserInterface = UserInterfaceImpl{}

type UserInterface interface {
	GetUser(email string) (UserData, error)
	GetUserList() ([]UserData, error)
	UpdateUserList(updatedList []UserData) error
}

type UserInterfaceImpl struct {
}

func (impl UserInterfaceImpl) GetUser(email string) (UserData, error) {

	var userList []UserData
	var theUser UserData

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

func (impl UserInterfaceImpl) GetUserList() ([]UserData, error) {

	var userList []UserData

	bytes, err := ioutil.ReadFile(constant.USERFILE)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(bytes, &userList)

	return userList, nil
}

func (impl UserInterfaceImpl) UpdateUserList(updatedList []UserData) error {
	bytes, err := json.Marshal(updatedList)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(constant.USERFILE, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
