// Copyright (C) OpenAI. Derived from the MIT licensed go-sqlite3 driver.
//
// Use of this source code is governed by an MIT-style license. The original
// license is in the LICENSE file.

//go:build cgo
// +build cgo

package sqlite3

/*
#cgo CFLAGS: -DSQLITE_ENABLE_SESSION -DSQLITE_ENABLE_PREUPDATE_HOOK
#ifndef USE_LIBSQLITE3
#include "sqlite3-binding.h"
#else
#include <sqlite3.h>
#endif
#include <stdlib.h>
#include <stdint.h>
int sqlite3GoSessionApply(sqlite3*, int, void*, uintptr_t);
int sqlite3GoSetDefensive(sqlite3*, int);
int sqlite3GoSetDQS(sqlite3*, int);
int sqlite3GoEnableLoadExtension(sqlite3*, int);
*/
import "C"

import (
	"errors"
	"runtime/cgo"
	"unsafe"
)

const (
	SQLITE_CHANGESET_OMIT    = 0
	SQLITE_CHANGESET_REPLACE = 1
	SQLITE_CHANGESET_ABORT   = 2
)

// SQLiteSession records table changes made through a SQLiteConn.
type SQLiteSession struct {
	conn    *SQLiteConn
	session *C.sqlite3_session
}

// SQLiteChangesetCallbacks are invoked synchronously while applying a changeset.
type SQLiteChangesetCallbacks struct {
	Filter   func(string) bool
	Conflict func(int) int
	Panic    any
}

// Session starts recording changes in the named database.
func (c *SQLiteConn) Session(database string) (*SQLiteSession, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.db == nil {
		return nil, errors.New("sqlite database is closed")
	}
	db := C.CString(database)
	defer C.free(unsafe.Pointer(db))
	var session *C.sqlite3_session
	if code := C.sqlite3session_create(c.db, db, &session); code != C.SQLITE_OK {
		return nil, c.lastError()
	}
	return &SQLiteSession{conn: c, session: session}, nil
}

// Attach adds a table to the set of tables recorded by the session.
func (s *SQLiteSession) Attach(table string) error {
	if s == nil || s.session == nil || s.conn == nil {
		return errors.New("sqlite session is closed")
	}
	s.conn.mu.Lock()
	defer s.conn.mu.Unlock()
	var name *C.char
	if table != "" {
		name = C.CString(table)
		defer C.free(unsafe.Pointer(name))
	}
	if code := C.sqlite3session_attach(s.session, name); code != C.SQLITE_OK {
		return s.conn.lastError()
	}
	return nil
}

func (s *SQLiteSession) changeset(patchset bool) ([]byte, error) {
	if s == nil || s.session == nil || s.conn == nil {
		return nil, errors.New("sqlite session is closed")
	}
	s.conn.mu.Lock()
	defer s.conn.mu.Unlock()
	var size C.int
	var data unsafe.Pointer
	var code C.int
	if patchset {
		code = C.sqlite3session_patchset(s.session, &size, &data)
	} else {
		code = C.sqlite3session_changeset(s.session, &size, &data)
	}
	if code != C.SQLITE_OK {
		return nil, s.conn.lastError()
	}
	defer C.sqlite3_free(data)
	if size == 0 {
		return []byte{}, nil
	}
	return C.GoBytes(data, size), nil
}

// Changeset returns the recorded changes in changeset format.
func (s *SQLiteSession) Changeset() ([]byte, error) { return s.changeset(false) }

// Patchset returns the recorded changes in patchset format.
func (s *SQLiteSession) Patchset() ([]byte, error) { return s.changeset(true) }

// Close releases the native session before its database connection closes.
func (s *SQLiteSession) Close() error {
	if s == nil || s.session == nil {
		return nil
	}
	if s.conn != nil {
		s.conn.mu.Lock()
		C.sqlite3session_delete(s.session)
		s.conn.mu.Unlock()
	}
	s.session = nil
	return nil
}

