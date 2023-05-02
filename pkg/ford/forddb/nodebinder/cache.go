package nodebinder

import (
	"reflect"
	"strings"
	"sync"
)

var globalFieldNameCache = newFieldNameCache()

type fieldNameCache struct {
	sync.RWMutex

	cache map[reflect.Type]map[string]reflect.StructField
}

func newFieldNameCache() *fieldNameCache {
	return &fieldNameCache{
		cache: map[reflect.Type]map[string]reflect.StructField{},
	}
}

func (c *fieldNameCache) GetAll(t reflect.Type) map[string]reflect.StructField {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fields := c.cache[t]

	if fields == nil {
		c.RWMutex.Lock()
		defer c.RWMutex.Unlock()

		if c.cache[t] == nil {
			fields = c.extractFields(t)
			c.cache[t] = fields
		}
	}

	return fields
}

func (c *fieldNameCache) Get(t reflect.Type, name string) reflect.StructField {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fields := c.cache[t]

	if fields == nil {
		c.RWMutex.Lock()
		defer c.RWMutex.Unlock()

		if c.cache[t] == nil {
			fields = c.extractFields(t)
			c.cache[t] = fields
		}
	}

	return fields[name]
}

func (c *fieldNameCache) extractFields(t reflect.Type) map[string]reflect.StructField {
	fields := map[string]reflect.StructField{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		name := field.Name
		hasTag := false

		if tag := field.Tag.Get("json"); tag != "" {
			parts := strings.Split(tag, ",")
			name = parts[0]
			hasTag = true
		}

		if name == "-" {
			continue
		}

		isEmbedded := field.Anonymous

		if hasTag && name != "" {
			isEmbedded = false
		}

		if isEmbedded {
			embeddedFields := c.extractFields(field.Type)

			for name, embeddedName := range embeddedFields {
				embeddedName.Index = append([]int{i}, embeddedName.Index...)
				fields[name] = embeddedName
			}
		} else {
			fields[name] = field
		}
	}

	return fields
}
