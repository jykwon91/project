package db

var Document DocumentInterface = DocumentInterfaceImpl{}

type DocumentInterface interface {

}

type DocumentInterfaceImpl struct {
}

type DocumentData struct {
	DocumentID    string
	DocumentType  string //receipt, contract, contact, personal
	DocumentBytes []byte
}
