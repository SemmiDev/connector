package g_learning_connector_test

import (
	"fmt"
	"github.com/alitto/pond/v2"
	"testing"
)

func TestWorker(t *testing.T) {
	pool := pond.NewPool(25)

	for i := 1; i <= 100; i++ {
		pool.Submit(func() {
			fmt.Println("hello world " + fmt.Sprintf("%d", i))
		})
	}

	pool.StopAndWait()
}
