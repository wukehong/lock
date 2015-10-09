// +build !lockfileex

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
	"io"
	"path/filepath"
	"syscall"
)

const (
	_FILE_ATTRIBUTE_TEMPORARY  = 0x100
	_FILE_FLAG_DELETE_ON_CLOSE = 0x04000000
)

func init() {
	// sane default
	lockFn = lockCreateFile
}

type handleUnlocker struct {
	h syscall.Handle
}

func (hu *handleUnlocker) Close() error {
	return syscall.Close(hu.h)
}

func lockCreateFile(name string) (io.Closer, error) {
	absName, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}
	pName, err := syscall.UTF16PtrFromString(absName)
	if err != nil {
		return nil, err
	}
	// http://msdn.microsoft.com/en-us/library/windows/desktop/aa363858%28v=vs.85%29.aspx
	h, err := syscall.CreateFile(pName,
		syscall.GENERIC_WRITE, // open for write
		0,   // no sharing
		nil, // don't let children inherit
		syscall.CREATE_ALWAYS, // create if not exists, truncate if does
		syscall.FILE_ATTRIBUTE_NORMAL|_FILE_ATTRIBUTE_TEMPORARY|_FILE_FLAG_DELETE_ON_CLOSE,
		0)
	if err != nil {
		return nil, err
	}
	return &handleUnlocker{h}, nil
}
