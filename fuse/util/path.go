package util

import (
	"sync"
)

// FusePathManager the map to save nodeid and the path of nodeid
type FusePathManager struct {
	pathDict map[uint64]string

	lk sync.RWMutex
}

// Init the function to initialize the map, and set "/" path
func (fp *FusePathManager) Init() {
	fp.pathDict = make(map[uint64]string)

	fp.pathDict[1] = "/"
}

// Set set nodeid and path to map, key: nodeid, val: path
func (fp *FusePathManager) Set(nodeid uint64, path string) {
	fp.lk.Lock()
	fp.pathDict[nodeid] = path
	fp.lk.Unlock()
}

// Get get the path by nodeid
func (fp *FusePathManager) Get(nodeid uint64) string {
	fp.lk.RLock()
	path := fp.pathDict[nodeid]
	fp.lk.RUnlock()
	return path
}

// Del delete the nodeid in map
func (fp *FusePathManager) Del(nodeid uint64) {
	fp.lk.Lock()
	delete(fp.pathDict, nodeid)
	fp.lk.Unlock()
}
