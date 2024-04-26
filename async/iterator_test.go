package async

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	ch := make(chan string, 10)
	ch <- "hello"
	ch <- "hello"
	ch <- "hello"
	ch <- "hello"
	close(ch)
	for {
		mm, ok := <-ch
		fmt.Println(mm, ok)
		if !ok {
			break
		}
	}
}
