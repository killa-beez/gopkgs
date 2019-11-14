package pool

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	p := New(1, 2)
	assert.Equal(t, 2, p.numWorkers)
	assert.NotNil(t, p.ch)
}

type arrSetter struct {
	arr   []int
	index int
	value int
}

func (as arrSetter) Perform(ctx context.Context) {
	as.arr[as.index] = as.value
}

func TestPerformDoesNotRunBeforeStart(t *testing.T) {
	arr := make([]int, 3)

	p := New(1, 1)
	p.Add(arrSetter{arr, 0, 0})
	p.Wait()

	assert.Equal(t, arr, []int{0, 0, 0})
}

func TestPerform(t *testing.T) {
	arr := make([]int, 3)

	p := New(3, 3)
	p.Start(context.Background())
	p.Add(arrSetter{arr, 0, 1})
	p.Add(arrSetter{arr, 1, 2})
	p.Add(arrSetter{arr, 2, 3})
	p.Wait()

	assert.Equal(t, arr, []int{1, 2, 3})
}
