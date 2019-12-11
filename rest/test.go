package rest

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jykwon91/project/db"
	"github.com/jykwon91/project/util/appError"
	"github.com/jykwon91/project/util/constant"
	"github.com/jykwon91/project/util/logger"
)

func TestCreatePayment(resp http.ResponseWriter, req *http.Request) *appError.AppError {
	claims, err := AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
	if err != nil {
		return &appError.AppError{err, "Failed to authenticate user", "Failed to authenticate user", 403}
	}

	email, err := GetEmailFromClaims(claims)
	if err != nil {
		return &appError.AppError{err, "Failed to authenticate user", "Failed to authenticate user", 403}
	}

	userList, err := db.User.GetUserList()
	if err != nil {
		return &appError.AppError{err, "Failed to get userlist from database.", "Failed to get user list", 500}
	}

	var tenantName, rentalAddress string
	var foundUser db.UserData
	for _, user := range userList {
		if strings.EqualFold(user.Email, email) {
			foundUser = user
			tenantName = user.FirstName + " " + user.LastName
			rentalAddress = user.RentalAddress.Street + "," + user.RentalAddress.Zipcode + "," + user.RentalAddress.City + "," + user.RentalAddress.State
		}
	}

	now := time.Now()
	secs := now.Unix()
	date := strconv.FormatInt(secs, 10)

	payment := db.PaymentData{TenantName: tenantName, RentalAddress: rentalAddress, LandLordID: foundUser.LandLordID, TenantID: foundUser.UserID, BTTransactionID: "", Category: constant.RENT, PaymentMethod: "credit", Status: constant.OPEN, Amount: 100, DueDate: date}

	for i, user := range userList {
		if strings.EqualFold(user.UserID, foundUser.LandLordID) {
			userList[i].PaymentList = append(userList[i].PaymentList, payment)
		}
	}

	err = db.User.UpdateUserList(userList)
	if err != nil {
		return &appError.AppError{err, "Creating test payment failed.", "Failed to add to database", 500}
	}

	logger.Logger(nil, resp, req)
	return nil
}
