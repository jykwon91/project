package db

var Notification NotificationInterface = NotificationInterfaceImpl{}

type NotificationInterface interface {

}

type NotificationInterfaceImpl struct {
}

type NotificationData struct {
	NotificationID string
	CreatedOn      string //epoch time
	Message        string
	From           string
}
