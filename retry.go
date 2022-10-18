package abs

import "time"

func retry[TResult any](cb func() (TResult, error), delay time.Duration) TResult {
	for {
		result, err := cb()
		if err != nil {
			time.Sleep(delay)
			continue
		}

		return result
	}
}
