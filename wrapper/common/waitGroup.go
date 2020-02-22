package common

import (
	"sync"
)

var s sync.Once

// 全局唯一的 waitGroup 
var wg *sync.WaitGroup

// WaitGroup 获取全局唯一的 waitgroup
func WaitGroup() *sync.WaitGroup {
	s.Do(func() {
		wg = new(sync.WaitGroup)
	})
	return wg
}
