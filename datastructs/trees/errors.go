package trees

import "fmt"

type AlreadyExistsError[K any] struct {
	GivenKey K
}

func (err AlreadyExistsError[K]) Error() string {
	return fmt.Sprintf("item for key %+v already exists", err.GivenKey)
}

func (err AlreadyExistsError[K]) Is(target error) bool {
	_, ok := target.(AlreadyExistsError[K])
	return ok
}

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
