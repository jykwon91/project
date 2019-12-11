package task

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/braintree-go/braintree-go"
	"github.com/google/uuid"

	"github.com/jykwon91/project/db"
	"github.com/jykwon91/project/util/constant"
	"github.com/jykwon91/project/util/convert"
	"github.com/jykwon91/project/util/email"
	"github.com/jykwon91/project/util/logger"
)

func InitTasks() {
	go daily()
}

func daily() {
	checkRentDue()
	checkRentLate()
	checkPendingPayments()
}

func checkRentDue() {
	for t := range time.NewTicker(constant.TWELVE_HOURS).C {
		//TEST: create test payment obj
		//for t := range time.NewTicker(8 * time.Second).C {
		fmt.Println(t)

		userList, err := db.User.GetUserList()
		if err != nil {
			logger.AltLogger(err.Error())
		}

		//TODO: Add check for charge in same month
		//TODO: Change status to late if needed
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

func checkRentLate() {
	for t := range time.NewTicker(constant.TWELVE_HOURS).C {

		userList, err := db.User.GetUserList()
		if err != nil {
			logger.AltLogger(err.Error())
		}

		for i, user := range userList {
			if strings.EqualFold(user.UserType, constant.TENANT) {
				for j, payment := range user.PaymentList {
					if !strings.EqualFold(payment.Status, constant.PAID) && t.Day() > 1 {
						daysLate := calculateDaysLate()

						if daysLate > 7 {
							//TODO: Send eviction notice if greater than 7 days
							//      Automatically print letter or send via mail and email
						}

						userList[i].PaymentList[j] = calculateLateFee(payment, user.RentalPaymentAmt, daysLate)
					}
				}
			}
		}

		err = db.User.UpdateUserList(userList)
		if err != nil {
			logger.AltLogger(err.Error())
		}
	}
}

func calculateDaysLate() int64 {
	_, _, day := time.Now().Date()
	daysLate := day - 1
	return int64(daysLate)
}

func calculateLateFee(payment db.PaymentData, rentalPaymentAmt int64, daysLate int64) db.PaymentData {

	payment.Status = constant.LATE

	days := convert.DaysInMonth(time.Now().Month().String())
	rentPerDay := rentalPaymentAmt / days

	lateFee := 7500 + (rentPerDay * daysLate)

	payment.LateFeeAmount = lateFee
	payment.Amount = rentalPaymentAmt + lateFee

	return payment
}

func checkPendingPayments() {
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
