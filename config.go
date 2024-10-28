package ucodesdk

import (
	"time"
)

type Config struct {
	AppId          string
	BaseURL        string
	AuthBaseURL    string
	FunctionName   string
	ProjectId      string
	RequestTimeout time.Duration
}
