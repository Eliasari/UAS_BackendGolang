package service

import (
	"net/http/httptest"
	"testing"

	"uas-prestasi/app/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetLecturers_Success(t *testing.T) {
	// mock DB
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// mock rows
	rows := sqlmock.NewRows([]string{"id", "user_id", "username"}).
		AddRow("l1", "u1", "dosenA")

	mock.ExpectQuery(`
		SELECT l.id, l.user_id, u.username
		FROM lecturers l
		JOIN users u ON u.id = l.user_id
	`).WillReturnRows(rows)

	// init repo & service
	lecturerRepo := repository.NewLecturerRepository(db)
	service := NewStudentService(nil, lecturerRepo)

	// fiber app
	app := fiber.New()
	app.Get("/lecturers", service.GetLecturers)

	// request
	req := httptest.NewRequest("GET", "/lecturers", nil)
	resp, err := app.Test(req)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetLecturers_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery(`FROM lecturers`).
		WillReturnError(assert.AnError)

	lecturerRepo := repository.NewLecturerRepository(db)
	service := NewStudentService(nil, lecturerRepo)

	app := fiber.New()
	app.Get("/lecturers", service.GetLecturers)

	req := httptest.NewRequest("GET", "/lecturers", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

func TestGetAdvisees_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery(`FROM students`).
		WithArgs("l1").
		WillReturnError(assert.AnError)

	lecturerRepo := repository.NewLecturerRepository(db)
	service := NewStudentService(nil, lecturerRepo)

	app := fiber.New()
	app.Get("/lecturers/:id/advisees", service.GetAdvisees)

	req := httptest.NewRequest("GET", "/lecturers/l1/advisees", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

