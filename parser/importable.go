package parser

import (
	"fmt"
	"slices"
	"strings"
)

type Importable interface {
	Name() string
	Prefix() string
	PrefixedName() string
	Global() bool
	Public() bool
	Equal(compare Importable) bool
}

type ImportableBase struct {
	name   string
	prefix string
	global bool
}

func NewImportableBase(name string, prefix string, global bool) ImportableBase {
	return ImportableBase{
		name,
		prefix,
		global,
	}
}

func (b ImportableBase) Name() string {
	return b.name
}

func (b ImportableBase) Prefix() string {
	return b.prefix
}

func (b ImportableBase) PrefixedName() string {
	name := b.name
	prefix := b.prefix

	if len(prefix) > 0 {
		prefix = fmt.Sprintf("%s_", prefix)

		// Only prefix if it doesn't already have the prefix.
		if !strings.HasPrefix(name, prefix) {
			name = fmt.Sprintf("%s%s", prefix, name)
		}
	}
	return name
}

func (b ImportableBase) Global() bool {
	return b.global
}

func (b ImportableBase) Public() bool {
	return isPublic(b.Name())
}

func (b ImportableBase) Equal(compare Importable) bool {
	return b.PrefixedName() == compare.PrefixedName()
}

type Importables[T Importable] [][]T

func (is *Importables[T]) add(importables ...T) {
	for _, importable := range importables {
		index := -1
		stack := []T{}

		for i, stackTemp := range *is {
			for _, importableTemp := range stackTemp {
				if importableTemp.Equal(importable) {
					index = i
					stack = stackTemp
					break
				}
			}

			if index >= 0 {
				break
			}
		}
		stack = append(stack, importable)

		if index < 0 {
			(*is) = append((*is), stack)
		} else {
			(*is)[index] = stack
		}
	}
}

func (is Importables[T]) find(name string, prefix string) (T, bool) {
	importable, index := is.findIndex(name, prefix)
	return importable, index >= 0
}

func (is Importables[T]) findIndex(name string, prefix string) (T, int) {
	for i, stack := range is {
		lastIndex := len(stack) - 1

		if lastIndex >= 0 {
			importable := stack[lastIndex]

			if importable.Name() == name && (importable.Prefix() == prefix || importable.Prefix() == "") { // TODO: Find out if it's a good idea to check for importable.Prefix() == "".
				return importable, i
			}
		}
	}
	var none T
	return none, -1
}

func (is Importables[T]) clone() Importables[T] {
	clone := Importables[T]{}

	for _, stack := range is {
		clone = append(clone, slices.Clone(stack))
	}
	return clone
}
