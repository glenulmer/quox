package main

import (
	"flag"
	"net/http"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	. "klpm/lib/output"
	. "klpm/lib/wrapdb"
)

const (
	envPort              = `PM_PORT`
	envLayout            = `PM_LAYOUT`
	envDBName            = `PM_DBNAME`
	envDBUser            = `PM_DBUSER`
	envDBPass            = `PM_DBPASS`
	dbPingTries          = 5
	dbPingPause          = 200 * time.Millisecond
)

const layoutPhone, layoutDesktop = `phone`, `desktop`

func Bootstrap() {
	var portFlag, layoutFlag, dbNameFlag, dbUserFlag, dbPassFlag string
	flag.StringVar(&portFlag, `port`, ``, Str(`Web port number (or $`, envPort, `)`))
	flag.StringVar(&layoutFlag, `layout`, ``, Str(`Layout name (`, layoutPhone, ` or `, layoutDesktop, `, or $`, envLayout, `)`))
	flag.StringVar(&dbNameFlag, `db`, ``, Str(`Database name (or $`, envDBName, `)`))
	flag.StringVar(&dbUserFlag, `user`, ``, Str(`Database user (or $`, envDBUser, `)`))
	flag.StringVar(&dbPassFlag, `pass`, ``, Str(`Database password (or $`, envDBPass, `)`))
	flag.Parse()

	pick := func(raw ...string) string {
		for _, x := range raw {
			x = strings.TrimSpace(x)
			if x != `` { return x }
		}
		return ``
	}

	port := pick(portFlag, os.Getenv(envPort), `4444`)
	if n, e := strconv.Atoi(port); e != nil || n < 1 || n > 65535 {
		panic(Str(`Invalid port configuration: `, port))
	}
	layoutDefault := layoutPhone
	if port == `3333` { layoutDefault = layoutDesktop }
	layout := pick(layoutFlag, os.Getenv(envLayout), layoutDefault)
	switch layout {
	case layoutPhone, layoutDesktop:
	default:
		panic(Str(`Invalid layout configuration: `, layout))
	}

	dbName := pick(dbNameFlag, os.Getenv(envDBName))
	dbUser := pick(dbUserFlag, os.Getenv(envDBUser))
	dbPass := pick(dbPassFlag, os.Getenv(envDBPass))
	auth := TrueAuth
	if os.Getenv(`HOME`) == `/home/glen` {
		auth = FakeAuth
	}

	missing := make([]string, 0, 3)
	if dbName == `` { missing = append(missing, envDBName) }
	if dbUser == `` { missing = append(missing, envDBUser) }
	if dbPass == `` { missing = append(missing, envDBPass) }
	if len(missing) > 0 {
		panic(Str(`Missing DB configuration: `, strings.Join(missing, `, `)))
	}

	App = App_t{
		DB: OpenDB(dbUser, dbPass, dbName),
		port: port,
		layout: layout,
		Auth: auth,
		sessionStore: NewSessionStore(),
	}

	LoadStaticData()
	App.staticVersion = ComputeStaticVersion()
}

func OpenDB(user, pass, name string) *DB_t {
	var lastErr error
	for try := 1; try <= dbPingTries; try++ {
		dbx := Mysql(user, pass, name).
			Option(`charset`, `utf8`).
			Option(`collation`, `utf8mb4_german2_ci`).
			Open()
		if e := dbx.Ping(); e == nil {
			return dbx
		} else {
			lastErr = e
			_ = dbx.Close()
		}
		if try < dbPingTries { time.Sleep(dbPingPause) }
	}
	panic(Error(`database ping failed: `, lastErr))
}

func ComputeStaticVersion() string {
	root := `./static`
	var count, newest, sumSize int64
	e := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil { return err }
		if d.IsDir() { return nil }
		info, e := d.Info()
		if e != nil { return e }
		count++
		sumSize += info.Size()
		mod := info.ModTime().UnixNano()
		if mod > newest { newest = mod }
		return nil
	})
	if e != nil { panic(Error(`could not compute static version from `, root, `: `, e)) }
	if count == 0 { panic(Error(`could not compute static version from `, root, `: no files`)) }
	return Str(count, `-`, sumSize, `-`, newest)
}

func FakeAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := GetState(r)
		if !state.LoggedIn() {
			state.user = UserInfo_t{ 1, `Glen`, `glen.ulmer@gmail.com` }
			SetState(r, state)
		}
		next.ServeHTTP(w, r)
	}
}

func TrueAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !GetState(r).LoggedIn() {
			http.Redirect(w, r, `/signin`, http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}
