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
	"errors"
)

var (
	ErrLocked = errors.New("already locked")
	ErrUnlock = errors.New("already unlocked")
)

type Flocker interface {
	// Lock is a none-blocking call to try and take the file lock.
	// If we are already locked, this function short-circuits and returns immediately
	//
	// 非阻塞锁定文件，如果Locked()为true，则立即返回nil
	Lock() error

	// Unlock is a function to unlock the file.
	// If we are already unlocked, this function short-circuits and returns immediately
	//
	// 非阻塞解锁文件，如果Locked()为false，则立即返回ErrUnlock
	Unlock() error

	// Locked is a function to return the current lock state (locked: true, unlocked: false).
	//
	//返回当前的锁定状态
	Locked() bool

	// Path is a function to return the path as provided in NewFlock().
	//
	//返回当前的文件
	Path() string
}
