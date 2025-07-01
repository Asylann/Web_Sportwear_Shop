package db

import (
	"WebSportwareShop/internal/models"
	"context"
)

func CreateProduct(ctx context.Context, p *models.Product) error {
	_, err := stmtCreateProduct.ExecContext(ctx, p.Name, p.Description, p.Category_id, p.Size, p.Price, p.ImageURL)
	return err
}

func DeleteProduct(ctx context.Context, id int) error {
	_, err := stmtDeleteProduct.ExecContext(ctx, id)
	return err
}
func ListOfProducts(ctx context.Context) ([]models.Product, error) {
	rows, err := stmtListOfProduct.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		if err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Size, &p.Price, &p.ImageURL, &p.Category_id); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}
func UpdateProduct(ctx context.Context, p *models.Product) error {
	_, err := stmtUpdateProduct.ExecContext(ctx, p.ID, p.Name, p.Description, p.Size, p.Price, p.ImageURL, p.Category_id)
	return err
}
func GetProduct(ctx context.Context, id int) (models.Product, error) {
	var p models.Product
	err := stmtGetProduct.QueryRowContext(ctx, id).Scan(&p.ID, &p.Name, &p.Description, &p.Size, &p.Price, &p.ImageURL, &p.Category_id)
	if err != nil {
		return models.Product{}, err
	}
	return p, err
}
