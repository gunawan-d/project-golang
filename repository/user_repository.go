package repository

import (
	"log"

	"github.com/jmoiron/sqlx"
	"project-golang/models"
)

const (
	// Query untuk mengambil data user berdasarkan email
	QueryGetUserByEmail = `SELECT name, email, idcard FROM test_profile WHERE email = ?`
	// Query untuk memperbarui atau menyisipkan ID Card
	QueryInsertIDCard = `INSERT INTO test_profile (name, email, idcard) VALUES (?, ?, ?)`
)

type UserRepository struct {
	DB *sqlx.DB
}

// GetUserByEmail mengambil data user berdasarkan email
func (repo *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := models.User{}
	err := repo.DB.QueryRow(QueryGetUserByEmail, email).Scan(&user.Name, &user.Email, &user.IDCard)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		return nil, err
	}
	return &user, nil
}

// UpdateIDCard menyisipkan atau memperbarui data ID Card ke dalam database
func (repo *UserRepository) UpdateIDCard(req models.UpdateIDCardRequest) error {
	_, err := repo.DB.Exec(QueryInsertIDCard, req.Name, req.Email, req.IDCard)
	if err != nil {
		log.Printf("Error updating ID card: %v", err)
	}
	return err
}
