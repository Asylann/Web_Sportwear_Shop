package db

import (
	"WebSportwareShop/internal/models"
	"context"
)

func GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var u models.User
	err := db.QueryRowContext(ctx, "SELECT * FROM users WHERE email=$1", email).Scan(&u.ID, &u.Email, &u.Password, &u.RoleId)
	if err != nil {
		return models.User{}, err
	}
	return u, err
}

func CreateUser(ctx context.Context, u *models.User) (int, error) {
	var id int
	err := db.QueryRowContext(ctx, "INSERT INTO users (email, password, role_id) VALUES($1,$2,$3) RETURNING id", u.Email, u.Password, u.RoleId).Scan(&id)
	u.ID = id
	return id, err
}

func DeleteUser(ctx context.Context, id int) error {
	_, err := db.ExecContext(ctx, "DELETE FROM users WHERE id=$1", id)
	return err
}
func ListOfUsers(ctx context.Context) ([]models.User, error) {
	rows, err := db.QueryContext(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := []models.User{}
	for rows.Next() {
		var u models.User
		if err = rows.Scan(&u.ID, &u.Email, &u.Password, &u.RoleId); err != nil {
			return nil, err
		}
		categories = append(categories, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}
func UpdateUser(ctx context.Context, u *models.User) error {
	_, err := db.ExecContext(ctx, "UPDATE users SET email= $2, password=$3, role_id=$4 WHERE id= $1", u.ID, u.Email, u.Password, u.RoleId)
	return err
}
func GetUser(ctx context.Context, id int) (models.User, error) {
	var u models.User
	err := db.QueryRowContext(ctx, "SELECT * FROM users WHERE id=$1", id).Scan(&u.ID, &u.Email, &u.Password, &u.RoleId)
	if err != nil {
		return models.User{}, err
	}
	return u, err
}

func GetUserEmail(ctx context.Context, id int) (string, error) {
	var u models.User
	err := db.QueryRowContext(ctx, "SELECT email FROM users WHERE ID=$1", id).Scan(&u.Email)
	if err != nil {
		return "", err
	}
	return u.Email, err
}

func GetEtagVersionByName(ctx context.Context, name string) (int, error) {
	var version int
	err := db.QueryRowContext(ctx, "SELECT version FROM etag_versions WHERE name=$1", name).Scan(&version)
	return version, err
}

func ChangeEtagVersionByName(ctx context.Context, name string) error {
	var version int
	err := db.QueryRowContext(ctx, "SELECT version FROM etag_versions WHERE name=$1", name).Scan(&version)
	if err != nil {
		return err
	}
	version++
	_, err = db.ExecContext(ctx, "UPDATE etag_versions SET version=$2 WHERE name=$1", name, version)
	return err
}
