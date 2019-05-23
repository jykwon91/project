package rest

import (
	"net/http"
	"encoding/json"
	"strings"
	"time"
	"strconv"

	"github.com/google/uuid"

	"github.com/jykwon91/project/db"
	"github.com/jykwon91/project/util/appError"
	"github.com/jykwon91/project/util/logger"
	//"github.com/jykwon91/project/util/email"
)

func SendNotification(resp http.ResponseWriter, req *http.Request) *appError.AppError {

        type ReqNotification struct {
                Message     string
                AddressList []string
        }

        _, err := AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
        if err != nil {
                return &appError.AppError{err, "Sending notification failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
        }

        userList, err := db.User.GetUserList()
        if err != nil {
                return &appError.AppError{err, "Sending notification failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
        }

        var reqNotification ReqNotification
        err = ReadReqBody(req, &reqNotification)
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

func GetAllNotifications(resp http.ResponseWriter, req *http.Request) *appError.AppError {
        claims, err := AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
        if err != nil {
                return &appError.AppError{err, "Getting notifications failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
        }

        email, err := GetEmailFromClaims(claims)
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
