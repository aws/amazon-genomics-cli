package main

import (
	"fmt"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/format"
	"github.com/rs/zerolog/log"
)

func setFormatter(f formatVars) string {
	configClient, err := newConfigClient()
	if err != nil {
		log.Error().Err(err)
		return ""
	}
	if f.format == "" {
		f.format = defaultFormat
		configFormat, err := configClient.GetFormat()
		if err != nil {
			log.Error().Err(err)
		} else {
			f.format = configFormat
		}
	}
	if err := ValidateFormat(format.FormatterType(f.format)); err != nil {
		fmt.Println(err.Error())
	}
	format.SetFormatter(format.FormatterType(f.format))
	return f.format
}

func ValidateFormat(f format.FormatterType) error {
	if err := f.ValidateFormatter(); err != nil {
		return err
	}
	return nil
}
