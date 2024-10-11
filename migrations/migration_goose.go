package migrations

import (
	"database/sql"
)

func Up(tx *sql.Tx) error {
	// Thực hiện các câu lệnh SQL để tạo bảng
	_, err := tx.Exec(`CREATE TABLE users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL UNIQUE
    );`)
	return err
}

func Down(tx *sql.Tx) error {
	// Thực hiện các câu lệnh SQL để xóa bảng
	_, err := tx.Exec(`DROP TABLE users;`)
	return err
}
