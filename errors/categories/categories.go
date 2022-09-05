package categories

import (
	"strconv"
)

// A Category is an unsigned 32-bit error code
type Category uint32

type Categorized interface {
	error
	CodeInt()
	CodeString()
	//Category() Category
}

const (
	// Unknown
	Unknown Category = 0

	// TemporaryError
	Temporary Category = 1

	// SystemError
	System Category = 2

	// DataError
	Data Category = 3

	//_maxCode = 5
)

func (c Category) String() string {
	switch c {
	case Unknown:
		return "Unknown"
	case Temporary:
		return "Temporary"
	case System:
		return "System"
	case Data:
		return "Data"
	default:
		return "Code(" + strconv.FormatInt(int64(c), 10) + ")"
	}
}
