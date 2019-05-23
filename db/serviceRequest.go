package db

var ServiceRequest ServiceRequestInterface = ServiceRequestInterfaceImpl{}

type ServiceRequestInterface interface {

}

type ServiceRequestInterfaceImpl struct {
}

type ServiceRequestData struct {
	Status        string
	RequestID     string
	RequestTime   string
	StartTime     string
	CompletedTime string
	Message       string
	RentalAddress AddressData
	TenantName    string
}
