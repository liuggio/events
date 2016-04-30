package events

import (
	"sync"
	"time"
	log "github.com/Sirupsen/logrus"
)

type Dispatcher interface {
	On(on string, fn func(interface{}) error)
	Raise(on string, data interface{})
	Wait()
	Close()
}

type EventHappened struct {
	Time time.Time
	On   string
	Data interface{}
}

type EventRepo interface  {
	Set(string, interface{}) error
	GetAll() (map[string][]EventHappened, error)
	Contains(string) bool
}

type listener struct {
	on  string
	fn  func(interface{}) error
}

type ListenerError struct {
	Err error
	On  string
	In  interface{}
	Fn  func(interface{}) error
}

func newError(err error, on string, in interface{}, fn func(interface{}) error) ListenerError {
	return ListenerError{err, on, in, fn}
}

type ErrorListenerFn func(err ListenerError)

func newListener(on string, fn func(interface{}) error) *listener {
	return &listener{
		on: on,
		fn: fn,
	}
}

func (l *listener) execute(e *events, d interface{}) {
	log.Info("Executing on:", l.on)

	e.wg.Add(1)
	go func(e *events) {
		err := l.fn(d)
		defer e.wg.Done()

		if err != nil {
			e.errCh <- newError(err, l.on, d, l.fn)
		}
	}(e)
}

type events struct {
	// Mutex to prevent race conditions within the Emitter.
	*sync.RWMutex
	// Map of listener on event name
	liteners  map[string][]listener
	errorFn   ErrorListenerFn
	wg        sync.WaitGroup
	errCh     chan ListenerError

	eventRepo EventRepo
}

func New() (e *events) {
	e = new(events)
	e.RWMutex = new(sync.RWMutex)
	e.liteners = make(map[string][]listener)
	e.errCh = make(chan ListenerError)

	go func(errCh chan ListenerError) {
		for {
			select {
				case err := <- errCh:
					if err.Err != nil {
						e.raiseError(err)
					}
				}
		}
	}(e.errCh)

	log.Info("created")
	return
}

func NewWithErrListener(el ErrorListenerFn) (e *events){
	e = New()
	e.OnError(el)
	return
}

func (e *events) OnError(el ErrorListenerFn) {
	e.errorFn = el
	return
}

func (e *events) AddEventRepo(es EventRepo) {
	e.eventRepo = es
	return
}

func (e *events) GetEventRepo() EventRepo {
	return e.eventRepo
}


func (e *events) On(on string, fn func(interface{}) error) {
	l := newListener(on, fn)
	e.RWMutex.Lock()
	e.liteners[on] = append(e.liteners[on], *l)
	e.RWMutex.Unlock()
	log.Info("Added listener on:", on)
}

func (e *events) Raise(on string, data interface{}) {

	if e.eventRepo != nil {
		e.eventRepo.Set(on, data)
	}

	for _, l := range e.liteners[on] {
		l.execute(e, data)
	}
}

func (e *events) Wait() {
	e.wg.Wait()
}

func (e *events) Close() {
	close(e.errCh)
}

func (e *events) raiseError(err ListenerError) {
	if e.errorFn != nil {
		e.wg.Add(1)
		defer e.wg.Done()
		e.errorFn(err)
	}
}