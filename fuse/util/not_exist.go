package util

import (
	"sync"
	"time"
)

// The file node in NotExistManager
type notExistFile struct {
	name string
	// timeout in ns
	negativeTime int64
}

// isNegative : if negative not timeout return true
// else return false
func (file *notExistFile) isNegative() bool {

	if file.negativeTime > time.Now().UnixNano() {
		return true
	}

	return false
}

// NotExistManager used to cache the not exist file path
// Because file system maybe will scan the mountpoint,
// and will look up a lot of not exist path
// TODO: should add a goroutine to delete timeout node?
type NotExistManager struct {
	NegativeTimeout int

	dict map[string]*notExistFile

	lk sync.RWMutex
}

// Init the function to initialize the map
func (manager *NotExistManager) Init(NegativeTimeout int) {
	manager.dict = make(map[string]*notExistFile)
	manager.NegativeTimeout = NegativeTimeout
}

// Set the not exist filepath, negativeTimeout is how long will timeout.
// The unit is seconds
func (manager *NotExistManager) Set(filepath string, negativeTimeout int) {

	seconds := time.Duration(negativeTimeout) * time.Second

	negativeTime := time.Now().UnixNano() + seconds.Nanoseconds()

	file := &notExistFile{name: filepath, negativeTime: negativeTime}

	manager.lk.Lock()
	manager.dict[filepath] = file
	manager.lk.Unlock()
}

// SetDefault set the not exist filepath, negativeTimeout use default
func (manager *NotExistManager) SetDefault(filepath string) {
	manager.Set(filepath, manager.NegativeTimeout)
}

// Del the file path in map
func (manager *NotExistManager) Del(filepath string) {

	manager.lk.Lock()
	delete(manager.dict, filepath)
	manager.lk.Unlock()

}

// IsNotExist if filepath is not exist in cache(so you may test file is exist), return true
// else return false
func (manager *NotExistManager) IsNotExist(filepath string) bool {

	result := false

	var file *notExistFile

	manager.lk.RLock()

	file = manager.dict[filepath]

	if file == nil || !file.isNegative() {
		result = true
	}

	manager.lk.RUnlock()

	if result == true && file != nil {
		manager.Del(filepath)
	}

	return result
}
