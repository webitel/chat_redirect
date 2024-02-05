package cache

import (
	"context"
	"errors"
)

type CacheValue struct {
	value any
}

type CacheStore interface {
	Get(ctx context.Context, key string) (*CacheValue, error)
	Set(ctx context.Context, key string, value any, expiresAfter int64) error
	Delete(ctx context.Context, key string) error
}

func (v *CacheValue) String() (string, error) {
	if v.value == nil {
		return "", errors.New("value is nil")
	}
	if v, ok := v.value.(string); ok {
		return v, nil
	} else {
		return "", errors.New("unable to convert value")
	}

}

func (v *CacheValue) Raw() any {
	return v.value
}
func (v *CacheValue) Set(value any) error {
	if value != nil {
		v.value = value
		return nil
	} else {
		return errors.New("accepted value is nil")
	}
}

func NewCacheValue(value any) (*CacheValue, error) {
	var cv CacheValue
	err := cv.Set(value)
	if err != nil {
		return nil, err
	}
	return &cv, err
}
