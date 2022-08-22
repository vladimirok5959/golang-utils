package apiserv

import (
	"net/http"
	"sync"
)

type Params struct {
	sync.RWMutex
	list map[*http.Request][]Param
}
