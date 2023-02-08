package git

import "sync"

func memoize[T any](callback func() (T, error)) func() (T, error) {
	var (
		result T
		err    error
	)

	once := &sync.Once{}

	return func() (T, error) {
		once.Do(func() {
			result, err = callback()
		})

		return result, err
	}
}
