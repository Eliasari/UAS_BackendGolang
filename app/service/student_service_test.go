package service

import (
	"database/sql"
	"net/http/httptest"
	"strings"
	"testing"

	"uas-prestasi/app/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetStudents_Success(t *testing.T) {
	// 1️⃣ setup mock DB
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// 2️⃣ mock query result
	rows := sqlmock.NewRows([]string{"id", "user_id", "advisor_id"}).
		AddRow("s1", "u1", nil)

	mock.ExpectQuery(`SELECT id, user_id, advisor_id FROM students`).
		WillReturnRows(rows)

	// 3️⃣ init repo & service
	studentRepo := repository.NewStudentRepository(db)
	service := NewStudentService(studentRepo, nil)

	// 4️⃣ fiber app
	app := fiber.New()
	app.Get("/students", service.GetStudents)

	// 5️⃣ execute request
	req := httptest.NewRequest("GET", "/students", nil)
	resp, err := app.Test(req)

	// 6️⃣ assertion
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetStudent_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery(`SELECT id, user_id, advisor_id FROM students WHERE id=\$1`).
		WithArgs("s1").
		WillReturnError(sql.ErrNoRows)

	studentRepo := repository.NewStudentRepository(db)
	service := NewStudentService(studentRepo, nil)

	app := fiber.New()
	app.Get("/students/:id", service.GetStudent)

	req := httptest.NewRequest("GET", "/students/s1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestGetStudent_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "user_id", "advisor_id"}).
		AddRow("s1", "u1", nil)

	mock.ExpectQuery(`SELECT id, user_id, advisor_id FROM students WHERE id=\$1`).
		WithArgs("s1").
		WillReturnRows(rows)

	studentRepo := repository.NewStudentRepository(db)
	service := NewStudentService(studentRepo, nil)

	app := fiber.New()
	app.Get("/students/:id", service.GetStudent)

	req := httptest.NewRequest("GET", "/students/s1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestAssignAdvisor_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec(`UPDATE students SET advisor_id=\$1 WHERE id=\$2`).
		WithArgs("l1", "s1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	studentRepo := repository.NewStudentRepository(db)
	service := NewStudentService(studentRepo, nil)

	app := fiber.New()
	app.Post("/students/:id/advisor", service.AssignAdvisor)

	body := `{"advisor_id":"l1"}`
	req := httptest.NewRequest("POST", "/students/s1/advisor", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestAssignAdvisor_InvalidPayload(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	studentRepo := repository.NewStudentRepository(db)
	service := NewStudentService(studentRepo, nil)

	app := fiber.New()
	app.Post("/students/:id/advisor", service.AssignAdvisor)

	req := httptest.NewRequest("POST", "/students/s1/advisor", strings.NewReader(`{`))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

