package rui

import (
	"sort"
	"strings"
	"sync"
)

// Properties interface of properties map
type Properties interface {
	// Get returns a value of the property with name defined by the argument.
	// The type of return value depends on the property. If the property is not set then nil is returned.
	Get(tag string) any
	getRaw(tag string) any
	// Set sets the value (second argument) of the property with name defined by the first argument.
	// Return "true" if the value has been set, in the opposite case "false" are returned and
	// a description of the error is written to the log
	Set(tag string, value any) bool
	setRaw(tag string, value any)
	// Remove removes the property with name defined by the argument
	Remove(tag string)
	// Clear removes all properties
	Clear()
	// AllTags returns an array of the set properties
	AllTags() []string
}

type propertyList struct {
	properties sync.Map
}

func (properties *propertyList) init() {
}

func (properties *propertyList) Get(tag string) any {
	return properties.getRaw(strings.ToLower(tag))
}

func (properties *propertyList) getRaw(tag string) any {
	if value, ok := properties.properties.Load(tag); ok {
		return value
	}
	return nil
}

func (properties *propertyList) setRaw(tag string, value any) {
	properties.properties.Store(tag, value)
}

func (properties *propertyList) Remove(tag string) {
	properties.properties.Delete(strings.ToLower(tag))
}

func (properties *propertyList) remove(tag string) {
	properties.properties.Delete(tag)
}

func (properties *propertyList) Clear() {
	properties.properties = sync.Map{}
}

func (properties *propertyList) AllTags() []string {
	tags := make([]string, 0)
	properties.properties.Range(func(key, value any) bool {
		tags = append(tags, key.(string))
		return true
	})
	sort.Strings(tags)
	return tags
}

func parseProperties(properties Properties, object DataObject) {
	count := object.PropertyCount()
	for i := 0; i < count; i++ {
		if node := object.Property(i); node != nil {
			switch node.Type() {
			case TextNode:
				properties.Set(node.Tag(), node.Text())

			case ObjectNode:
				properties.Set(node.Tag(), node.Object())

			case ArrayNode:
				properties.Set(node.Tag(), node.ArrayElements())
			}
		}
	}
}
