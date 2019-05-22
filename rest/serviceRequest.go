package rest

import (
	"net/http"
	"strings"
	"encoding/json"

	"github.com/jykwon91/project/db"
	"github.com/jykwon91/project/util/constant"
	"github.com/jykwon91/project/util/appError"
	"github.com/jykwon91/project/util/logger"
)

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

