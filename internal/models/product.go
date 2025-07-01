package models

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Size        int     `json:"size"`
	Category_id int     `json:"category_Id"`
	ImageURL    string  `json:"imageURL"`
}
