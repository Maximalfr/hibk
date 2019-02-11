package database

import (
	"database/sql"
	"log"

	"github.com/Maximalfr/hibk/pkg/hibk/models"
)

// initUser initializes the database for the user table
func initUser(db *sql.DB) {
	users := `CREATE TABLE IF NOT EXISTS users(
					id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
					username VARCHAR(30) NOT NULL UNIQUE,
                    password CHAR(60) NOT NULL
				)`

	inits := []string{users}

	for _, ex := range inits {
		_, err := db.Exec(ex)
		if err != nil {
			log.Println(err)
		}
	}
}

// GetUser returns a user structure filled with database values
func GetUser(username string) (user models.DbUser, err error) {
	db, err := open()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM users WHERE username=?", username)
	if err != nil {
		log.Println(err)
	}

	// One unique user so one row
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			log.Println(err)
		}
	}
	return
}

// RegisterUser adds a new user in the database.
// The password needs to be hashed
func RegisterUser(username string, password string) error {
	db, err := open()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO users(username, password) VALUES(?,?)", username, password)
	if err != nil {
		return err
	}
	return nil
}

// ChangePassword changes the password for a given user
func ChangePassword(username string, password string) (err error) {
	db, err := open()
	defer db.Close()
	_, err = db.Exec("UPDATE users SET password=? WHERE username=?", password, username)
	return
}
