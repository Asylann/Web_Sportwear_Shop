// frontend/js/seller.js

// ensure this matches the form id in your HTML
const FORM_ID = "add-product-form";


document.addEventListener("DOMContentLoaded", () => {
    utils.checkAuth();
    utils.checkRole(2, 3);
    loadCategories();
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
        .then(res => res.json())
        .then(json => {
            const select = document.getElementById("product-category");
            json.data.forEach(cat => {
                const opt = document.createElement("option");
                opt.value = cat.id;
                opt.textContent = cat.name;
                select.appendChild(opt);
            });
        })
        .catch(err => console.error("Failed to load categories:", err));
}

function handleProductCreation(event) {
    event.preventDefault();
    console.log("[seller.js] handleProductCreation fired");

    const token = localStorage.getItem("token");
    const name = document.getElementById("product-name").value;
    const description = document.getElementById("product-description").value;
    const price = parseFloat(document.getElementById("product-price").value);
    const imageURL = document.getElementById("product-image").value;
    const categoryId = parseInt(document.getElementById("product-category").value, 10);

    if (!categoryId) {
        alert("Please select a category!");
        return;
    }

    const payload = { name, description, price, category_id: categoryId, imageURL };
    console.log("→ Sending payload:", payload);

    fetch("http://localhost:8080/products", {  // or "/products" if that's your route
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
            if (!res.ok) throw new Error(body.error || `HTTP ${res.status}`);
            alert("Product created!");
            document.getElementById(FORM_ID).reset();
        })
        .catch(err => {
            console.error("Error creating product:", err);
            alert("Error: " + err.message);
        });
}
