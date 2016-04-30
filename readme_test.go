package events

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

var wasError = false

func TestRunReadmeExample(t *testing.T) {

	errorListener := func(e ListenerError){
		wasError = true
	}

	helloListener := func(a interface{}) error {
		fmt.Print("Hello ", a)
		return nil
	}

	e := New()
	e.On("hailed", helloListener)
	e.OnError(errorListener)

	e.Raise("hailed", "World")
	e.Wait()
	assert.False(t, wasError)
}