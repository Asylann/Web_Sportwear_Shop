package db

import (
	"WebSportwareShop/internal/models"
	"context"
)

func CreateCategory(ctx context.Context, c *models.Category) error {
	_, err := stmtCreateCategory.ExecContext(ctx, c.Name, c.Description)
	return err
}

func DeleteCategory(ctx context.Context, id int) error {
	_, err := stmtDeleteCategory.ExecContext(ctx, id)
	return err
}
func ListOfCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := stmtListOfCategory.QueryContext(ctx)
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
	_, err := stmtUpdateCategory.ExecContext(ctx, c.ID, c.Name, c.Description)
	return err
}
func GetCategory(ctx context.Context, id int) (models.Category, error) {
	var p models.Category
	err := stmtGetCategory.QueryRowContext(ctx, id).Scan(&p.ID, &p.Name, &p.Description)
	if err != nil {
		return models.Category{}, err
	}
	return p, err
}
