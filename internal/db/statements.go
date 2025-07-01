package db

import (
	"context"
	"database/sql"
	"log"
)

var (
	stmtCreateProduct *sql.Stmt
	stmtDeleteProduct *sql.Stmt
	stmtListOfProduct *sql.Stmt
	stmtUpdateProduct *sql.Stmt
	stmtGetProduct    *sql.Stmt

	stmtCreateCategory *sql.Stmt
	stmtDeleteCategory *sql.Stmt
	stmtListOfCategory *sql.Stmt
	stmtUpdateCategory *sql.Stmt
	stmtGetCategory    *sql.Stmt

	stmtCreateUser     *sql.Stmt
	stmtDeleteUser     *sql.Stmt
	stmtListOfUser     *sql.Stmt
	stmtUpdateUser     *sql.Stmt
	stmtGetUser        *sql.Stmt
	stmtGetUserByEmail *sql.Stmt
)

func initStmt(db *sql.DB) {
	var err error
	stmtCreateProduct, err = db.PrepareContext(context.Background(),
		`INSERT INTO products (name, description, category_id,size, price, imageurl)
	VALUES($1,$2,$3,$4,$5,$6)`)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	stmtDeleteProduct, err = db.PrepareContext(context.Background(),
		`DELETE FROM products WHERE id=$1`)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	stmtListOfProduct, err = db.PrepareContext(context.Background(),
		`SELECT * FROM products`)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	stmtUpdateProduct, err = db.PrepareContext(context.Background(),
		`UPDATE products SET name= $2, description=$3,size=$4, price=$5, imageurl=$6,category_id=$7 WHERE id= $1`)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	stmtGetProduct, err = db.PrepareContext(context.Background(),
		`SELECT * FROM products WHERE id=$1`)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	stmtCreateCategory, err = db.PrepareContext(context.Background(),
		`INSERT INTO categories (name, description)
	VALUES($1,$2)`)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	stmtDeleteCategory, err = db.PrepareContext(context.Background(),
		`DELETE FROM categories WHERE id=$1`)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	stmtListOfCategory, err = db.PrepareContext(context.Background(),
		`SELECT * FROM categories`)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	stmtUpdateCategory, err = db.PrepareContext(context.Background(),
		`UPDATE categories SET name= $2, description=$3 WHERE id= $1`)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	stmtGetCategory, err = db.PrepareContext(context.Background(),
		`SELECT * FROM categories WHERE id=$1`)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	stmtCreateUser, err = db.PrepareContext(context.Background(),
		`INSERT INTO users (email, password, role_id)
	VALUES($1,$2,$3)`)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	stmtDeleteUser, err = db.PrepareContext(context.Background(),
		`DELETE FROM users WHERE id=$1`)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	stmtListOfUser, err = db.PrepareContext(context.Background(),
		`SELECT * FROM users`)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	stmtUpdateUser, err = db.PrepareContext(context.Background(),
		`UPDATE users SET email= $2, password=$3, role_id=$4 WHERE id= $1`)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	stmtGetUser, err = db.PrepareContext(context.Background(),
		`SELECT * FROM users WHERE id=$1`)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	stmtGetUserByEmail, err = db.PrepareContext(context.Background(),
		`SELECT * FROM users WHERE email=$1`)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
