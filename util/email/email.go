package email

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/jykwon91/project/util/constant"
	"github.com/jykwon91/project/db"
	"github.com/jykwon91/project/util/convert"
	"github.com/jykwon91/project/util/logger"
)

func EmailCompletedPaymentConfirmation(completedPaymentList []db.PaymentData) {

	email, err := ioutil.ReadFile(constant.EMAILFILE)
	if err != nil {
		logger.AltLogger(err.Error())
	}

	emailPass, err := ioutil.ReadFile(constant.EMAILPASSFILE)
	if err != nil {
		logger.AltLogger(err.Error())
	}
	var userList []db.UserData

	bytes, err := ioutil.ReadFile(constant.USERFILE)
	if err != nil {
		logger.AltLogger(err.Error())
	}
	json.Unmarshal(bytes, &userList)

	t := time.Now()
	year := strconv.Itoa(t.Year())
	month := t.Month().String()
	//day := strconv.Itoa(t.Day())
	//theDate := month + " " + day + " " + year

	for _, theUser := range userList {
		for _, payment := range completedPaymentList {
			if strings.EqualFold(payment.TenantID, theUser.UserID) {
				m := gomail.NewMessage()
				m.Embed("/home/jkwon/Git/project/signature.jpg")
				m.SetHeader("From", string(email))
				m.SetHeader("To", theUser.Email)
				m.SetHeader("Subject", "Payment confirmation for "+month+" "+year)
				m.SetBody("text/html", `
                                                <p> This is to confirm your payment has finished processing.</p>                
                                `)

				d := gomail.NewDialer("smtp.gmail.com", 587, string(email), string(emailPass))

				// Send the email to Bob, Cora and Dan.
				if err := d.DialAndSend(m); err != nil {
					logger.AltLogger(err.Error())
				}
			}
		}
	}
}

func EmailPaymentConfirmation(paymentInfo db.PaymentData) error {

	email, err := ioutil.ReadFile(constant.EMAILFILE)
	if err != nil {
		return err
	}

	emailPass, err := ioutil.ReadFile(constant.EMAILPASSFILE)
	if err != nil {
		return err
	}
	var userList []db.UserData

	bytes, err := ioutil.ReadFile(constant.USERFILE)
	if err != nil {
		return err
	}
	json.Unmarshal(bytes, &userList)

	t := time.Now()
	year := strconv.Itoa(t.Year())
	month := t.Month().String()
	day := strconv.Itoa(t.Day())
	theDate := month + " " + day + " " + year

	amountDue := convert.IntToDollar(paymentInfo.Amount)
	amountPaid := convert.IntToDollar(paymentInfo.Amount)

	for _, theUser := range userList {
		if strings.EqualFold(theUser.UserID, paymentInfo.TenantID) {
			name := theUser.FirstName + " " + theUser.LastName
			m := gomail.NewMessage()
			//m.Embed("/home/jkwondd/Git/project/signature.jpg")
			m.SetHeader("From", string(email))
			m.SetHeader("To", theUser.Email)
			m.SetHeader("Subject", "Rent Receipt for "+month+" "+year)
			m.SetBody("text/html", `
                                <p>Thank you for your payment! This is a confirmation email. We will process your payment within 24 hours.</p><br>
                                <table width='600' style='border:1px solid #333'>
                                <tbody>
                                        <tr><td align='left'><b>Transaction number:</b>`+paymentInfo.PaymentID+`</td></tr>
                                        <tr><td align='left'><b>Brain Tree Transaction number:</b>`+paymentInfo.BTTransactionID+`</td></tr>
                                        <tr><td align='left'><b>Name:</b> `+name+`</td></tr>
                                        <tr><td align='left'><b>Date:</b> `+theDate+`</td></tr>
                                        <tr><td align='left'><b>Transaction Type:</b>Card</td></tr>
                                        <tr>
                                                <td align='center'>
                                                        <table align='center' width='300' border='0' cellspacing='0' cellpadding='0' style='border:1px solid #ccc; padding:10px 0px 10px 10px'>
                                                                <tr>
                                                                        <td><b>Category:</b></td>
                                                                        <td>`+paymentInfo.Category+`</td>
                                                                </tr>
                                                                <tr>
                                                                        <td><b>Amount Due:</b></td>
                                                                        <td>$`+amountDue+`</td>
                                                                </tr>
                                                                <tr>
                                                                        <td><b>Amount Paid:</b></td>
                                                                        <td>$`+amountPaid+`</td>
                                                                </tr>
                                                                <tr>
                                                                        <td><b>Received by:</b></td>
                                                                        <td>Jason Kwon</td>
                                                                </tr>
                                                        </table>
                                        </tr>
                                        <br>
                                </tbody>`)

			d := gomail.NewDialer("smtp.gmail.com", 587, string(email), string(emailPass))

			if err := d.DialAndSend(m); err != nil {
				return err
			}
		}
	}

	return nil
}

func EmailRentDueNotification() {
	email, err := ioutil.ReadFile(constant.EMAILFILE)
	if err != nil {
		logger.AltLogger(err.Error())
	}

	emailPass, err := ioutil.ReadFile(constant.EMAILPASSFILE)
	if err != nil {
		logger.AltLogger(err.Error())
	}
	var userList []db.UserData

	bytes, err := ioutil.ReadFile(constant.USERFILE)
	if err != nil {
		logger.AltLogger(err.Error())
	}
	json.Unmarshal(bytes, &userList)

	t := time.Now()
	year := strconv.Itoa(t.Year())
	month := t.Month().String()
	day := strconv.Itoa(t.Day())
	theDate := month + " " + day + " " + year

	for _, theUser := range userList {
		if strings.EqualFold(theUser.UserType, constant.TENANT) {
			m := gomail.NewMessage()
			m.SetHeader("From", string(email))
			m.SetHeader("To", theUser.Email)
			m.SetHeader("Subject", "Rent due for "+month+" "+year)
			m.SetBody("text/html", `
                                <p>This is a reminder that rent is due today(`+theDate+`)</p><br>
                                <p>Log into www.rentalmgmt.co to pay.</p>
                        `)

			d := gomail.NewDialer("smtp.gmail.com", 587, string(email), string(emailPass))
			if err := d.DialAndSend(m); err != nil {
				logger.AltLogger(err.Error())
			}
		}
	}

}

func EmailServiceReq(landlord db.UserData, serviceReq db.ServiceRequestData) {
	email, err := ioutil.ReadFile(constant.EMAILFILE)
	if err != nil {
		logger.AltLogger(err.Error())
	}

	emailPass, err := ioutil.ReadFile(constant.EMAILPASSFILE)
	if err != nil {
		logger.AltLogger(err.Error())
	}

	address := serviceReq.RentalAddress.Street + " " + serviceReq.RentalAddress.Zipcode
	now := time.Now().Format("2006-01-02 15:04:05")

	m := gomail.NewMessage()
	//m.Embed("/home/jkwondd/Git/project/signature.jpg")
	m.SetHeader("From", string(email))
	m.SetHeader("To", landlord.Email)
	m.SetHeader("Subject", "Service request from "+serviceReq.TenantName)
	m.SetBody("text/html", `
                <p>Date: `+now+`</p><br>
                <p>Address: `+address+`</p><br>
                <p>Message: `+serviceReq.Message+`</p><br>
                `)

	d := gomail.NewDialer("smtp.gmail.com", 587, string(email), string(emailPass))

	if err := d.DialAndSend(m); err != nil {
		logger.AltLogger(err.Error())
	}
}
