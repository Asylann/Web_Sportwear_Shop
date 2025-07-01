# WEB SHOPWEAR SHOP

### It demonstrates a production‑ready REST API with user authentication (JWT), role‑based access control (customer, seller, admin), product & category management, shopping cart & order workflows, and schema migrations.

## Key Features

* #### User Management: Signup, login (JWT), and profile endpoints

* #### RBAC: Customers browse and purchase; sellers manage their products; admins manage users and orders

* #### Product Catalog: CRUD operations on products and categories

* #### Database Migrations: Versioned SQL migrations powered by golang-migrate

## Tech Stack

* #### Backend: Go, Gorilla Mux, database/sql + lib/pq, JWT (github.com/golang-jwt/jwt)

* #### Database: PostgreSQL with prepared statements and connection pooling

* #### Migrations: golang-migrate