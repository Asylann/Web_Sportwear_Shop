const API_BASE_URL = "http://localhost:8080";

// Fetch wrapper with JWT Authorization header
function apiRequest(endpoint, method = "GET", body = null) {
    const headers = new Headers();
    headers.append("Content-Type", "application/json");

    const token = localStorage.getItem("token");
    if (token) {
        headers.append("Authorization", `Bearer ${token}`);
    }

    const options = {
        method,
        headers,
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

// Product APIs
function getProducts() {
    return apiRequest("/products");
}

function createProduct(product) {
    return apiRequest("/createProduct", "POST", product);
}

function deleteProduct(id) {
    return apiRequest(`/deleteProduct/${id}`, "DELETE");
}

// User APIs
function loginUser(credentials) {
    return apiRequest("/login", "POST", credentials);
}

function signupUser(user) {
    return apiRequest("/signup", "POST", user);
}

function getUsers() {
    return apiRequest("/users");
}

// Expose functions for other scripts
window.api = {
    getProducts,
    createProduct,
    deleteProduct,
    loginUser,
    signupUser,
    getUsers,
};
