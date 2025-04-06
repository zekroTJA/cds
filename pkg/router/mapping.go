package router

// A lot of this code is more or less directly copied
// from Go's net/http package.
// See: https://github.com/golang/go/blob/master/src/net/http/routing_tree.go
// License: https://github.com/golang/go/blob/master/LICENSE

const maxSlice = 8

type entry[K comparable, V any] struct {
	key K
	val V
}

type mapping[K comparable, V any] struct {
	slicePairs []entry[K, V]
	mapPairs   map[K]V
}

func (t *mapping[K, V]) add(key K, val V) {
	if t.mapPairs == nil && len(t.slicePairs) < maxSlice {
		t.slicePairs = append(t.slicePairs, entry[K, V]{key: key, val: val})
		return
	}

	if t.mapPairs == nil {
		t.mapPairs = make(map[K]V, maxSlice+1)
		for _, entry := range t.slicePairs {
			t.mapPairs[entry.key] = entry.val
		}
		t.slicePairs = nil
	}

	t.mapPairs[key] = val
}

func (t *mapping[K, V]) get(key K) (val V, ok bool) {
	if t.mapPairs != nil {
		val, ok = t.mapPairs[key]
		return val, ok
	}

	for _, entry := range t.slicePairs {
		if entry.key == key {
			return entry.val, true
		}
	}

	return val, false
}
