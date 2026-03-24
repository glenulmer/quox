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
	var b bytes.Buffer
	e := gob.NewEncoder(&b).Encode([]string{
		strconv.Itoa(x.user.login),
		x.user.greet,
		x.user.email,
	})
	if e != nil { return nil, e }
	return b.Bytes(), nil
}

func (x *State_t)GobDecode(in []byte) error {
	if len(in) == 0 {
		*x = InitState()
		return nil
	}
	wire := []string{}
	e := gob.NewDecoder(bytes.NewReader(in)).Decode(&wire)
	if e != nil { return e }
	if len(wire) == 0 {
		*x = InitState()
		return nil
	}
	login, _ := strconv.Atoi(wire[0])
	x.user.login = login
	if len(wire) > 1 { x.user.greet = wire[1] }
	if len(wire) > 2 { x.user.email = wire[2] }
	return nil
}
