package configreader

import (
	"errors"
)

var (
	ErrInvalidCmdArgs     = errors.New("invalid command line args")
	ErrConfigFileNotFound = errors.New("config file not found")
	ErrConfigFileParse    = errors.New("config file parse error")
	ErrConfigAddress      = errors.New("invalid address")
)
