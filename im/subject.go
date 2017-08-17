package im

import (
	"sync"
)

//A Topic means a public
//discussion
type Topic struct {
	sync.RWMutex
	ID   [16]byte
	Name string
}
