package db

import (
	"strings"
	"time"
	"strconv"
	"io/ioutil"
	"encoding/json"

	"github.com/jykwon91/project/util/constant"
)

var Payment PaymentInterface = PaymentInterfaceImpl{}

type PaymentInterface interface {
	UpdatePayment(paymentID string, landlordID string, tenantID string, amount int64, btTransactionID string, paymentMethod string) (PaymentData, error)
	AddToPendingPaymentList(pendingPayment PaymentData) error
}

type PaymentInterfaceImpl struct {
}

// could this be tied to Document
type PaymentData struct {
	TenantName      string
	RentalAddress   string
	PaymentID       string
	LandLordID      string
	TenantID        string
	BTTransactionID string
	Category        string //rent, service, repairs, utility, etc
	PaymentMethod   string
	Amount          int64  //cents
	LateFeeAmount   int64  //cents
	Status          string //open, processing, paid, late
	PaidDate        string
	DueDate         string
	Description     string
}

func (impl PaymentInterfaceImpl) UpdatePayment(paymentID string, landLordID string, tenantID string, amount int64, btTransactionID string, paymentMethod string) (PaymentData, error) {

        userList, err := User.GetUserList()
        if err != nil {
                return PaymentData{}, err
        }

        var pendingPayment PaymentData
        for i, user := range userList {
                if strings.EqualFold(user.UserID, landLordID) {
                        for j, payment := range user.PaymentList {
                                if strings.EqualFold(paymentID, payment.PaymentID) {
                                        if amount == payment.Amount {
                                                now := time.Now()
                                                secs := now.Unix()
                                                date := strconv.FormatInt(secs, 10)
                                                userList[i].PaymentList[j].Status = constant.PROCESSING
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

        err = User.UpdateUserList(userList)
        if err != nil {
                return PaymentData{}, err
        }

        return pendingPayment, nil
}

func (impl PaymentInterfaceImpl) AddToPendingPaymentList(pendingPayment PaymentData) error {

        var pendingPaymentList []PaymentData

        bytes, err := ioutil.ReadFile(constant.PENDINGPAYMENTSFILE)
        if err != nil {
                return err
        }
        json.Unmarshal(bytes, &pendingPaymentList)

        pendingPaymentList = append(pendingPaymentList, pendingPayment)

        bytes, err = json.Marshal(pendingPaymentList)
        if err != nil {
                return err
        }

        err = ioutil.WriteFile(constant.PENDINGPAYMENTSFILE, bytes, 0644)
        if err != nil {
                return err
        }

        return nil
}
