package repository

import (
	"fmt"
	"log"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type userRepositoryDB struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryDB{db: db}
}

// NewDatabase creates a new database connection
func NewDatabase(host string, port int, user, password, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Successfully connected to the database.")
	return db, nil
}

// CloseDatabase closes the database connection
func CloseDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get DB instance: %w", err)
	}
	return sqlDB.Close()
}

func (r *userRepositoryDB) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

func (r *userRepositoryDB) Login(username, password string) (bool, error) {
	var user User
	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false, nil
	}

	return true, nil
}

func (r *userRepositoryDB) GetAll() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepositoryDB) GetUser(identifier string) (*User, error) {
	var user User

	// Check if identifier is a UUID
	if userUUID, parseErr := uuid.FromString(identifier); parseErr == nil {
		// Find user by UUID
		if err := r.db.First(&user, "uuid = ?", userUUID).Error; err != nil {
			return nil, err
		}
	} else {
		// Find user by username
		if err := r.db.First(&user, "username = ?", identifier).Error; err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (r *userRepositoryDB) UpdateUser(userUUID uuid.UUID, user *User) error {
	return r.db.Model(&User{}).Where("uuid = ?", userUUID).Updates(user).Error
}

func (r *userRepositoryDB) DeleteUser(identifier string) error {
	// Check if identifier is a UUID
	if userUUID, parseErr := uuid.FromString(identifier); parseErr == nil {
		// Delete user by UUID
		return r.db.Where("uuid = ?", userUUID).Delete(&User{}).Error
	} else {
		// Delete user by username
		return r.db.Where("username = ?", identifier).Delete(&User{}).Error
	}
}
