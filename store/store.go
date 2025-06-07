package store

import "sync"

var syncMap sync.Map

func Save(key, value any) {
	syncMap.Store(key, value)
}

func Delete(key any) {
	syncMap.Delete(key)
}

func DeleteMany(keys ...any) {
	for _, key := range keys {
		syncMap.Delete(key)
	}
}

func Clear() {
	syncMap.Clear()
}

func Load(key any) (value any, ok bool) {
	return syncMap.Load(key)
}

func LoadDefault(key, defaultValue any) (value any) {
	if data, ok := syncMap.Load(key); ok {
		return data
	}
	return defaultValue
}
