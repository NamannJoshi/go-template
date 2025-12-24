package utils

import (
	"fmt"
	"strings"
	"time"
)

type DateOnly struct {
	Time time.Time
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	fmt.Print(t)
	d.Time = t
	return nil
}
