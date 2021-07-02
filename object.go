package cachettl

import (
	"reflect"
	"time"
)

type objectWithTTL struct {
	Data       reflect.Value
	Type       reflect.Type
	Ttl        int64
	CreateTime time.Time
}

func (o *objectWithTTL) checkValid() bool {
	if time.Now().Truncate(time.Millisecond).Sub(o.CreateTime) < time.Duration(o.Ttl)*time.Second {
		return true
	} else {
		return false
	}
}
