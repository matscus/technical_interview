package users

import (
	"database/sql"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var (
	u, _ = uuid.Parse("43f73baa-8d4e-4938-9106-86a908b446c3")
	user = User{
		ID:        u,
		FirstName: "Bus",
		LastName:  "Bakinsky",
		Phone:     "03",
		Email:     "bus$pbuspark.az",
		Gender:    "bus",
	}
)

func Test_User_New(t *testing.T) {
	test := New()
	if test.ID.String() == "" {
		t.Errorf("New return ID = %s; ", test.ID)
	}
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func Test_User_Get(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	DB = sqlx.NewDb(mockDB, "sqlmock")
	rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "phone", "email", "gender"}).
		AddRow(user.ID, user.FirstName, user.LastName, user.Phone, user.Email, user.Gender)
	mock.ExpectQuery("SELECT").WithArgs(user.ID).WillReturnRows(rows)
	u, err := Get(user.ID)
	assert.NotNil(t, u)
	assert.NoError(t, err)
}

func Test_User_GetAll(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	DB = sqlx.NewDb(mockDB, "sqlmock")
	rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "phone", "email", "gender"}).
		AddRow(user.ID, user.FirstName, user.LastName, user.Phone, user.Email, user.Gender)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	u, err := GetAll()
	assert.NotNil(t, u)
	assert.NoError(t, err)
}

func Test_User_Create(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	DB = sqlx.NewDb(mockDB, "sqlmock")
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	err := user.Create()
	assert.NoError(t, err)
}

func Test_User_Delete(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	DB = sqlx.NewDb(mockDB, "sqlmock")
	mock.ExpectExec("DELETE").WillReturnResult(sqlmock.NewResult(1, 1))
	err := user.Delete()
	assert.NoError(t, err)
}

func Test_User_Update(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	DB = sqlx.NewDb(mockDB, "sqlmock")
	mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	err := user.Update()
	assert.NoError(t, err)
}
