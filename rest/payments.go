package rest

import (
	"net/http"
	"strings"
	"time"
	"strconv"
	"encoding/json"
	"context"

	"github.com/braintree-go/braintree-go"
	"github.com/gorilla/mux"

	"github.com/jykwon91/project/util/constant"
	"github.com/jykwon91/project/db"
	"github.com/jykwon91/project/util/email"
	"github.com/jykwon91/project/util/appError"
	"github.com/jykwon91/project/util/logger"
)

func GetPaymentOverview(resp http.ResponseWriter, req *http.Request) *appError.AppError {

        vars := mux.Vars(req)
        landLordID := vars["landLordID"]
        type PaymentOverview struct {
                CurrentPayPeriod string
                CurrentAmountDue int64
                TotalLateAmount  int64
                TotalLateFees    int64
                TotalDue         int64
        }

        claims, err := AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
        if err != nil {
                return &appError.AppError{err, "Getting payment overview failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
        }

        email, err := GetEmailFromClaims(claims)
        if err != nil {
                return &appError.AppError{err, "Getting payment overview failed. Server down. Please contact customer support or try again later.", "Failed to get email", 500}
        }

        userList, err := db.User.GetUserList()
        if err != nil {
                return &appError.AppError{err, "Getting payment overview failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
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
                                if !strings.EqualFold(payment.Status, constant.PAID) && strings.EqualFold(payment.TenantID, tenantID) {
                                        epoch, err := strconv.ParseInt(payment.DueDate, 10, 64)
                                        if err != nil {
                                                return &appError.AppError{err, "Getting payment overview failed. Server down. Please contact customer support or try again later.", "Failed to convert epoch", 500}
                                        }
                                        date := time.Unix(epoch, 0)
                                        if date.Month() == now.Month() {
                                                paymentOverview.CurrentPayPeriod = now.Month().String() + " " + strconv.Itoa(now.Year())
                                                paymentOverview.CurrentAmountDue = paymentOverview.CurrentAmountDue + payment.Amount
                                                paymentOverview.TotalDue = paymentOverview.TotalDue + payment.Amount
                                        } else if strings.EqualFold(payment.Status, constant.LATE) {
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
                return &appError.AppError{err, "Getting payment overview failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
        }

        resp.Header().Set("Content-Type", "application/json")
        resp.WriteHeader(200)
        resp.Write(bytes)

        logger.Logger(nil, resp, req)
        return nil
}

func GetPaymentList(resp http.ResponseWriter, req *http.Request) *appError.AppError {

        claims, err := AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
        if err != nil {
                return &appError.AppError{err, "Getting payment list failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
        }

        //TODO: Function to get user type from claims
        var userType string
        for k, v := range claims {
                if strings.EqualFold(k, "userType") {
                        userType = v.(string)
                }
        }

        email, err := GetEmailFromClaims(claims)
        if err != nil {
                return &appError.AppError{err, "Getting payment list failed. Server down. Please contact customer support or try again later.", "Failed to get email", 500}
        }

        userList, err := db.User.GetUserList()
        if err != nil {
                return &appError.AppError{err, "Getting payment list failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
        }

        var paymentList []db.PaymentData
        var tenantID string
        var landLordID string

        for _, user := range userList {
                if strings.EqualFold(userType, constant.TENANT) && strings.EqualFold(user.Email, email) {
                        tenantID = user.UserID
                        landLordID = user.LandLordID
                        break
                } else if strings.EqualFold(userType, constant.LANDLORD) && strings.EqualFold(user.Email, email) {
                        paymentList = user.PaymentList
                        break
                }
        }

        if strings.EqualFold(userType, constant.TENANT) {
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
                return &appError.AppError{err, "Getting service request list failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
        }

        resp.Header().Set("Content-Type", "application/json")
        resp.WriteHeader(200)
        resp.Write(bytes)

	logger.Logger(nil, resp, req)
	return nil
}

func TenantPayment(resp http.ResponseWriter, req *http.Request) *appError.AppError {
        vars := mux.Vars(req)
        tokenKey := vars["tokenKey"]
        type ReqBody struct {
                PaymentID  string
                LandLordID string
                TenantID   string
                Amount     int64
        }

        var reqBody ReqBody
        err := ReadReqBody(req, &reqBody)
        if err != nil {
                return &appError.AppError{err, "Payment failed. Please contact customer support or try again later.", "Failed to read request", 500}
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
                Type:               "sale",
                Amount:             braintree.NewDecimal(reqBody.Amount, 2),
                PaymentMethodNonce: tokenKey,
        })
        if err != nil {
                return &appError.AppError{err, "Payment failed. Server down. Please contact customer support.", "Payment failed", 403}
        }

        pendingPayment, err := db.Payment.UpdatePayment(reqBody.PaymentID, reqBody.LandLordID, reqBody.TenantID, reqBody.Amount, t.Id, string(t.PaymentInstrumentType))
        if err != nil {
                return &appError.AppError{err, "Creating payment history failed. Server down. Please contact customer support.", "Creating Payment failed", 500}
        }

        err = db.Payment.AddToPendingPaymentList(pendingPayment)
        if err != nil {
                return &appError.AppError{err, "Updating payment failed. Please contact customer support.", "Adding payment to pending payment list failed", 500}
        }

        err = email.EmailPaymentConfirmation(pendingPayment)
        if err != nil {
                return &appError.AppError{err, "Updating payment failed. Please contact customer support.", "Adding payment to pending payment list failed", 500}
        }

        logger.Logger(nil, resp, req)
	return nil
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
