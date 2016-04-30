package events
import (
	"time"
	"sync"
)

type InMemoryESRepo struct {
	store map[string][]EventHappened
	*sync.RWMutex
}

func (e *events) AddInMemoryEventRepo() {
	e.eventRepo = NewInMemESRepo()
	return
}

func NewInMemESRepo() *InMemoryESRepo {
	return &InMemoryESRepo{
		make(map[string][]EventHappened),
		new(sync.RWMutex),
	}
}

func (repo *InMemoryESRepo) Set(eventName string, data interface{}) error {
	eh := EventHappened{time.Now(), eventName, data}
	repo.RWMutex.Lock()
	repo.store[eventName] = append(repo.store[eventName], eh)
	repo.RWMutex.Unlock()
	return nil
}

func (repo *InMemoryESRepo) GetAll() (map[string][]EventHappened, error) {
	repo.RWMutex.RLock()
	defer repo.RWMutex.RUnlock()
	return repo.store, nil
}

func (repo *InMemoryESRepo) Contains(eventName string) bool {
	var e []EventHappened
	var ok bool
	repo.RWMutex.RLock()
	if e, ok = repo.store[eventName]; !ok {
		repo.RWMutex.RUnlock()
		return ok;
	}
	repo.RWMutex.RUnlock()
	return len(e) > 0
}
