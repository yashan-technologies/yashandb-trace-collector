package errdef

import (
	"errors"

	"ytc/i18n"
)

var (
	ErrEndLessStart        = errors.New("start time should be less than end time")
	ErrStartShouldLessCurr = errors.New("start time should be less current time")
)

type ErrGreaterMaxDur struct {
	MaxDuration string
}

func NewGreaterMaxDur(max string) *ErrGreaterMaxDur {
	return &ErrGreaterMaxDur{MaxDuration: max}
}

func (e ErrGreaterMaxDur) Error() string {
	return i18n.TWithData("err.greater_max_dur", map[string]interface{}{"MaxDuration": e.MaxDuration})
}

type ErrLessMinDur struct {
	MinDuration string
}

func NewLessMinDur(min string) *ErrLessMinDur {
	return &ErrLessMinDur{MinDuration: min}
}

func (e ErrLessMinDur) Error() string {
	return i18n.TWithData("err.less_min_dur", map[string]interface{}{"MinDuration": e.MinDuration})
}
