package rest

import (
	"net/http"
	"strings"
	"encoding/json"
	"time"
	"strconv"

	"github.com/google/uuid"

	"github.com/jykwon91/project/db"
	"github.com/jykwon91/project/util/constant"
	"github.com/jykwon91/project/util/appError"
	"github.com/jykwon91/project/util/logger"
	"github.com/jykwon91/project/util/email"
)

func SendServiceRequest(resp http.ResponseWriter, req *http.Request) *appError.AppError {
        type ReqService struct {
                Message       string
                TenantName    string
                RentalAddress db.AddressData
        }

        _, err := AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
        if err != nil {
                return &appError.AppError{err, "Sending service request failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
        }

        userList, err := db.User.GetUserList()
        if err != nil {
                return &appError.AppError{err, "Sending service request failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
        }

        var reqService ReqService
        err = ReadReqBody(req, &reqService)
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

func GetServiceRequestList(resp http.ResponseWriter, req *http.Request) *appError.AppError {
        claims, err := AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
        if err != nil {
                return &appError.AppError{err, "Getting service request list failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
        }

        var userType string
        for k, v := range claims {
                if strings.EqualFold(k, "userType") {
                        userType = v.(string)
                }
        }

        email, err := GetEmailFromClaims(claims)
        if err != nil {
                return &appError.AppError{err, "Getting service request list failed. Server down. Please contact customer support or try again later.", "Failed to get email", 500}
        }

        userList, err := db.User.GetUserList()
        if err != nil {
                return &appError.AppError{err, "Getting service request list failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
        }

        var serviceRequestList []db.ServiceRequestData

        if strings.EqualFold(userType, "tenant") {
                var tenantRentalAddress db.AddressData
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
        } else if strings.EqualFold(userType, constant.LANDLORD) {
                for i, user := range userList {
                        if strings.EqualFold(user.Email, email) {
                                serviceRequestList = userList[i].ServiceRequestList
                        }
                }
        }

        bytes, err := json.Marshal(serviceRequestList)
        if err != nil {
                return &appError.AppError{err, "Getting service request list failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
        }

        resp.Header().Set("Content-Type", "application/json")
        resp.WriteHeader(200)
        resp.Write(bytes)

        logger.Logger(nil, resp, req)
        return nil
}

func UpdateServiceRequest(resp http.ResponseWriter, req *http.Request) *appError.AppError {

        claims, err := AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
        if err != nil {
                return &appError.AppError{err, "Sending service request failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
        }

        email, err := GetEmailFromClaims(claims)
        if err != nil {
                return &appError.AppError{err, "Sending service request failed. Server down. Please contact customer support or try again later.", "Failed to get email", 403}
        }

        userList, err := db.User.GetUserList()
        if err != nil {
                return &appError.AppError{err, "Sending service request failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
        }

        var updateServiceReqObj db.ServiceRequestData
        err = ReadReqBody(req, &updateServiceReqObj)
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

