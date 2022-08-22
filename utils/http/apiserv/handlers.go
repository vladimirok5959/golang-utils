package apiserv

import (
	"regexp"
	"sync"
)

type Handlers struct {
	sync.RWMutex
	list map[*regexp.Regexp]Handler
}
