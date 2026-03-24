package main

import "bytes"
import "encoding/gob"
import "strconv"

type UserInfo_t struct {
	login int
	greet, email string
}

type QuoteVars_t map[string]string

type State_t struct {
	user UserInfo_t
	quote QuoteVars_t
}

func InitState() State_t {
	return State_t{ quote: make(map[string]string) }
}

func (x State_t)LoggedIn() bool { return x.user.login > 0 }

func (x State_t)GobEncode() ([]byte, error) {
	qcount := len(x.quote)
	var b bytes.Buffer
	wire := make([]string, 0, 4+qcount*2)
	wire = append(wire,
		strconv.Itoa(x.user.login),
		x.user.greet,
		x.user.email,
		strconv.Itoa(qcount),
	)
	for k, v := range x.quote { wire = append(wire, k, v) }
	e := gob.NewEncoder(&b).Encode(wire)
	if e != nil { return nil, e }
	return b.Bytes(), nil
}

func (x *State_t)GobDecode(in []byte) error {
	if len(in) == 0 { *x = InitState(); return nil }
	wire := []string{}
	e := gob.NewDecoder(bytes.NewReader(in)).Decode(&wire)
	if e != nil { return e }
	if len(wire) == 0 { *x = InitState(); return nil }
	*x = InitState()
	login, _ := strconv.Atoi(wire[0])
	x.user.login = login
	if len(wire) > 1 { x.user.greet = wire[1] }
	if len(wire) > 2 { x.user.email = wire[2] }
	qcount := 0
	if len(wire) > 3 { qcount, _ = strconv.Atoi(wire[3]) }
	at := 4
	for i := 0; i < qcount; i++ {
		if at+1 >= len(wire) { break }
		x.quote[wire[at]] = wire[at+1]
		at += 2
	}
	return nil
}
