package db

var Address AddressInterface = AddressInterfaceImpl{}

type AddressInterface interface {
}

type AddressData struct {
	AddressID    string
	Street       string
	Zipcode      string
	City         string
	State        string
	PropertyType string
}

type AddressInterfaceImpl struct {
}
