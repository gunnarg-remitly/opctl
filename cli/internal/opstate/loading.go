package opstate

import "time"

// Loading has a string representation that dynamically changes over time to
// display a visual loading spinner
type Loading struct {
	state       int
	lastChanged time.Time
}

var loadingRunes = []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}

func (l *Loading) String() string {
	now := time.Now()
	ms := now.UnixNano() / int64(time.Millisecond)
	r := loadingRunes[(ms/int64(100))%int64(len(loadingRunes))]
	return string(r)
}
