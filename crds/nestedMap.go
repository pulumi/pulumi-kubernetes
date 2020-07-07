package main

// NestedMap contains an arbitrarily nested structure. Has nice methods that
// allow for quickly traversing down keys/values, assuming 'object' is of a
// suitable type.
type NestedMap struct {
	object interface{}
}

// get returns a new NestedMap, with a new object of the value of object[key].
// Assumes that object can be casted to map[interface{}]interface{}.
func (n NestedMap) get(key interface{}) NestedMap {
	o := n.object
	m := (o).(map[interface{}]interface{})
	value := m[key]
	return NestedMap{value}
}

// Get returns a new NestedMap, with a new object of the value of following
// each given key. n.Get(key1, key2, ..., keyn) is equivalent to calling
// n.get(key1).get(key2).get(keyn). Assumes that the object can be casted to
// map[interface{}]interface{} for each intermediate value.
func (n NestedMap) Get(keys ...interface{}) NestedMap {
	for _, key := range keys {
		n = n.get(key)
	}
	return n
}

// Index returns a new NestedMap, with a new object of the value of
// object[index]. Assumes that object can be casted to []interface{}.
func (n NestedMap) Index(index int) NestedMap {
	o := n.object
	m := o.([]interface{})
	value := m[index]
	return NestedMap{value}
}

// Int returns the object casted to an integer.
func (n NestedMap) Int() int {
	return (n.object).(int)
}

// String returns the object casted to a string.
func (n NestedMap) String() string {
	return (n.object).(string)
}

// Map returns the object casted to an arbitrary map[interface{}]interface{}.
func (n NestedMap) Map() map[interface{}]interface{} {
	return (n.object).(map[interface{}]interface{})
}

// StringArray returns the object casted to an array of strings.
func (n NestedMap) StringArray() []string {
	o := n.object
	genericArray := o.([]interface{})
	stringArray := make([]string, len(genericArray))
	for i, value := range genericArray {
		stringArray[i] = value.(string)
	}
	return stringArray
}
