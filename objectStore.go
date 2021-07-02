package cachettl

import (
	"reflect"
	"sync"
	"time"
)

type ObjectStore struct {
	store map[string]*objectWithTTL
	mutex sync.Mutex
}

func NewObjectStore() *ObjectStore {
	return &ObjectStore{
		store: make(map[string]*objectWithTTL),
	}
}

func (s *ObjectStore) Add(key string, data interface{}, ttl int64) error {
	if len(key) == 0 {
		return ErrKeyIsBlanc
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	obj := s.store[key]
	if obj != nil {
		return ErrObjExist
	}

	s.store[key] = &objectWithTTL{
		Data:       reflect.ValueOf(data),
		Type:       reflect.TypeOf(data),
		Ttl:        ttl,
		CreateTime: time.Now().Truncate(time.Millisecond),
	}

	return nil
}

func (s *ObjectStore) Get(key string, outObj interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	obj := s.store[key]

	if obj == nil {
		return ErrObjNotFound
	}

	if !obj.checkValid() {
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

func (s *ObjectStore) Update(key string, data interface{}, ttl int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	obj := s.store[key]
	if obj == nil {
		return ErrObjNotFound
	}

	obj.Data = reflect.ValueOf(data)
	obj.Ttl = ttl
	obj.CreateTime = time.Now().Truncate(time.Millisecond)

	s.store[key] = obj

	return nil
}

func (s *ObjectStore) Delete(key string) {
	s.mutex.Lock()

	delete(s.store, key)

	s.mutex.Unlock()
}
