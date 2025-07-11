const API_BASE_URL = "http://localhost:8080";

// Fetch wrapper with JWT Authorization header
function apiRequest(endpoint, method = "GET", body = null) {
    const headers = new Headers();
    headers.append("Content-Type", "application/json");

    const token = localStorage.getItem("token");
    if (!token) {
        throw new Error("No authentication token found");
    }

    const options = {
        method,
        headers,
        credential: "include",
    };

    if (body) {
        options.body = JSON.stringify(body);
    }

    return fetch(`${API_BASE_URL}${endpoint}`, options)
        .then((response) => {
            if (!response.ok) {
                return response.json().then((data) => {
                    throw new Error(data.error || "Something went wrong");
                });
            }
            return response.json();
        })
        .catch((error) => {
            console.error("API Error:", error.message);
            throw error;
        });
}

// Product APIs - Fixed endpoints to match backend
function getProducts() {
    return apiRequest("/products");
}

function createProduct(product) {
    return apiRequest("/products", "POST", product);
}

function updateProduct(id, product) {
    return apiRequest(`/products/${id}`, "PUT", product);
}

function deleteProduct(id) {
    return apiRequest(`/products/${id}`, "DELETE");
}

function getProduct(id) {
    return apiRequest(`/products/${id}`);
}

// Category APIs
function getCategories() {
    return apiRequest("/categories");
}

function createCategory(category) {
    return apiRequest("/categories", "POST", category);
}

function updateCategory(id, category) {
    return apiRequest(`/categories/${id}`, "PUT", category);
}

function deleteCategory(id) {
    return apiRequest(`/categories/${id}`, "DELETE");
}

function getCategory(id) {
    return apiRequest(`/categories/${id}`);
}

// User APIs - Fixed endpoints to match backend
function loginUser(credentials) {
    return apiRequest("/login", "POST", credentials);
}

function signupUser(user) {
    return apiRequest("/signup", "POST", user);
}

function getUsers() {
    return apiRequest("/users");
}

function getUser(id) {
    return apiRequest(`/users/${id}`);
}

function updateUser(id, user) {
    return apiRequest(`/users/${id}`, "PUT", user);
}

function deleteUser(id) {
    return apiRequest(`/users/${id}`, "DELETE");
}

// Helper function to make authenticated requests
function makeAuthenticatedRequest(endpoint, method = "GET", body = null) {
    const token = localStorage.getItem("token");

    if (!token) {
        throw new Error("No authentication token found");
    }

    const headers = {
        "Content-Type": "application/json",
        credential : "include",
    };

    const options = {
        method,
        headers
    };

    if (body) {
        options.body = JSON.stringify(body);
    }

    return fetch(`${API_BASE_URL}${endpoint}`, options)
        .then(response => {
            if (!response.ok) {
                if (response.status === 401) {
                    // Token expired or invalid
                    localStorage.clear();
                    window.location.href = '/index.html';
                    throw new Error("Session expired. Please login again.");
                }
                return response.json().then(data => {
                    throw new Error(data.error || `HTTP ${response.status}`);
                });
            }
            return response.json();
        });
}

// Expose functions for other scripts
window.api = {
    // Generic request function
    apiRequest,
    makeAuthenticatedRequest,

    // Product functions
    getProducts,
    createProduct,
    updateProduct,
    deleteProduct,
    getProduct,

    // Category functions
    getCategories,
    createCategory,
    updateCategory,
    deleteCategory,
    getCategory,

    // User functions
    loginUser,
    signupUser,
    getUsers,
    getUser,
    updateUser,
    deleteUser
};