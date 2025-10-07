package db

import (
	"WebSportwareShop/internal/models"
	"context"
)

func CreateProduct(ctx context.Context, p *models.Product) error {
	_, err := db.ExecContext(ctx, "INSERT INTO products (name, description, category_id,size, price, imageurl, seller_id)\n\tVALUES($1,$2,$3,$4,$5,$6,$7)", p.Name, p.Description, p.CategoryID, p.Size, p.Price, p.ImageURL, p.SellerID)
	return err
}

func DeleteProduct(ctx context.Context, id int) error {
	_, err := db.ExecContext(ctx, "DELETE FROM products WHERE id=$1", id)
	return err
}
func ListOfProducts(ctx context.Context) ([]models.Product, error) {
	rows, err := db.QueryContext(ctx, "SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Size, &p.Price, &p.ImageURL, &p.CategoryID, &p.SellerID); err != nil {
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
	_, err := db.ExecContext(ctx, "UPDATE products SET name= $2, description=$3,size=$4, price=$5, imageurl=$6,category_id=$7, seller_id=$8 WHERE id= $1", p.ID, p.Name, p.Description, p.Size, p.Price, p.ImageURL, p.CategoryID, p.SellerID)
	return err
}
func GetProduct(ctx context.Context, id int) (models.Product, error) {
	var p models.Product
	err := db.QueryRowContext(ctx, "SELECT * FROM products WHERE id=$1", id).Scan(&p.ID, &p.Name, &p.Description, &p.Size, &p.Price, &p.ImageURL, &p.CategoryID, &p.SellerID)
	if err != nil {
		return models.Product{}, err
	}
	return p, err
}

func ListOfProductsByCategory(ctx context.Context, id int) ([]models.Product, error) {
	rows, err := db.QueryContext(ctx, "SELECT * FROM products WHERE category_id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Size, &p.Price, &p.ImageURL, &p.CategoryID, &p.SellerID); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func ListOfProductsBySellerID(ctx context.Context, id int) ([]models.Product, error) {
	rows, err := db.QueryContext(ctx, "SELECT * FROM products WHERE seller_id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Size, &p.Price, &p.ImageURL, &p.CategoryID, &p.SellerID); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
