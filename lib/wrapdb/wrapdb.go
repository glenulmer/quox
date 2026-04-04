package wrapdb

import (
	. "quo2/lib/output"
	sql "database/sql"
	. "github.com/go-sql-driver/mysql"
)

type DB_t sql.DB
type Tx_t sql.Tx
type Row_t struct { e error; q string; sqlrow *sql.Row }
type Rows_t sql.Rows

type Pack_t interface {
	Error() error
	Message() string
	HasError() bool
	IsRow() bool
	Query() string
}

type tPackRows struct { e error; q string; rows *Rows_t }
type tPackRow struct { e error; q string }

func (p tPackRows)Error() error { return p.e }
func (p tPackRows)HasError() bool { return p.e != nil }
func (p tPackRows)Query() string { return p.q }
func (tPackRows)IsRow() bool { return false }

func (p tPackRow)Error() error { return p.e }
func (p tPackRow)HasError() bool { return p.e != nil }
func (p tPackRow)Query() string { return p.q }
func (tPackRow)IsRow() bool { return true }

func _PackRows(err error, query string, rows *sql.Rows) tPackRows { return tPackRows{ e:err, q:query, rows:wrapRows(rows) } }
func _PackRow(err error, query string ) tPackRow { return tPackRow{ e:err, q: query } }

func wrapRows(rows *sql.Rows) *Rows_t {
	if rows == nil { return nil }
	r := Rows_t(*rows)
	return &r
}

var doLogSQL = false

func ShowQueryOn() { doLogSQL = true }
func ShowQueryOff() { doLogSQL = false }
func logSQL(q string) {
	if !doLogSQL { return }
	Log(q)
}

func (r *Rows_t)Columns() []string {
	cc, _ := r.unwrap().Columns()
	return cc
}

// Close, Err, Next, Scan
func (r *Rows_t)Close() error { return r.unwrap().Close() }
func (r *Rows_t)Err() error { return r.unwrap().Err() }
func (r *Rows_t)Next() bool { return r.unwrap().Next() }
func (r *Rows_t)Scan(dest ...interface{}) error { return r.unwrap().Scan(dest...) }
func (r *Rows_t)NextResultSet() bool { return r.unwrap().NextResultSet() }

func (p tPackRows)Close() error { return p.rows.unwrap().Close() }
func (p tPackRows)Err() error { return p.rows.unwrap().Err() }
func (p tPackRows)Next() bool { return p.rows.unwrap().Next() }
func (p tPackRows)Scan(dest ...interface{}) error { return p.rows.unwrap().Scan(dest...) }
func (p tPackRows)NextResultSet() bool { return p.rows.unwrap().NextResultSet() }

type XQL interface { // for DB and Tx
	Query(q string) tPackRows // *Rows_t
	QueryRow(q string) *Row_t
	Call(proc string, parms ...interface{}) tPackRows // *Rows_t
	CallRow(proc string, parms ...interface{}) *Row_t
}

func callStr(proc string, parms ...interface{}) string {
	var sql Builder
	sql.Add("call ", proc, "(")
	nargs := len(parms) - 1
	for k, arg := range parms {
		switch v := arg.(type) {
		case string:
			sql.Add(DQ(v))
		case bool:
			sql.Add(If(v,`1`,`0`))
		default:
			sql.Add(v)
		}
		if k < nargs { sql.Add(",") }
	}
	sql.Add(")")
	logSQL(sql.String())
	return sql.String()
}

func (z *DB_t)Ping() error { return z.unwrap().Ping() }

func (z *DB_t)Query(q string) tPackRows {
	r, e := z.unwrap().Query(q)
	return _PackRows(e, q, r)
}

func (z *Tx_t)Query(q string) tPackRows {
	r, e := z.unwrap().Query(q)
	return _PackRows(e, q, r)
}

func (z *DB_t)QueryRow(q string) *Row_t {
	var r = z.unwrap().QueryRow(q)
	pr := Row_t{ q: q, sqlrow: r }
	return &pr
}

func (z *Tx_t)QueryRow(q string) *Row_t {
	var r = z.unwrap().QueryRow(q)
	pr := Row_t{ q: q, sqlrow: r }
	return &pr
}

func (z *DB_t)Call(proc string, parms ...interface{}) tPackRows {
	return z.Query(callStr(proc, parms...))
}

func (z *Tx_t)Call(proc string, parms ...interface{}) tPackRows {
	return z.Query(callStr(proc, parms...))
}

func (z *Row_t)Scan(dest ...interface{}) tPackRow {
	z.e = z.sqlrow.Scan(dest...)
	return _PackRow(z.e, z.q)
}

func (z *Row_t)Error() error {
	if z.e != nil { return z.e }
	if z.sqlrow == nil { return nil }
	z.e = z.sqlrow.Err()
	return z.e
}

func (z *Row_t)HasError() bool { return z.Error() != nil }
func (z *Row_t)Query() string { return z.q }
func (z *Row_t)Message() string { return mysqlerr(z.q, z.Error()) }

func (z *DB_t)CallRow(proc string, parms ...interface{}) *Row_t {
	return z.QueryRow(callStr(proc, parms...))
}

func (z *Tx_t)CallRow(proc string, parms ...interface{}) *Row_t {
	return z.QueryRow(callStr(proc, parms...))
}

// Transaction support

func (z *DB_t)Begin() (*Tx_t, error) {
	tx, e := z.unwrap().Begin()
	return (*Tx_t)(tx), e
}

func (z *Tx_t)Commit() error {
	return z.unwrap().Commit()
}

