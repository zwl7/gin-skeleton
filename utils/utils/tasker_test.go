package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestTasker(t *testing.T) {
	_tasker := NewTasker(6)
	time.Sleep(time.Duration(1) * time.Second)
	for _i := 0; _i < 1000; _i++ {
		_n := _i
		go func(n int) {
			//time.Sleep(time.Duration(n%3) * time.Second)
			for _j := 0; _j < 100; _j++ {
				_m := _j
				_tasker.AddTask(_m, func(k, l int) func() {
					return func() {
						//time.Sleep(time.Duration(200) * time.Millisecond)
						if k == 999 {
							fmt.Println("========== TestTasker ===========", k, l)
						}
					}
				}(n, _m))
			}
		}(_n)
	}
	time.Sleep(time.Duration(25*100*200) * time.Millisecond)
}
