package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/mosesbenjamin/sre-from-local-to-prod/internal/database"
)

type Student struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
}

func databaseStudentToStudent(student database.Student) Student {
	return Student{
		ID:        student.ID,
		CreatedAt: student.CreatedAt,
		UpdatedAt: student.UpdatedAt,
		Email:     student.Email,
	}
}

func databaseStudentsToStudents(students []database.Student) []Student {
	result := make([]Student, len(students))
	for i, student := range students {
		result[i] = databaseStudentToStudent(student)
	}
	return result
}
