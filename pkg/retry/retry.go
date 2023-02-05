package retry

import "time"

func Retry(fn func() error, times int, sleepDuration time.Duration) (err error) {
	for err = fn(); err != nil && times > 1; times, err = times-1, fn() {
		time.Sleep(sleepDuration)
	}
	return err
}
