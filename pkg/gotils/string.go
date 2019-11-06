package gotils

import (
	"fmt"
	"sync"
)

var (
	mu      sync.RWMutex
	counter = 0
)

// GenerateUniqueString make unique string
func GenerateUniqueString(prefix string) string {
	var str string
	mu.Lock()
	counter = counter + 1
	str = fmt.Sprintf("%s_uniq_%d", prefix, counter)
	mu.Unlock()

	return str
}
