package cmd

import "time"

const (
	BaseCommandGet      = "get"
	BaseCommandImport   = "import"
	BaseCommandDescribe = "describe"
	CoordinatorFlagName = "coordinator"
	TimeOut             = 10 * time.Second
)