// ApplyChangeset applies a changeset to the main database.
func (c *SQLiteConn) ApplyChangeset(changeset []byte, callbacks *SQLiteChangesetCallbacks) error {
	if c == nil || !c.dbConnOpen() {
		return errors.New("sqlite database is closed")
	}
	if len(changeset) > int(^uint32(0)>>1) {
		return errors.New("changeset exceeds SQLite's maximum input size")
	}
	var data unsafe.Pointer
	if len(changeset) > 0 {
		data = C.CBytes(changeset)
		defer C.free(data)
	}
	var handle cgo.Handle
	var context C.uintptr_t
	if callbacks != nil && (callbacks.Filter != nil || callbacks.Conflict != nil) {
		handle = cgo.NewHandle(callbacks)
		defer handle.Delete()
		context = C.uintptr_t(handle)
	}
	c.mu.Lock()
	code := C.sqlite3GoSessionApply(c.db, C.int(len(changeset)), data, context)
	c.mu.Unlock()
	if callbacks != nil && callbacks.Panic != nil {
		return errors.New("SQLite changeset callback failed")
	}
	if code != C.SQLITE_OK {
		return c.lastError()
	}
	return nil
}

// SetDefensive enables SQLite's defensive database configuration.
func (c *SQLiteConn) SetDefensive(enabled bool) error {
	if c == nil || !c.dbConnOpen() {
		return errors.New("sqlite database is closed")
	}
	value := C.int(0)
	if enabled {
		value = 1
	}
	c.mu.Lock()
	code := C.sqlite3GoSetDefensive(c.db, value)
	c.mu.Unlock()
	if code != C.SQLITE_OK {
		return c.lastError()
	}
	return nil
}

// SetDoubleQuotedStringLiterals controls DQS handling for DDL and DML.
func (c *SQLiteConn) SetDoubleQuotedStringLiterals(enabled bool) error {
	if c == nil || !c.dbConnOpen() {
		return errors.New("sqlite database is closed")
	}
	value := C.int(0)
	if enabled {
		value = 1
	}
	c.mu.Lock()
	code := C.sqlite3GoSetDQS(c.db, value)
	c.mu.Unlock()
	if code != C.SQLITE_OK {
		return c.lastError()
	}
	return nil
}

// EnableLoadExtension toggles SQLite extension loading on the connection.
func (c *SQLiteConn) EnableLoadExtension(enabled bool) error {
	if c == nil || !c.dbConnOpen() {
		return errors.New("sqlite database is closed")
	}
	value := C.int(0)
	if enabled {
		value = 1
	}
	c.mu.Lock()
	code := C.sqlite3GoEnableLoadExtension(c.db, value)
	c.mu.Unlock()
	if code != C.SQLITE_OK {
		return c.lastError()
	}
	return nil
}

//export goSQLite3SessionFilter
func goSQLite3SessionFilter(handle C.uintptr_t, table *C.char) (result C.int) {
	defer func() {
		if failure := recover(); failure != nil {
			if callback, ok := cgo.Handle(handle).Value().(*SQLiteChangesetCallbacks); ok {
				callback.Panic = failure
			}
			result = 0
		}
	}()
	callback := cgo.Handle(handle).Value().(*SQLiteChangesetCallbacks)
	if callback.Filter == nil || callback.Filter(C.GoString(table)) {
		return 1
	}
	return 0
}

//export goSQLite3SessionConflict
func goSQLite3SessionConflict(handle C.uintptr_t, kind C.int) (result C.int) {
	defer func() {
		if failure := recover(); failure != nil {
			if callback, ok := cgo.Handle(handle).Value().(*SQLiteChangesetCallbacks); ok {
				callback.Panic = failure
			}
			result = C.SQLITE_CHANGESET_ABORT
		}
	}()
	callback := cgo.Handle(handle).Value().(*SQLiteChangesetCallbacks)
	if callback.Conflict == nil {
		return C.SQLITE_CHANGESET_ABORT
	}
	return C.int(callback.Conflict(int(kind)))
}
