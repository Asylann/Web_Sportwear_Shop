// frontend/js/admin.js

document.addEventListener("DOMContentLoaded", () => {
    console.log("Admin page loading...");

    // Check authentication and role with detailed logging
    if (!utils.checkAuth()) {
        console.error("Authentication failed");
        return;
    }

    if (!utils.checkRole(3)) {
        console.error("Role check failed for admin");
        return;
    }

    console.log("Admin dashboard loaded successfully.");

    // Initialize admin dashboard
    initializeAdminDashboard();
});

function initializeAdminDashboard() {
    console.log("Initializing admin dashboard...");

    const roleId = parseInt(localStorage.getItem("roleId"), 10);

    console.log("Role ID:", roleId);

    if (roleId !== 3) {
        console.error("Access denied - not admin role");
        utils.showAlert("Access denied. Admin privileges required.", "error");
        setTimeout(() => {
            utils.redirectToDashboard();
        }, 2000);
        return;
    }

    // Display admin info
    displayAdminInfo();

    // Test API connectivity
    testApiConnection();

    console.log("Admin dashboard initialized successfully.");
}

function displayAdminInfo() {
    const userStatusElement = document.getElementById("user-status");
    if (userStatusElement) {
        userStatusElement.textContent = "👤 Admin user authenticated";
    }

    // Add admin greeting
    const adminGreeting = document.createElement("div");
    adminGreeting.innerHTML = `
        <div class="alert alert-success" style="margin-bottom: 20px;">
            <strong>Welcome, Admin!</strong> You have full access to all system features.
        </div>
    `;

    const mainContent = document.querySelector(".dashboard-main");
    if (mainContent) {
        mainContent.insertBefore(adminGreeting, mainContent.firstChild);
    }
}

function testApiConnection() {
    console.log("Testing API connection...");

    // Test with a simple API call
    fetch("https://localhost:8080/categories", {
        credentials : "include",
    })
        .then(response => {
            console.log("API test response status:", response.status);
            if (response.ok) {
                console.log("API connection successful");
            } else {
                console.error("API connection failed:", response.status);
            }
            return response.json();
        })
        .then(data => {
            console.log("API test data:", data);
        })
        .catch(error => {
            console.error("API test error:", error);
            utils.showAlert("API connection failed. Please check your connection.", "error");
        });
}

// Helper function to handle admin-specific errors
function handleAdminError(error, operation) {
    console.error(`Admin operation failed (${operation}):`, error);

    if (error.message.includes("401") || error.message.includes("Unauthorized")) {
        utils.showAlert("Session expired. Please log in again.", "error");
        setTimeout(() => {
            utils.redirectToLogin();
        }, 2000);
    } else if (error.message.includes("403") || error.message.includes("Forbidden")) {
        utils.showAlert("Access denied. Admin privileges required.", "error");
    } else {
        utils.showAlert(`${operation} failed: ${error.message}`, "error");
    }
}

// Helper function to confirm admin actions
function confirmAdminAction(action, callback) {
    const confirmed = confirm(`Are you sure you want to ${action}? This action cannot be undone.`);
    if (confirmed && typeof callback === 'function') {
        callback();
    }
    return confirmed;
}

// Load functions with better error handling
function loadUsers() {
    console.log("Loading users...");

    utils.makeAuthenticatedRequest("/users")
        .then(data => {
            console.log("Users loaded:", data);
            displayUsers(data);
        })
        .catch(error => {
            console.error("Failed to load users:", error);
            handleAdminError(error, "Load users");
        });
}

function loadProducts() {
    console.log("Loading products...");

    utils.makeAuthenticatedRequest("/products")
        .then(data => {
            console.log("Products loaded:", data);
            displayProducts(data);
        })
        .catch(error => {
            console.error("Failed to load products:", error);
            handleAdminError(error, "Load products");
        });
}

function loadCategories() {
    console.log("Loading categories...");

    utils.makeAuthenticatedRequest("/categories")
        .then(data => {
            console.log("Categories loaded:", data);
            displayCategories(data);
        })
        .catch(error => {
            console.error("Failed to load categories:", error);
            handleAdminError(error, "Load categories");
        });
}

function displayUsers(data) {
    const container = document.getElementById("usersList");
    if (!container) return;

    container.innerHTML = '';

    if (data.data && data.data.length > 0) {
        const table = document.createElement("table");
        table.className = "table";
        table.innerHTML = `
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Email</th>
                    <th>Role</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                ${data.data.map(user => `
                    <tr>
                        <td>${user.id}</td>
                        <td>${user.email}</td>
                        <td>
                            <span class="badge ${getBadgeClass(user.role_id)}">
                                ${getRoleName(user.role_id)}
                            </span>
                        </td>
                        <td>
                            <button onclick="deleteUser(${user.id})" class="button-danger">Delete</button>
                        </td>
                    </tr>
                `).join('')}
            </tbody>
        `;
        container.appendChild(table);
    } else {
        container.innerHTML = '<p>No users found.</p>';
    }
}

function displayProducts(data) {
    const container = document.getElementById("productsList");
    if (!container) return;

    container.innerHTML = '';

    if (data.data) {
        const products = Array.isArray(data.data) ? data.data : Object.values(data.data);

        if (products.length > 0) {
            products.forEach(product => {
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
                container.appendChild(card);
            });
        } else {
            container.innerHTML = '<p>No products found.</p>';
        }
    } else {
        container.innerHTML = '<p>No products found.</p>';
    }
}

function displayCategories(data) {
    const select = document.getElementById("category");
    if (!select) return;

    // Clear existing options except the first one
    select.innerHTML = '<option value="">-- Select Category --</option>';

    if (data.data && data.data.length > 0) {
        data.data.forEach(category => {
            const option = document.createElement("option");
            option.value = category.id;
            option.textContent = category.name;
            select.appendChild(option);
        });
    }
}

function getRoleName(roleId) {
    switch (roleId) {
        case 1: return "Customer";
        case 2: return "Seller";
        case 3: return "Admin";
        case "3": return "Admin"
        default: return "Unknown";
    }
}

function getBadgeClass(roleId) {
    switch (roleId) {
        case 1: return "badge-customer";
        case 2: return "badge-seller";
        case 3: return "badge-admin";
        case "3" : return "badge-admin"
        default: return "";
    }
}

// Make functions available globally
window.loadUsers = loadUsers;
window.loadProducts = loadProducts;
window.loadCategories = loadCategories;
window.handleAdminError = handleAdminError;
window.confirmAdminAction = confirmAdminAction;