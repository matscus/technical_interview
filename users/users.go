package users

import (
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	db     *sqlx.DB
	schema = `
	CREATE TABLE IF NOT EXISTS users (
		id        uuid PRIMARY KEY,
		firstname varchar NOT NULL,
		lastname  varchar NOT NULL,
		phone     varchar NOT NULL,
		email      varchar NOT NULL,
		gender    varchar(6) NOT NULL
	);
	`
)

func InitDB(connStr string) {
	var err error
	db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Printf("Database connection error %s\n", err)
	}
	log.Printf("Database connection")
	go func() {
		for {
			err := db.Ping()
			if err != nil {
				log.Printf("database ping error %s\n", err)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	initScheme()
}

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	FirstName string    `json:"firstname" db:"firstname"`
	LastName  string    `json:"lastname" db:"lastname"`
	Phone     string    `json:"phone" db:"phone"`
	Email     string    `json:"email" db:"email"`
	Gender    string    `json:"gender" db:"gender"`
}

//New  - create and return new fake user
func New() User {
	return User{
		ID:        uuid.New(),
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Phone:     gofakeit.Phone(),
		Email:     gofakeit.Email(),
		Gender:    gofakeit.Gender(),
	}
}

func GetAllUsers() (users []User, err error) {
	return users, db.Select(&users, "SELECT * FROM users")
}

func GetUser(uuid uuid.UUID) (user User, err error) {
	return user, db.Get(&user, "SELECT * FROM users where id=$1", uuid)
}

func (u *User) Create() error {

	return nil
}

func (u *User) ChangeUser() error {

	return nil
}

func initScheme() {
	user := []User{}
	err := db.Select(&user, "select * from users limit 1")
	if err != nil {
		log.Println("schema users is not exist, start created...")
		db.MustExec(schema)
		log.Println("creater complited")
		insertTestValues()
	} else {
		log.Println("schema users is exist, start drop...")
		db.MustExec("drop table users")
		log.Println("drop complited")
		log.Println("start created schema...")
		db.MustExec(schema)
		log.Println("created complited")
		insertTestValues()
	}
}

func insertTestValues() {
	log.Println("start insert test values")
	tx := db.MustBegin()
	for i := 0; i < 100; i++ {
		u := New()
		_, err := tx.NamedExec("INSERT INTO users (id,firstname,lastname,phone,email,gender) VALUES (:id,:firstname,:lastname,:phone,:email,:gender)", &u)
		if err != nil {
			log.Println(err)
		}
	}
	tx.Commit()
	log.Println("insert complited")
}
