package users

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/praveenmenon/golang_demo/api/v1/models"
	"log"
	"net/http"
)

type listController struct{}

var List listController

func (l listController) List_users(rw http.ResponseWriter, req *http.Request) {

	var list models.UserList
	flag := 1
	db, err := sql.Open("postgres", "password=password host=localhost dbname=golang_demo sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	get_users, err := db.Query("SELECT id,firstname,lastname,email FROM users")
	if err != nil || get_users == nil {
		log.Fatal(err)
	}
	var user_id int
	var firstname string
	var lastname string
	var email string
	var no_of_users int

	for get_users.Next() {
		err := get_users.Scan(&user_id, &firstname, &lastname, &email)
		if err != nil {
			log.Fatal(err)
		}
		profile := models.UserDetails{user_id, firstname, lastname, email}
		list.User_Details = append(list.User_Details, profile)
		no_of_users++
		flag = 0
	}
	defer get_users.Close()
	if flag == 0 {
		b, err := json.Marshal(models.UserList{
			Success:      "true",
			No_Of_Users:  no_of_users,
			User_Details: list.User_Details,
		})
		if err != nil {
			log.Fatal(err)
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write(b)
		goto users_end
	}
	if flag == 1 {
		b, err := json.Marshal(models.ErrorMessage{
			Success: "false",
			Error:   "No Users",
		})

		if err != nil {
			log.Fatal(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(b)
	}
users_end:
	db.Close()
}
