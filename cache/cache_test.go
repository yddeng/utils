/*
 *
 */
package cache_test

import (
	"fmt"
	"github.com/tagDong/dutil/cache"
	"testing"
)

func TestNew(t *testing.T) {
	c := cache.New(nil, 2, string("1"), cache.LFU)
	c.Set("1", "1")
	c.Set("1", "2")
	c.Set("3", "3")
	fmt.Println(c.Size())
	v, ok := c.Get("1")
	fmt.Println(v, ok)
	c.Set("4", "4")
	fmt.Println(c.Size())
	v, ok = c.Get("1")
	fmt.Println(v, ok)
	v, ok = c.Get("3")
	fmt.Println(v, ok)
}
