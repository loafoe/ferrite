package types

import (
	"github.com/philips-software/go-hsdp-api/iron"
)

type Schedule struct {
	iron.Schedule
	CodeID string
}
