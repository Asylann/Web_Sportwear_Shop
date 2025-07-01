package db

import (
	"WebSportwareShop/internal/models"
	"context"
)

func GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var u models.User
	err := stmtGetUserByEmail.QueryRowContext(ctx, email).Scan(&u.ID, &u.Email, &u.Password, &u.RoleId)
	if err != nil {
		return models.User{}, err
	}
	return u, err
}

func CreateUser(ctx context.Context, u *models.User) error {
	_, err := stmtCreateUser.ExecContext(ctx, u.Email, u.Password, u.RoleId)
	return err
}

func DeleteUser(ctx context.Context, id int) error {
	_, err := stmtDeleteUser.ExecContext(ctx, id)
	return err
}
func ListOfUsers(ctx context.Context) ([]models.User, error) {
	rows, err := stmtListOfUser.QueryContext(ctx)
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
	_, err := stmtUpdateUser.ExecContext(ctx, u.ID, u.Email, u.Password, u.RoleId)
	return err
}
func GetUser(ctx context.Context, id int) (models.User, error) {
	var u models.User
	err := stmtGetUser.QueryRowContext(ctx, id).Scan(&u.ID, &u.Email, &u.Password, &u.RoleId)
	if err != nil {
		return models.User{}, err
	}
	return u, err
}
