package model

import "sort"

type KeySet map[ComponentKey]struct{}

func NewKeySet(keys ...ComponentKey) KeySet {
	set := make(KeySet, len(keys))
	for _, key := range keys {
		set[key] = struct{}{}
	}
	return set
}

func (s KeySet) Has(key ComponentKey) bool {
	_, ok := s[key]
	return ok
}

func (s KeySet) Add(key ComponentKey) {
	s[key] = struct{}{}
}

func (s KeySet) And(other KeySet) KeySet {
	out := make(KeySet)
	for key := range s {
		if other.Has(key) {
			out[key] = struct{}{}
		}
	}
	return out
}

func (s KeySet) Or(other KeySet) KeySet {
	out := make(KeySet, len(s)+len(other))
	for key := range s {
		out[key] = struct{}{}
	}
	for key := range other {
		out[key] = struct{}{}
	}
	return out
}

func (s KeySet) Minus(other KeySet) KeySet {
	out := make(KeySet)
	for key := range s {
		if !other.Has(key) {
			out[key] = struct{}{}
		}
	}
	return out
}

func (s KeySet) Xor(other KeySet) KeySet {
	return s.Minus(other).Or(other.Minus(s))
}

func (s KeySet) Sorted(level Level) []ComponentKey {
	keys := make([]ComponentKey, 0, len(s))
	for key := range s {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Name != keys[j].Name {
			return keys[i].Name < keys[j].Name
		}
		if level == Level1 {
			return false
		}
		return keys[i].Version < keys[j].Version
	})
	return keys
}
