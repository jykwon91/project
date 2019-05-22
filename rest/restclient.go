package rest

import (
	"net/http"
	"fmt"
	"strings"
	"io/ioutil"
	"encoding/json"
	"log"

	"github.com/jykwon91/project/util/constant"
	"github.com/jykwon91/project/util/logger"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/jykwon91/project/util/appError"
)

type appHandler func(http.ResponseWriter, *http.Request) *appError.AppError

func InitRestClient() {
	router := mux.NewRouter()

        router.Handle("/users/authenticate", appHandler(AuthenticateUser))
        router.Handle("/tenant/payment/overview/{landLordID}", appHandler(GetPaymentOverview))
        router.Handle("/users/landlord/all", appHandler(GetLandLordList))
        router.Handle("/tenant/pay/{tokenKey}", appHandler(TenantPayment))
        router.Handle("/users/service/all", appHandler(GetServiceRequestList))
/*
        router.Handle("/users/register", appHandler(registerUser))
        router.Handle("/users/all", appHandler(getAllUsers))
        router.Handle("/users/notification/all", appHandler(getAllNotifications))
        router.Handle("/users/currentUser", appHandler(getCurrentUser))
        router.Handle("/users/update", appHandler(updateUser))
        router.Handle("/users/payment/all", appHandler(getPaymentList))

        router.Handle("/stateList", appHandler(getStateList))
        router.Handle("/users/landlord/property/register", appHandler(registerLandlordProperty))
        router.Handle("/landlord/property/all", appHandler(getAllLandLordProperties))
        router.Handle("/landlord/notification", appHandler(sendNotification))
        router.Handle("/landlord/service/request/update", appHandler(updateServiceRequest))
        router.Handle("/landlord/tenant/all", appHandler(getTenantList))
        router.Handle("/tenant/service/request", appHandler(sendServiceRequest))

        //router.Handle("/landlord/tenant/delete/{tenantID}", appHandler(deleteTenant))
        //router.Handle("/users/notification/delete/{notificationID}", appHandler(deleteNotification))
*/

        router.Handle("/test/create/payment", appHandler(TestCreatePayment))

        router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
        c := cors.New(cors.Options{
                AllowedOrigins:   []string{"http://10.0.0.152:8081", "http://localhost", "http://localhost:8080", "http://192.168.1.125", "http://192.168.1.125:8080", "http://rentalmgmt.co:8080", "http://rentalmgmt.co", "http://localhost:8081", "http://www.rentalmgmt.co"},
                AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
                AllowedHeaders:   []string{"Accept", "Accept-Language", "Content-Type", "Authorization", "Bearer"},
                AllowCredentials: true,
                Debug:            false,
        })

	router.Use(AuthMiddleware)

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", handler))
}

func (fn appHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
        if e := fn(resp, req); e != nil {
                logger.Logger(e, resp, req)
                resp.Header().Set("Content-Type", "application/json")
                resp.WriteHeader(int(e.Code))
                fmt.Fprintf(resp, `{"type":"error","message":"%s", "code":"%d"}`, e.ServerLogMessage, e.Code)
        }
}

// Middleware function which will be called for each request
func AuthMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
                if !strings.EqualFold(req.URL.Path, "/users/authenticate") && !strings.EqualFold(req.URL.Path, "/users/register") && !strings.EqualFold(req.URL.Path, "/stateList") && strings.EqualFold(req.URL.Path, "/tenant/pay") {
                        var tokenList []string
                        token := req.Header.Get("Authorization")

                        bytes, _ := ioutil.ReadFile(constant.TOKENFILE)
                        json.Unmarshal(bytes, &tokenList)

                        var authenticated bool
                        for _, tmpToken := range tokenList {
                                if strings.EqualFold(tmpToken, token) {
                                        authenticated = true
                                }
                        }
                        if authenticated {
                                next.ServeHTTP(resp, req)
                        } else {
                                http.Error(resp, "\"Forbidden. You do not have permission to view this content.\"", http.StatusForbidden)
                        }
                } else {
                        next.ServeHTTP(resp, req)
                }
        })
}
