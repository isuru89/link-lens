package analyzer

import "fmt"

const (
	ErrorInvalidUrl        = "InvalidUrl"
	RemoteFetchError       = "RemoteFetchError"
	UnsuccessfulStatusCode = "UnsuccessfulStatusCode"
	InvalidContentType     = "InvalidContentType"
)

type AnalysisError struct {
	ErrorCode string

	Cause error
}

func (e *AnalysisError) Error() string {
	return fmt.Sprintf("[%s] %s", e.ErrorCode, e.Cause.Error())
}
