package constant

import "time"

const (
	LANDLORD            = "landlord"
	TENANT              = "tenant"
	RENT                = "rent"
	PAID                = "paid"
	OPEN                = "open"
	PROCESSING          = "processing"
	LATE                = "late"
	ERROR               = "error"
	LOGFILE             = "/home/jkwon/Git/project/log/server.log"
	TOKENFILE           = "/home/jkwon/Git/project/database/Tokens"
	PENDINGPAYMENTSFILE = "/home/jkwon/Git/project/database/PendingPayments"
	STATELISTFILE       = "/home/jkwon/Git/project/database/StateList"
	USERFILE            = "/home/jkwon/Git/project/database/Users"
	EMAILFILE           = "/home/jkwon/Git/project/etc/businessEmail"
	EMAILPASSFILE       = "/home/jkwon/Git/project/etc/businessEmailPass"
	TOKENPASSFILE       = "/home/jkwon/Git/project/etc/tokenpass"
	TWELVE_HOURS        = 43200 * time.Second
)
