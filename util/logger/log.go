package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jykwon91/project/util/constant"
	"github.com/jykwon91/project/util/appError"
)

func Logger(e *appError.AppError, resp http.ResponseWriter, req *http.Request) {
	var logMessage string
	origin := strings.Replace(req.Header.Get("Origin"), "http://", "", -1)
	now := time.Now().Format("2006-01-02 15:04:05")
	if e != nil {
		log.Printf("[HTTP %d][%s][%s] - %s - ERROR: %s: %s", e.Code, req.Method, origin, req.RequestURI, e.Message, e.Error)
		logMessage = "[" + now + "][HTTP " + strconv.FormatUint(e.Code, 10) + "][" + req.Method + "][" + origin + "] - " + req.RequestURI + " - ERROR: " + e.Message + ": " + e.Error.Error() + "\n"
	} else {
		log.Printf("[HTTP 200][%s][%s] - %s ", req.Method, origin, req.RequestURI)
		logMessage = "[" + now + "][HTTP 200][" + req.Method + "][" + origin + "] - " + req.RequestURI + "\n"
	}

	f, err := os.OpenFile(constant.LOGFILE, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		AltLogger(err.Error())
	}

	_, err = f.WriteString(logMessage)
	if err != nil {
		AltLogger(err.Error())
	}

	f.Close()
}

func AltLogger(errStr string) {
	logMessage := "ERROR: " + errStr
	f, err := os.OpenFile(constant.LOGFILE, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Can't print to server log")
	}

	_, err = f.WriteString(logMessage)
	if err != nil {
		fmt.Printf("Can't print to server log")
	}

	f.Close()
}
