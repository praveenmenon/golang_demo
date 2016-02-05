package models

// Registration struct [account/sign_up]
type Register struct {
	Id                    int    `valid:"numeric"`
	Firstname             string `valid:"alphanum,required"`
	Lastname              string `valid:"alphanum,required"`
	Email                 string `valid:"email,required"`
	Password              string `valid:"alphanum,required"`
	Password_confirmation string `valid:"alphanum,required"`
}

type UserDetails struct {
	Id        int
	Firstname string
	Lastname  string
	Email     string
}

type Notification struct {
	SenderId   int `valid:"numeric,required"`
	RecieverId int `valid:"numeric,required"`
}

type SignIn struct {
	Success string
	Message string
	User    Register
}

// Message struct [controllers/account]
// Common for sign_up, session and password
type Message struct {
	Success string
	Message string
	User    Register
}

type ErrorMessage struct {
	Success string
	Error   string
}

// User profile Struct

type ProfileErrorMessage struct {
	Success string
	Error   string
}

type UserListMessage struct {
	Success  string
	Message  string
	User_ids []int
}

type UserList struct {
	Success      string
	No_Of_Users  int
	User_Details []UserDetails
}

type ShowErrorMessage struct {
	Success string
	Error   string
}
