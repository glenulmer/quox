package main

type UserInfo_t struct {
	login int
	greet, email string
}

type UIBagVars_t map[string]string

type State_t struct {
	user UserInfo_t
	quote UIBagVars_t
}

type SessionVars_t struct {
	user UserInfo_t
	quote QuoteVars_t
	device string
	deviceConfirmed bool
}

func InitState() State_t {
	return State_t{ quote: make(map[string]string) }
}

func (x State_t)LoggedIn() bool { return x.user.login > 0 }

func CloneUIBagVars(in UIBagVars_t) UIBagVars_t {
	out := make(UIBagVars_t, len(in))
	for k, v := range in { out[k] = v }
	return out
}
