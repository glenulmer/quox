package main

type UserInfo_t struct {
	login int
	greet, email string
}

type State_t struct {
	user UserInfo_t
}

func InitState() State_t {
	return State_t{ }
}

func (x State_t)LoggedIn() bool { return x.user.login > 0 }
