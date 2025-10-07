package db

import (
	"WebSportwareShop/internal/models"
	"context"
)

func CreateCategory(ctx context.Context, c *models.Category) error {
	_, err := db.ExecContext(ctx, "INSERT INTO categories (name, description)VALUES($1,$2)", c.Name, c.Description)
	return err
}

func DeleteCategory(ctx context.Context, id int) error {
	_, err := db.ExecContext(ctx, "DELETE FROM categories WHERE id=$1", id)
	return err
}
func ListOfCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := db.QueryContext(ctx, "SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := []models.Category{}
	for rows.Next() {
		var c models.Category
		if err = rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}
func UpdateCategory(ctx context.Context, c *models.Category) error {
	_, err := db.ExecContext(ctx, "UPDATE categories SET name= $2, description=$3 WHERE id= $1", c.ID, c.Name, c.Description)
	return err
}
func GetCategory(ctx context.Context, id int) (models.Category, error) {
	var p models.Category
	err := db.QueryRowContext(ctx, "SELECT * FROM categories WHERE id=$1", id).Scan(&p.ID, &p.Name, &p.Description)
	if err != nil {
		return models.Category{}, err
	}
	return p, err
}