func (z *Tx_t)Rollback() error {
	return z.unwrap().Rollback()
}

func (z *DB_t)WithTx(do func(*Tx_t) error) error {
	tx, e := z.Begin()
	if e != nil { return e }

	e = do(tx)
	if e != nil {
		tx.Rollback()
		return e
	}

	return tx.Commit()
}

////////////////////////

type tMysql struct {
	user, pass, name, serv string
	options Builder
}

func Mysql(user, pass, dbname string) *tMysql {
	return &tMysql{ user: user, pass: pass, name: dbname }
}
const nilMsg = "nil tMysql (wrapdb)"

func (x *tMysql)Server(serv string) *tMysql { x.serv = serv; return x }
func (x *tMysql)Option(name string, value interface{}) *tMysql {
	ch := "&"
	if x.options.Len() == 0 { ch = "?" }
	x.options.Add(ch, name,"=",value)
	return x
}

func (x *tMysql)Open() *DB_t {
	conn := Str(x.user, ":", x.pass, "@", x.serv, "/", x.name, x.options.String())
	db, e := sql.Open("mysql", conn)
	if e != nil { panic(e) }
	db.Ping()
	return (*DB_t)(db)
}

func (z *DB_t)Close() error { return z.unwrap().Close() }

func (rr *Rows_t)ScanSlice() (slice []string) {
	var s [24]string

	r := rr.unwrap()
	cols, e := r.Columns()
	if e != nil {
		Log(e.Error())
		return
	}
	ncols := len(cols)
	if ncols == 0 { return }
	if ncols > len(s) { Log(Error(`More than 24 columns in query`)); return }
	switch ncols {
		case  1: e = r.Scan(&s[0])
		case  2: e = r.Scan(&s[0],&s[1])
		case  3: e = r.Scan(&s[0],&s[1],&s[2])
		case  4: e = r.Scan(&s[0],&s[1],&s[2],&s[3])
		case  5: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4])
		case  6: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5])
		case  7: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5],&s[6])
		case  8: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5],&s[6],&s[7])
		case  9: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5],&s[6],&s[7],&s[8])
		case 10: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5],&s[6],&s[7],&s[8],&s[9])
		case 11: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5],&s[6],&s[7],&s[8],&s[9],&s[10])
		case 12: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5],&s[6],&s[7],&s[8],&s[9],&s[10],&s[11])
		case 13: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5],&s[6],&s[7],&s[8],&s[9],&s[10],&s[11],&s[12])
		case 14: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5],&s[6],&s[7],&s[8],&s[9],&s[10],&s[11],&s[12],&s[13])
		case 15: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5],&s[6],&s[7],&s[8],&s[9],&s[10],&s[11],&s[12],&s[13],&s[14])
		case 16: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5],&s[6],&s[7],&s[8],&s[9],&s[10],&s[11],&s[12],&s[13],&s[14],&s[15])
		case 17: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5],&s[6],&s[7],&s[8],&s[9],&s[10],&s[11],&s[12],&s[13],&s[14],&s[15],&s[16])
		case 18: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5],&s[6],&s[7],&s[8],&s[9],&s[10],&s[11],&s[12],&s[13],&s[14],&s[15],&s[16],&s[17])
		case 19: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5],&s[6],&s[7],&s[8],&s[9],&s[10],&s[11],&s[12],&s[13],&s[14],&s[15],&s[16],&s[17],&s[18])
		case 20: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5],&s[6],&s[7],&s[8],&s[9],&s[10],&s[11],&s[12],&s[13],&s[14],&s[15],&s[16],&s[17],&s[18],&s[19])
		case 21: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5],&s[6],&s[7],&s[8],&s[9],&s[10],&s[11],&s[12],&s[13],&s[14],&s[15],&s[16],&s[17],&s[18],&s[19],&s[20])
		case 22: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5],&s[6],&s[7],&s[8],&s[9],&s[10],&s[11],&s[12],&s[13],&s[14],&s[15],&s[16],&s[17],&s[18],&s[19],&s[20],&s[21])
		case 23: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5],&s[6],&s[7],&s[8],&s[9],&s[10],&s[11],&s[12],&s[13],&s[14],&s[15],&s[16],&s[17],&s[18],&s[19],&s[20],&s[21],&s[22])
		case 24: e = r.Scan(&s[0],&s[1],&s[2],&s[3],&s[4],&s[5],&s[6],&s[7],&s[8],&s[9],&s[10],&s[11],&s[12],&s[13],&s[14],&s[15],&s[16],&s[17],&s[18],&s[19],&s[20],&s[21],&s[22],&s[23])
	}
	if e != nil {
		Log(e.Error())
		return
	}

	for k := 0; k < ncols; k++ {
		slice = append(slice, s[k])
	}

	return slice
}

/// Helpers - error printing, unwrapping.

func mysqlerr(q string, err error) string {
	if err != nil {
		if mysqlErr, ok := err.(*MySQLError); ok {
			return mysqlErr.Message + `::` + q
		}
	}
	return ""
}

func (p tPackRows)Message() string { return mysqlerr(p.q, p.e) }
func (p tPackRow)Message() string { return mysqlerr(p.q, p.e) }

func (z *DB_t)unwrap() *sql.DB { return (*sql.DB)(z); }
func (z *Tx_t)unwrap() *sql.Tx { return (*sql.Tx)(z); }
func (z *Rows_t)unwrap() *sql.Rows { return (*sql.Rows)(z); }
func (z *Row_t)unwrap() *sql.Row { return z.sqlrow; }
