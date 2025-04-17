package trees

import "fmt"

type NotFoundError[K any] struct {
	GivenKey K
}

func (err NotFoundError[K]) Error() string {
	return fmt.Sprintf("item for key %+v does not exists", err.GivenKey)
}

func (err NotFoundError[K]) Is(target error) bool {
	_, ok := target.(NotFoundError[K])
	return ok
}
