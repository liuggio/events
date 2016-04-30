package events

import (
	"testing"
	"errors"
	"time"
	log "src/github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var pingD = false
var pingW = false
var pingE = false

func TestShouldDispatch(t *testing.T) {

	log.SetLevel(log.DebugLevel)
	e := New()
	e.On("event", func(a interface{}) error {
		var err error
		pingD = true
		return err
	})
	var data interface{}
	data = 1
	e.Raise("event", data)

	e.Wait()
	assert.True(t, pingD)
}

func TestShouldWait(t *testing.T) {

	log.SetLevel(log.DebugLevel)
	e := New()
	e.On("event", func(a interface{}) error {
		var err error
		pingW=true
		time.Sleep(2*time.Second)
		return err
	})
	var data interface{}
	data = "in"
	e.Raise("event", data)
	e.Wait()
	assert.True(t, pingW)
}

func TestOnErrorShouldRaiseErrorEvent(t *testing.T) {

	errList := func(e ListenerError){
		pingE = true
		log.Errorf("error on listener: error[%s] on[%s] input[%v], fn[%v]", e.Err, e.On, e.In, e.Fn)
	}

	log.SetLevel(log.DebugLevel)
	e := New()
	e.OnError(errList)
	e.On("event", func(a interface{}) error {
		return errors.New("error found :)")
	})
	var data interface{}
	data = "in"
	e.Raise("event", data)
	e.Wait()

	assert.True(t, pingE)
}
