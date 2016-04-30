package events

import (
	"testing"
	log "src/github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)


func TestShouldStoreAndRetrieveEvents(t *testing.T) {

	log.SetLevel(log.DebugLevel)
	e := New()
	e.AddInMemoryEventRepo()
	e.On("event", func(a interface{}) error {
		var err error
		pingD = true
		return err
	})
	var data interface{}
	data = 1
	e.Raise("event", data)
	e.Wait()
	assert.True(t, e.GetEventRepo().Contains("event"))
}


func TestShouldStoreAndRetrieveEvensts(t *testing.T) {

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
	assert.Nil(t, e.GetEventRepo())
}