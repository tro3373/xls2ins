package cmd

import (
	"errors"
	"path/filepath"
	"regexp"
)

type Config struct {
	BookConfigs []BookConfig
}

type BookConfig struct {
	BookNameRegExp string
	SheetConfigs   []SheetConfig
}

type SheetConfig struct {
	SheetName  string
	StartRow   int
	SqlFormat  string
	SqlArgCols []string
}

func NewConfigError(msg string) error {
	return errors.New("ConfigError: " + msg)
}

func (c Config) FindBookConfig(filePath string) *BookConfig {
	bookName := filepath.Base(filePath)
	for _, e := range c.BookConfigs {
		exp := regexp.MustCompile(e.BookNameRegExp)
		if exp.MatchString(bookName) {
			return &e
		}
	}
	return nil
}

func (s SheetConfig) validate() error {
	if s.SheetName == "" {
		return NewConfigError("SheetName is required")
	}
	if s.StartRow == 0 {
		return NewConfigError("StartRow is required")
	}
	if s.SqlFormat == "" {
		return NewConfigError("SqlFormat is required")
	}
	if len(s.SqlArgCols) == 0 {
		return NewConfigError("SqlArgCols is required")
	}
	return nil
}
