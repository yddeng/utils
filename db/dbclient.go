package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func pgsqlOpen(host string, port int, dbname string, user string, password string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", host, port, dbname, user, password)
	return sql.Open("postgres", connStr)
}

func mysqlOpen(host string, port int, dbname string, user string, password string) (*sql.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbname)
	return sql.Open("mysql", connStr)
}

func sqlOpen(sqlType string, host string, port int, dbname string, user string, password string) (*sql.DB, error) {
	if sqlType == "mysql" {
		return mysqlOpen(host, port, dbname, user, password)
	} else {
		return pgsqlOpen(host, port, dbname, user, password)
	}
}

type Client struct {
	db      *sql.DB
	sqlType string
}

func NewClient(sqlType string, host string, port int, dbname string, user string, password string) (c *Client, err error) {
	c = new(Client)
	c.sqlType = sqlType
	c.db, err = sqlOpen(sqlType, host, port, dbname, user, password)
	if err != nil {
		return
	}

	err = c.db.Ping()
	if err != nil {
		return
	}
	return
}

type Error struct {
	Code ErrCode
	Err  error
}

func (e *Error) IsOK() bool {
	return e.Code == ERR_OK
}

type ErrCode int32

const (
	ERR_OK = ErrCode(iota)
	ERR_BUSY
	ERR_RECORD_EXIST    // key已经存在
	ERR_RECORD_NOTEXIST // key不存在
	ERR_TIMEOUT
	ERR_SQLERROR
	ERR_MISSING_FIELDS //缺少字段
	ERR_MISSING_TABLE  //没有指定表
	ERR_MISSING_KEY    //没有指定key
	ERR_INVAILD_TABLE  //非法表
	ERR_INVAILD_FIELD  //非法字段
	ERR_CAS_NOT_EQUAL
)

var codeStr = map[ErrCode]string{
	ERR_OK:              "ERR_OK",
	ERR_BUSY:            "ERR_BUSY",
	ERR_RECORD_EXIST:    "ERR_RECORD_EXIST",
	ERR_RECORD_NOTEXIST: "ERR_RECORD_NOTEXIST",
	ERR_TIMEOUT:         "ERR_TIMEOUT",
	ERR_SQLERROR:        "ERR_SQLERROR",
	ERR_MISSING_FIELDS:  "ERR_MISSING_FIELDS",
	ERR_MISSING_TABLE:   "ERR_MISSING_TABLE",
	ERR_MISSING_KEY:     "ERR_MISSING_KEY",
	ERR_INVAILD_TABLE:   "ERR_INVAILD_TABLE",
	ERR_INVAILD_FIELD:   "ERR_INVAILD_FIELD",
	ERR_CAS_NOT_EQUAL:   "ERR_CAS_NOT_EQUAL",
}

func (e ErrCode) String() string {
	return codeStr[e]
}
