package apiserv

import "strconv"

type Param struct {
	Value string
}

func (p Param) Integer() int64 {
	if v, err := strconv.ParseInt(p.Value, 10, 64); err == nil {
		return v
	}
	return 0
}

func (p Param) String() string {
	return p.Value
}
