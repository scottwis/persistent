package persistent

import (
	"constraints"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

type Comparable[T any] interface {
	Equal(rhs T) bool
	comparable
}

type Ordered[T any] interface {
	Comparable[T]
	Less(rhs T) bool
}

type Key[T constraints.Ordered] struct {
	Key T
}

func (lhs Key[T]) Less(rhs Key[T]) bool {
	return lhs.Key < rhs.Key
}

func (lhs Key[T]) Equal(rhs Key[T]) bool {
	return lhs.Key == rhs.Key
}

func AsKey[T constraints.Ordered](key T) Key[T] {
	return Key[T]{Key: key}
}

func (lhs Key[T]) MarshalText() ([]byte, error) {
	//Key[T] only works for types satisfying the ordered constraint
	return []byte(fmt.Sprint(lhs.Key)), nil
}

func (lhs Key[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(lhs.Key)
}
func (lhs *Key[T]) UnmarshalText(text []byte) error {
	v := reflect.ValueOf(&lhs.Key)
	if v.Kind() == reflect.String {
		*(*string)(unsafe.Pointer(&lhs.Key)) = string(text)
	}

	_, err := fmt.Sscanf(string(text), "%v", &lhs.Key)
	return err
}

func (lhs *Key[T]) UnmarshalJSON(data []byte) error {
	s := strings.TrimSpace(string(data))
	if s != "" && s[0] == '"' {
		var s string
		err := json.Unmarshal(data, &s)
		if err != nil {
			return err
		}
		return lhs.UnmarshalText([]byte(s))
	}
	return json.Unmarshal(data, &lhs.Key)
}

// Iterator defines an interface for an iterator of a persistent data structure.
//
// It is safe to use multiple iterators over the same data structure concurrently. A single
// Iterator instance, however, is not concurrency safe. Sharing an iterator instance between go-routines requires
// explicit synchronization.
//
// Example:
//    var iter = tree.Iterator()
//    for iter.Next() {
//        fmt.Println(iter.Current())
//    }
type Iterator[T any] interface {
	//Next advances the iterator by one position, and returns true iif the new position is valid.
	Next() bool

	//Current returns the value at the current iterator position. Will panic if the iterator position is invalid.
	Current() T
}
