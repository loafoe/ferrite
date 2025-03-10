package types

import (
	"github.com/dip-software/go-dip-api/iron"
)

type Schedule struct {
	iron.Schedule
	CodeID string
}
