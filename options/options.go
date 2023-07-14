package options

import "fmt"

// Realization of functional option pattern
type Option[T any] func(target *T) error

// Apply given opts.
// Any of aopts must not be nil.
// In case of error on any position, opts execution will be stopped.
func ApplyOptions[T any](target *T, opts ...Option[T]) error {
	for idx, opt := range opts {
		if opt == nil {
			return fmt.Errorf("got nil opt on pos %d", idx)
		}
		if err := opt(target); err != nil {
			return fmt.Errorf("applying opt on pos %d: %w", idx, err)
		}
	}
	return nil
}
