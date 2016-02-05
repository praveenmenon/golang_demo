package account

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	_ "github.com/lib/pq"
	"github.com/praveenmenon/golang_demo/api/v1/controllers"
	"github.com/praveenmenon/golang_demo/api/v1/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

type registrationController struct{}

var Registration registrationController

func (r registrationController) Create(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	flag := 1
	var u models.Register

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &u)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", "password=password host=localhost dbname=golang_demo sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	users, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL,firstname varchar(100),lastname varchar(100),email varchar(320),password varchar(100),password_confirmation varchar(100), createdat timestamptz, updatedat timestamptz, PRIMARY KEY(id))")
	if err != nil || users == nil {
		log.Fatal(err)
	}

	res, err := db.Query("SELECT email FROM users ")
	if err != nil {
		log.Fatal(err)
	}

	fetch_id, err := db.Query("SELECT coalesce(max(id), 0) FROM users")
	if err != nil {
		log.Fatal(err)
	}

	if flag == 1 {
		email := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
		exp, err := regexp.Compile(email)
		if err != nil {
			os.Exit(1)
		}
		if u.Firstname == "" || u.Lastname == "" || u.Email == "" || !exp.MatchString(u.Email) || u.Password == "" || u.Password_confirmation == "" {

			result, err := govalidator.ValidateStruct(u)
			if err != nil {
				println("error: " + err.Error())
			}
			fmt.Println(result)
			flag = 0
			b, err := json.Marshal(models.ErrorMessage{
				Success: "false",
				Error:   err.Error(),
			})
			if err != nil {
				log.Fatal(err)
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(b)
			goto create_user_end
		}
	}
	if flag == 1 {
		for res.Next() {
			var email string
			err = res.Scan(&email)
			if err != nil {
				log.Fatal(err)
			}

			if email == u.Email {
				b, err := json.Marshal(models.ErrorMessage{
					Success: "false",
					Error:   "Email id already exist",
				})
				if err != nil {
					log.Fatal(err)
				}
				rw.Header().Set("Content-Type", "application/json")
				rw.Write(b)
				fmt.Println("Email id already exist")
				flag = 0
				goto create_user_end
			}
		}
		defer res.Close()
		if u.Password != u.Password_confirmation {
			b, err := json.Marshal(models.ErrorMessage{
				Success: "false",
				Error:   "Password and Password_confirmation do not match",
			})
			if err != nil {
				log.Fatal(err)
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(b)
			goto create_user_end
		}
		for fetch_id.Next() {
			var id int
			err = fetch_id.Scan(&id)
			if err != nil {
				log.Fatal(err)
			}
			id = id + 1

			var sStmt string = "insert into users (id,firstname, lastname, email, password, password_confirmation) values ($1,$2,$3,$4,$5,$6)"
			stmt, err := db.Prepare(sStmt)
			if err != nil {
				log.Fatal(err)
			}

			key := []byte("traveling is fun")
			password := []byte(u.Password)
			confirm_password := []byte(u.Password_confirmation)
			encrypt_password := controllers.Encrypt(key, password)
			encrypt_password_confirmation := controllers.Encrypt(key, confirm_password)

			user_res, err := stmt.Exec(id, u.Firstname, u.Lastname, u.Email, encrypt_password, encrypt_password_confirmation)
			if err != nil || user_res == nil {
				log.Fatal(err)
			}

			fmt.Printf("StartTime: %v\n", time.Now())
			fmt.Println("User created Successfully!")

			user := models.Register{id, u.Firstname, u.Lastname, u.Email, u.Password, u.Password_confirmation}

			b, err := json.Marshal(models.SignIn{
				Success: "true",
				Message: "User created Successfully!",
				User:    user,
			})

			if err != nil || res == nil {
				log.Fatal(err)
			}

			rw.Header().Set("Content-Type", "application/json")
			rw.Write(b)
			stmt.Close()
		}
		defer fetch_id.Close()
	}
create_user_end:
	db.Close()
}
