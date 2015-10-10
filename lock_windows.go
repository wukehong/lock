// +build windows

/*
Copyright 2013 The Go Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lock

import (
	"path/filepath"
	"sync"
	"syscall"
)

const (
	_FILE_ATTRIBUTE_TEMPORARY  = 0x100
	_FILE_FLAG_DELETE_ON_CLOSE = 0x04000000
)

type flock struct {
	path      string
	absPath   string
	utf16Path *uint16
	mu        sync.RWMutex
	fh        syscall.Handle
	locked    bool
}

func NewFlock(path string) FLocker {
	f := &flock{path: path}
	f.absPath, _ = filepath.Abs(path)
	f.utf16Path, _ = syscall.UTF16PtrFromString(f.absPath)
	return f
}

func (f *flock) Lock() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.locked {
		return nil
	}

	// http://msdn.microsoft.com/en-us/library/windows/desktop/aa363858%28v=vs.85%29.aspx
	fh, err := syscall.CreateFile(f.utf16Path,
		syscall.GENERIC_WRITE, // open for write
		0,   // no sharing
		nil, // don't let children inherit
		syscall.CREATE_ALWAYS, // create if not exists, truncate if does
		syscall.FILE_ATTRIBUTE_NORMAL|_FILE_ATTRIBUTE_TEMPORARY|_FILE_FLAG_DELETE_ON_CLOSE,
		0)
	f.locked = err == nil
	if f.locked {
		f.fh = fh
	}
	return err
}

func (f *flock) TryLock() error {
	return f.Lock()
}

func (f *flock) Unlock() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.locked {
		return ErrUnlock
	}

	err := syscall.Close(f.fh)
	if err == nil {
		f.locked = false
		f.fh = 0
	}
	return err
}

func (f *flock) Locked() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.locked
}

func (f *flock) Path() string {
	return f.path
}
