// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package orm

import (
	"database/sql"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/alfatih/irhabi/irhabi/common/log"
	"github.com/Sirupsen/logrus"
)

// NewLog set io.Writer to create a Logger.
func NewLog(out io.Writer) *logrus.Logger {
	l := log.New()
	l.Out = out
	return l
}

func logQuery(alias *alias, operaton, query string, t time.Time, err error, args ...interface{}) {
	go func() {
		sub := time.Now().Sub(t) / 1e5
		elsp := float64(int(sub)) / 10.0

		query = fmt.Sprintf(strings.Replace(query, "?", "`%v`", -1), args...)
		l := DebugLog.WithFields(logrus.Fields{
			"latecy": fmt.Sprintf("%1.1fms", elsp),
			"query":  query,
		})

		if err != nil {
			l.WithField("prefix", fmt.Sprintf("%s/%s", operaton, "FAIL")).Error(err)
		} else {
			l.WithField("prefix", fmt.Sprintf("%s/%s", operaton, "OK")).Info("Success")
		}
	}()
}

// statement query logger struct.
// if dev mode, use stmtQueryLog, or use stmtQuerier.
type stmtQueryLog struct {
	alias *alias
	query string
	stmt  stmtQuerier
}

var _ stmtQuerier = new(stmtQueryLog)

func (d *stmtQueryLog) Close() error {
	a := time.Now()
	err := d.stmt.Close()
	logQuery(d.alias, "st.Close", d.query, a, err)
	return err
}

func (d *stmtQueryLog) Exec(args ...interface{}) (sql.Result, error) {
	a := time.Now()
	res, err := d.stmt.Exec(args...)
	logQuery(d.alias, "st.Exec", d.query, a, err, args...)
	return res, err
}

func (d *stmtQueryLog) Query(args ...interface{}) (*sql.Rows, error) {
	a := time.Now()
	res, err := d.stmt.Query(args...)
	logQuery(d.alias, "st.Query", d.query, a, err, args...)
	return res, err
}

func (d *stmtQueryLog) QueryRow(args ...interface{}) *sql.Row {
	a := time.Now()
	res := d.stmt.QueryRow(args...)
	logQuery(d.alias, "st.QueryRow", d.query, a, nil, args...)
	return res
}

func newStmtQueryLog(alias *alias, stmt stmtQuerier, query string) stmtQuerier {
	d := new(stmtQueryLog)
	d.stmt = stmt
	d.alias = alias
	d.query = query
	return d
}

// database query logger struct.
// if dev mode, use dbQueryLog, or use dbQuerier.
type dbQueryLog struct {
	alias *alias
	db    dbQuerier
	tx    txer
	txe   txEnder
}

var _ dbQuerier = new(dbQueryLog)
var _ txer = new(dbQueryLog)
var _ txEnder = new(dbQueryLog)

func (d *dbQueryLog) Prepare(query string) (*sql.Stmt, error) {
	a := time.Now()
	stmt, err := d.db.Prepare(query)
	logQuery(d.alias, "db.Prepare", query, a, err)
	return stmt, err
}

func (d *dbQueryLog) Exec(query string, args ...interface{}) (sql.Result, error) {
	a := time.Now()
	res, err := d.db.Exec(query, args...)
	logQuery(d.alias, "db.Exec", query, a, err, args...)
	return res, err
}

func (d *dbQueryLog) Query(query string, args ...interface{}) (*sql.Rows, error) {
	a := time.Now()
	res, err := d.db.Query(query, args...)
	logQuery(d.alias, "db.Query", query, a, err, args...)
	return res, err
}

func (d *dbQueryLog) QueryRow(query string, args ...interface{}) *sql.Row {
	a := time.Now()
	res := d.db.QueryRow(query, args...)
	logQuery(d.alias, "db.QueryRow", query, a, nil, args...)
	return res
}

func (d *dbQueryLog) Begin() (*sql.Tx, error) {
	a := time.Now()
	tx, err := d.db.(txer).Begin()
	logQuery(d.alias, "db.Begin", "START TRANSACTION", a, err)
	return tx, err
}

func (d *dbQueryLog) Commit() error {
	a := time.Now()
	err := d.db.(txEnder).Commit()
	logQuery(d.alias, "tx.Commit", "COMMIT", a, err)
	return err
}

func (d *dbQueryLog) Rollback() error {
	a := time.Now()
	err := d.db.(txEnder).Rollback()
	logQuery(d.alias, "tx.Rollback", "ROLLBACK", a, err)
	return err
}

func (d *dbQueryLog) SetDB(db dbQuerier) {
	d.db = db
}

func newDbQueryLog(alias *alias, db dbQuerier) dbQuerier {
	d := new(dbQueryLog)
	d.alias = alias
	d.db = db
	return d
}
