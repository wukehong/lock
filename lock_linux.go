// +build darwin dragonfly freebsd linux netbsd openbsd

// Copyright 2015 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

// Package flock implements a thread-safe sync.Locker interface for file locking.
// It also includes a non-blocking TryLock() function to allow locking
// without blocking execution.
//
// Package flock is released under the BSD 3-Clause License. See the LICENSE file
// for more details.

//fork from https://github.com/theckman/go-flock/blob/master/flock.go

package lock

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

type flock struct {
	path    string
	absPath string
	mu      sync.RWMutex
	fh      *os.File
	locked  bool
}

func NewFlock(path string) Flocker {
	f := &flock{path: path}
	f.absPath, _ = filepath.Abs(path)
	return f
}

func (f *flock) Path() string {
	return f.path
}

func (f *flock) Locked() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.locked
}

func (f *flock) String() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.locked {
		return fmt.Sprintf("'%s' locked", f.path)
	} else {
		return fmt.Sprintf("'%s' unlock", f.path)
	}
}

func (f *flock) NBLock() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.locked {
		return nil
	}

	if f.fh == nil {
		if err := f.setFh(); err != nil {
			return err
		}
	}

	err := syscall.Flock(int(f.fh.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)

	f.locked = err == nil
	return err
}

func (f *flock) Unlock() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// if we aren't locked or if the lockfile instance is nil
	// just return a nil error because we are unlocked
	if !f.locked {
		return ErrUnlock
	}

	// mark the file as unlocked
	err := syscall.Flock(int(f.fh.Fd()), syscall.LOCK_UN)

	if err == nil {
		f.fh.Close()
		f.locked = false
		f.fh = nil
	}
	return err
}

func (f *flock) setFh() error {
	// open a new os.File instance
	// create it if it doesn't exist, truncate it if it does exist, open the file read-write
	fh, err := os.OpenFile(f.absPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.FileMode(0600))

	if err != nil {
		return err
	}

	// set the filehandle on the struct
	f.fh = fh
	return nil
}
