package apiserv

import "strconv"

type Param struct {
	value string
}

func (p Param) Integer() int64 {
	if v, err := strconv.ParseInt(p.value, 10, 64); err == nil {
		return v
	}
	return 0
}

func (p Param) String() string {
	return p.value
}
