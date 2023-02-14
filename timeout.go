package dockerdb

import "time"

func SetMaxWaitTime(sec time.Duration) {
	maxWaitTime = sec
}
