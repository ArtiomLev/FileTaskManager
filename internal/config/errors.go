package config

import "errors"

var (
	ErrConfigNotFound     = errors.New("config not found")
	ErrCannotCreateConfig = errors.New("failed to create config")
	ErrCannotCreateDir    = errors.New("failed to create config directory")
	ErrCannotReadConfig   = errors.New("cannot read config")
	ErrCannotParseConfig  = errors.New("cannot parse config")
)
