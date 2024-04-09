package rui

import (
	"sort"
	"strings"
	"sync"
)

var mutexProperties sync.Mutex

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
	properties map[string]any
}

func (properties *propertyList) init() {
	properties.properties = map[string]any{}
}

func (properties *propertyList) Get(tag string) any {
	return properties.getRaw(strings.ToLower(tag))
}

func (properties *propertyList) getRaw(tag string) any {
	mutexProperties.Lock()
	defer mutexProperties.Unlock()
	if value, ok := properties.properties[tag]; ok {
		return value
	}
	return nil
}

func (properties *propertyList) setRaw(tag string, value any) {
	mutexProperties.Lock()
	defer mutexProperties.Unlock()
	properties.properties[tag] = value
}

func (properties *propertyList) Remove(tag string) {
	mutexProperties.Lock()
	defer mutexProperties.Unlock()
	delete(properties.properties, strings.ToLower(tag))
}

func (properties *propertyList) remove(tag string) {
	mutexProperties.Lock()
	defer mutexProperties.Unlock()
	delete(properties.properties, tag)
}

func (properties *propertyList) Clear() {
	mutexProperties.Lock()
	defer mutexProperties.Unlock()
	properties.properties = map[string]any{}
}

func (properties *propertyList) AllTags() []string {
	mutexProperties.Lock()
	defer mutexProperties.Unlock()
	tags := make([]string, 0, len(properties.properties))
	for t := range properties.properties {
		tags = append(tags, t)
	}
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
