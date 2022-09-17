package errors

import "github.com/bjatkin/bear"

// Error Types
var (
	FileNotFound = bear.NewType("File Not Found")
)

// Exit Codes
const (
	BuildFailed = iota + 1
)
