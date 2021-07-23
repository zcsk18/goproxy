package core

import (
	"sync"
)

var Pool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024*500)
	},
}