package web

import (
	"sync"

	"github.com/gin-gonic/gin"
)

// WaitGroupWrapper is a handler wrapper which adds a handler to a sync.WaitGroup
func WaitGroupWrapper(wg *sync.WaitGroup) gin.HandlerFunc {
	return func(c *gin.Context) {
		wg.Add(1)
		defer wg.Done()

		c.Next()
	}
}
