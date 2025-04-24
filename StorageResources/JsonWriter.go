package StorageResources

import "sync"

// this type is for development and testing only
type JsonWriter struct {
	mutex sync.Mutex
}
