package utils

import "time"

func RetryConnectToDataBase(fn func() error, maxAttempts int, delay time.Duration) (err error) {
	for maxAttempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			maxAttempts--
			continue
		}
		return nil
	}
	return
}
