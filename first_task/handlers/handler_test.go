package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/matscus/technical_interview/first_task/users"
)

var (
	u, _ = uuid.Parse("43f73baa-8d4e-4938-9106-86a908b446c3")
	user = users.User{
		ID:        u,
		FirstName: "Bus",
		LastName:  "Bakinsky",
		Phone:     "03",
		Email:     "bus@buspark.az",
		Gender:    "bus",
	}
)

func Test_Handlers_GetAllUsers(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	users.DB = sqlx.NewDb(mockDB, "sqlmock")
	rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "phone", "email", "gender"}).
		AddRow(user.ID, user.FirstName, user.LastName, user.Phone, user.Email, user.Gender)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	req, err := http.NewRequest("GET", "/api/v1/users/getusers", nil)
	if err != nil {
		t.Fatal(err)
	}
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllUsers)
	handler.ServeHTTP(responseRecorder, req)
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `[{"id":"43f73baa-8d4e-4938-9106-86a908b446c3","firstname":"Bus","lastname":"Bakinsky","phone":"03","email":"bus@buspark.az","gender":"bus"}]`
	users := make([]users.User, 0, 1)
	if err := json.NewDecoder(responseRecorder.Body).Decode(&users); err != nil {
		log.Fatalln(err)
	}
	if users[0].ID != u {
		t.Errorf("handler returned unexpected body: got %v want %v",
			responseRecorder.Body.String(), expected)
	}
}

func Test_Handlers_GetUser(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	users.DB = sqlx.NewDb(mockDB, "sqlmock")
	rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "phone", "email", "gender"}).
		AddRow(user.ID, user.FirstName, user.LastName, user.Phone, user.Email, user.Gender)
	mock.ExpectQuery("SELECT").WithArgs(user.ID).WillReturnRows(rows)
	req, err := http.NewRequest("GET", "/api/v1/users/getuser", nil)
	if err != nil {
		t.Fatal(err)
	}
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUser)
	vars := map[string]string{
		"id": user.ID.String(),
	}
	req = mux.SetURLVars(req, vars)
	handler.ServeHTTP(responseRecorder, req)
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"id":"43f73baa-8d4e-4938-9106-86a908b446c3","firstname":"Bus","lastname":"Bakinsky","phone":"03","email":"bus@buspark.az","gender":"bus"}`
	usr := users.User{}
	if err := json.NewDecoder(responseRecorder.Body).Decode(&usr); err != nil {
		log.Fatalln(err)
	}
	if usr.ID != u {
		t.Errorf("handler returned unexpected body: got %v want %v",
			responseRecorder.Body.String(), expected)
	}
}

func Test_Handlers_DeleteUser(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	users.DB = sqlx.NewDb(mockDB, "sqlmock")
	mock.ExpectExec("DELETE").WillReturnResult(sqlmock.NewResult(1, 1))
	userJSON, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/api/v1/users/deleteusers", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteUser)
	handler.ServeHTTP(responseRecorder, req)
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func Test_Handlers_CreateRandomUser(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	users.DB = sqlx.NewDb(mockDB, "sqlmock")
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	req, err := http.NewRequest("POST", "/api/v1/users/createrandomuser", bytes.NewBuffer([]byte("{}")))
	if err != nil {
		t.Fatal(err)
	}
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateRandomUser)
	handler.ServeHTTP(responseRecorder, req)
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	usr := users.User{}
	if err := json.NewDecoder(responseRecorder.Body).Decode(&usr); err != nil {
		log.Fatalln(err)
	}
	if usr.FirstName == "" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			responseRecorder.Body.String(), usr)
	}
}
