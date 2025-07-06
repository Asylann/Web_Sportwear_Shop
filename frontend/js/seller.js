// frontend/js/seller.js

// ensure this matches the form id in your HTML
const FORM_ID = "add-product-form";

document.addEventListener("DOMContentLoaded", () => {
    // Check authentication and role
    if (!utils.checkAuth()) return;
    if (!utils.checkRole(2, 3)) return;

    // Load categories and set up form
    loadCategories();
    loadSellerProducts();

    const form = document.getElementById(FORM_ID);
    if (!form) {
        console.error(`Form with id="${FORM_ID}" not found!`);
        return;
    }
    form.addEventListener("submit", handleProductCreation);
});

function loadCategories() {
    const token = localStorage.getItem("token");
    fetch("http://localhost:8080/categories", {
        headers: { "Authorization": "Bearer " + token },
    })
        .then(res => {
            if (!res.ok) {
                throw new Error(`HTTP ${res.status}`);
            }
            return res.json();
        })
        .then(json => {
            const select = document.getElementById("product-category");
            if (!select) {
                console.error("Category select element not found");
                return;
            }

            // Clear existing options except the first one
            select.innerHTML = '<option value="">-- Select Category --</option>';

            if (json.data) {
                json.data.forEach(cat => {
                    const opt = document.createElement("option");
                    opt.value = cat.id;
                    opt.textContent = cat.name;
                    select.appendChild(opt);
                });
            }
        })
        .catch(err => {
            console.error("Failed to load categories:", err);
            utils.showAlert("Failed to load categories", "error");
        });
}

function loadSellerProducts() {
    const token = localStorage.getItem("token");
    const userId = localStorage.getItem("userId")
    fetch("http://localhost:8080/productsBySeller/"+userId, {
        headers: { "Authorization": "Bearer " + token },
    })
        .then(res => {
            if (!res.ok) {
                throw new Error(`HTTP ${res.status}`);
            }
            return res.json();
        })
        .then(json => {
            const container = document.getElementById("seller-product-list");
            if (!container) return;

            container.innerHTML = '';

            if (json.data) {
                const products = Array.isArray(json.data) ? json.data : Object.values(json.data);
                products.forEach(product => {
                    const productCard = createProductCard(product);
                    container.appendChild(productCard);
                });
            }
        })
        .catch(err => {
            console.error("Failed to load products:", err);
            utils.showAlert("Failed to load products", "error");
        });
}

function createProductCard(product) {
    const card = document.createElement("div");
    card.className = "product-card";
    card.innerHTML = `
        ${product.imageURL ? `<img src="${product.imageURL}" alt="${product.name}" />` : ''}
        <div class="info">
            <h3>${product.name}</h3>
            <p>${product.description}</p>
            <div class="price">${product.price} ₸</div>
            <button onclick="deleteProduct(${product.id})" class="button-danger">Delete</button>
        </div>
    `;
    return card;
}

function handleProductCreation(event) {
    event.preventDefault();
    console.log("[seller.js] handleProductCreation fired");

    const token = localStorage.getItem("token");
    const name = document.getElementById("product-name").value.trim();
    const description = document.getElementById("product-description").value.trim();
    const priceInput = document.getElementById("product-price").value;
    const imageURL = document.getElementById("product-image").value.trim();
    const categoryId = parseInt(document.getElementById("product-category").value, 10);
    const sizeInput = document.getElementById("product-size").value;
    const userId = parseInt(localStorage.getItem("userId"),10)
    // Validation
    if (!name || !description || !priceInput) {
        utils.showAlert("Please fill in all required fields", "error");
        return;
    }

    const price = parseFloat(priceInput);
    if (isNaN(price) || price <= 0) {
        utils.showAlert("Please enter a valid price", "error");
        return;
    }
    const size = parseFloat(sizeInput);
    if (isNaN(size) || size < 0) {
        utils.showAlert("Please enter a valid size", "error");
        return;
    }

    if (!categoryId) {
        utils.showAlert("Please select a category", "error");
        return;
    }

    const payload = {
        name,
        description,
        price,
        size : size,
        category_id: categoryId,
        imageURL: imageURL || "",
        seller_id : userId
    };

    console.log("→ Sending payload:", payload);

    fetch("http://localhost:8080/products", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + token
        },
        body: JSON.stringify(payload),
    })
        .then(async res => {
            console.log("← Response status:", res.status);
            const body = await res.json().catch(() => ({}));
            console.log("← Response body:", body);

            if (!res.ok) {
                throw new Error(body.error || `HTTP ${res.status}`);
            }

            utils.showAlert("Product created successfully!");
            document.getElementById(FORM_ID).reset();
            loadSellerProducts(); // Refresh the product list
        })
        .catch(err => {
            console.error("Error creating product:", err);
            utils.showAlert("Error: " + err.message, "error");
        });
}

// Global function to delete products
function deleteProduct(productId) {
    if (!confirm("Are you sure you want to delete this product?")) {
        return;
    }

    const token = localStorage.getItem("token");
    fetch(`http://localhost:8080/products/${productId}`, {
        method: "DELETE",
        headers: {
            "Authorization": "Bearer " + token
        }
    })
        .then(res => {
            if (!res.ok) throw new Error(`HTTP ${res.status}`);
            utils.showAlert('Product deleted successfully!');
            loadSellerProducts();
        })
        .catch(err => {
            console.error('Error deleting product:', err);
            utils.showAlert('Error deleting product', 'error');
        });
}

// Make deleteProduct available globally
window.deleteProduct = deleteProduct;
window.loadCategories = loadCategories;