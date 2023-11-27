// Code generated by go-enum DO NOT EDIT.
// Version: 0.5.7
// Revision: bf63e108589bbd2327b13ec2c5da532aad234029
// Build Date: 2023-07-25T23:27:55Z
// Built By: goreleaser

package trees

import (
	"errors"
	"fmt"
)

const (
	// DirectionNoDir is a Direction of type No_dir.
	DirectionNoDir Direction = iota
	// DirectionLeft is a Direction of type Left.
	DirectionLeft
	// DirectionRight is a Direction of type Right.
	DirectionRight
)

var ErrInvalidDirection = errors.New("not a valid Direction")

const _DirectionName = "no_dirleftright"

var _DirectionMap = map[Direction]string{
	DirectionNoDir: _DirectionName[0:6],
	DirectionLeft:  _DirectionName[6:10],
	DirectionRight: _DirectionName[10:15],
}

// String implements the Stringer interface.
func (x Direction) String() string {
	if str, ok := _DirectionMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Direction(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Direction) IsValid() bool {
	_, ok := _DirectionMap[x]
	return ok
}

var _DirectionValue = map[string]Direction{
	_DirectionName[0:6]:   DirectionNoDir,
	_DirectionName[6:10]:  DirectionLeft,
	_DirectionName[10:15]: DirectionRight,
}

// ParseDirection attempts to convert a string to a Direction.
func ParseDirection(name string) (Direction, error) {
	if x, ok := _DirectionValue[name]; ok {
		return x, nil
	}
	return Direction(0), fmt.Errorf("%s is %w", name, ErrInvalidDirection)
}