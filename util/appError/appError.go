package appError

type AppError struct {
	Error            error
	ServerLogMessage string
	Message          string
	Code             uint64
}
