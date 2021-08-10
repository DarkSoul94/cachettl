package cachettl

import (
	"context"
	"reflect"
	"sync"
	"time"
)

type ObjectStore struct {
	store    map[string]*objectWithTTL
	mutex    sync.Mutex
	shutdown context.CancelFunc
}

func NewObjectStore(cleanPeriod time.Duration) *ObjectStore {
	ctx, cancel := context.WithCancel(context.Background())

	newStore := &ObjectStore{
		store:    make(map[string]*objectWithTTL),
		shutdown: cancel,
	}

	go newStore.cleaner(cleanPeriod, ctx)

	return newStore
}

func (s *ObjectStore) Close() {
	s.shutdown()
}

func (s *ObjectStore) Add(key string, data interface{}, ttl int64) error {
	if len(key) == 0 {
		return ErrKeyIsBlank
	}

	newObj := &objectWithTTL{
		Data:       reflect.ValueOf(data),
		Type:       reflect.TypeOf(data),
		Ttl:        ttl,
		CreateTime: time.Now().Truncate(time.Millisecond),
	}

	s.mutex.Lock()
	s.store[key] = newObj
	s.mutex.Unlock()

	return nil
}

func (s *ObjectStore) Get(key string, outObj interface{}) error {
	s.mutex.Lock()
	obj, ok := s.store[key]
	s.mutex.Unlock()
	if !ok {
		return ErrObjNotFound
	}

	if !obj.checkValid() {
		s.Delete(key)
		return ErrObjNotValid
	}

	v := reflect.ValueOf(outObj)
	if v.Elem().Type() == obj.Type {
		v.Elem().Set(obj.Data)
	} else {
		return ErrInvalidType
	}

	return nil
}

func (s *ObjectStore) Delete(key string) {
	s.mutex.Lock()

	delete(s.store, key)

	s.mutex.Unlock()
}

func (s *ObjectStore) cleaner(period time.Duration, ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			s.mutex.Lock()
			store := s.store
			s.mutex.Unlock()

			for key, obj := range store {
				if !obj.checkValid() {
					s.Delete(key)
				}
			}
		}
		time.Sleep(period)
	}
}
