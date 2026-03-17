package storage

import (
	"logger-app/db"
	"logger-app/models"
)

func CreateUser(id, name string) error {
	_, err := db.DB.Exec("INSERT INTO users (id, name) VALUES (?, ?)", id, name)
	return err
}

func GetAllUsers() ([]models.User, error) {
	rows, err := db.DB.Query("SELECT id, name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetUserByID(id string) (*models.User, error) {
	var u models.User
	err := db.DB.QueryRow("SELECT id, name FROM users WHERE id = ?", id).Scan(&u.ID, &u.Name)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func UpdateUser(id, name string) error {
	_, err := db.DB.Exec("UPDATE users SET name = ? WHERE id = ?", name, id)
	return err
}

func DeleteUser(id string) error {
	// First delete events for this user to maintain referential integrity
	_, err := db.DB.Exec("DELETE FROM events WHERE user_id = ?", id)
	if err != nil {
		return err
	}
	_, err = db.DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func UserExists(userID string) (bool, error) {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
