package opstate

import "time"

// LoadingSpinner has a string representation that dynamically changes over time to
// display a visual loading spinner
type LoadingSpinner interface {
	String() string
}

// DotLoadingSpinner is a LoadingSpinner that uses brail dots
type DotLoadingSpinner struct {
	state       int
	lastChanged time.Time
}

var loadingRunes = []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}

func (l DotLoadingSpinner) String() string {
	now := time.Now()
	ms := now.UnixNano() / int64(time.Millisecond)
	r := loadingRunes[(ms/int64(100))%int64(len(loadingRunes))]
	return string(r)
}
