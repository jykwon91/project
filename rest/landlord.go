package rest

import (
	"net/http"
	"strings"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/jykwon91/project/db"
	"github.com/jykwon91/project/util/logger"
	"github.com/jykwon91/project/util/appError"
)

func GetTenantList(resp http.ResponseWriter, req *http.Request) *appError.AppError {

        claims, err := AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
        if err != nil {
                return &appError.AppError{err, "Getting tenant list failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
        }

        email, err := GetEmailFromClaims(claims)
        if err != nil {
                return &appError.AppError{err, "Getting tenant list failed. Server down. Please contact customer support or try again later.", "Failed to get email", 500}
        }

        userList, err := db.User.GetUserList()
        if err != nil {
                return &appError.AppError{err, "Getting tenant list failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
        }

        landLordID := ""
        for _, user := range userList {
                if strings.EqualFold(user.Email, email) {
                        landLordID = user.UserID
                        break
                }
        }

        if strings.EqualFold(landLordID, "") {
                return &appError.AppError{err, "Getting tenant list failed. Server down. Please contact customer support or try again later.", "Could not find land lord", 500}
        }

        var tenantList []db.UserData
        for _, user := range userList {
                if strings.EqualFold(user.LandLordID, landLordID) {
                        tenantList = append(tenantList, user)
                }
        }

        bytes, err := json.Marshal(tenantList)
        if err != nil {
                return &appError.AppError{err, "Getting all landlord properties failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
        }

        resp.Header().Set("Content-Type", "application/json")
        resp.WriteHeader(200)
        resp.Write(bytes)

        logger.Logger(nil, resp, req)
        return nil
}

func GetAllLandLordProperties(resp http.ResponseWriter, req *http.Request) *appError.AppError {

        claims, err := AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
        if err != nil {
                return &appError.AppError{err, "Getting all landlord properties failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
        }

        email, err := GetEmailFromClaims(claims)
        if err != nil {
                return &appError.AppError{err, "Getting all landlord properties failed. Server down. Please contact customer support or try again later.", "Failed to get email", 403}
        }

        userList, err := db.User.GetUserList()
        if err != nil {
                return &appError.AppError{err, "Getting all landlord properties failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
        }

        var addressList []db.AddressData
        for _, user := range userList {
                if strings.EqualFold(user.Email, email) {
                        addressList = user.OwnedPropertyAddressList
                        break
                }
        }

        bytes, err := json.Marshal(addressList)
        if err != nil {
                return &appError.AppError{err, "Getting all landlord properties failed. Server down. Please contact customer support or try again later", "Failed to marshal response body", 500}
        }

        resp.Header().Set("Content-Type", "application/json")
        resp.WriteHeader(200)
        resp.Write(bytes)

        logger.Logger(nil, resp, req)
        return nil
}

func RegisterLandlordProperty(resp http.ResponseWriter, req *http.Request) *appError.AppError {

        claims, err := AuthenticateTokenAndReturnClaims(req.Header.Get("Authorization"))
        if err != nil {
                return &appError.AppError{err, "Registering landlord property failed. Server down. Please contact customer support or try again later.", "Failed to authenticate user", 403}
        }

        var address db.AddressData
        err = ReadReqBody(req, &address)
        if err != nil {
                return &appError.AppError{err, "Registering landlord property failed. Please contact customer support or try again later.", "Failed to read request", 500}
        }

        address.AddressID = uuid.New().String()
        address.PropertyType = "rental"

        email, err := GetEmailFromClaims(claims)
        if err != nil {
                return &appError.AppError{err, "Registering landlord property failed. Server down. Please contact customer support or try again later.", "Failed to get email", 500}
        }

        userList, err := db.User.GetUserList()
        if err != nil {
                return &appError.AppError{err, "Registering landlord property failed. Server down. Please contact customer support or try again later.", "Failed to get user list", 500}
        }

        for i, user := range userList {
                if strings.EqualFold(user.Email, email) {
                        userList[i].OwnedPropertyAddressList = append(userList[i].OwnedPropertyAddressList, address)
                        break
                }
        }

        err = db.User.UpdateUserList(userList)
        if err != nil {
                return &appError.AppError{err, "Registering landlord property failed. Server down. Please contact customer support or try again later.", "Failed to update database", 500}
        }

        logger.Logger(nil, resp, req)
        return nil
}
