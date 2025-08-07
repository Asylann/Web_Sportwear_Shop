// apiClient.js

const API_BASE_URL = "https://localhost:8080";

// Generic fetch wrapper that sends the auth cookie automatically
async function apiRequest(endpoint, method = "GET", body = null) {
    const options = {
        method,
        headers: { "Content-Type": "application/json" },
        credentials: "include"    // <-- always send HttpsOnly cookie
    };

    if (body !== null) {
        options.body = JSON.stringify(body);
    }

    const res = await fetch(`${API_BASE_URL}${endpoint}`, options);
    if (!res.ok) {
        // handle 401 by clearing client state and redirecting
        if (res.status === 401) {
            clearSession();
            window.location.href = "/index.html";
            throw new Error("Session expired. Please log in again.");
        }
        const err = await res.json();
        throw new Error(err.error || `HTTPs ${res.status}`);
    }
    return res.json();
}

// --- Auth / User info ---
export async function loginUser(credentials) {
    // POST /login will set the auth_token cookie
    return apiRequest("/login", "POST", credentials);
}

export async function signupUser(user) {
    return apiRequest("/signup", "POST", user);
}

export async function getCurrentUser() {
    // returns { ID, Email, RoleId } from your /me endpoint
    return apiRequest("/me", "GET");
}

// --- Product APIs ---
export const getProducts     = () => apiRequest("/products");
export const createProduct   = (p) => apiRequest("/products", "POST", p);
export const updateProduct   = (id, p) => apiRequest(`/products/${id}`, "PUT", p);
export const deleteProduct   = (id) => apiRequest(`/products/${id}`, "DELETE");
export const getProduct      = (id) => apiRequest(`/products/${id}`);

// --- Category APIs ---
export const getCategories   = () => apiRequest("/categories");
export const createCategory  = (c) => apiRequest("/categories", "POST", c);
export const updateCategory  = (id, c) => apiRequest(`/categories/${id}`, "PUT", c);
export const deleteCategory  = (id) => apiRequest(`/categories/${id}`, "DELETE");
export const getCategory     = (id) => apiRequest(`/categories/${id}`);

// --- User management (admin only) ---
export const getUsers        = () => apiRequest("/users");
export const getUser         = (id) => apiRequest(`/users/${id}`);
export const updateUser      = (id, u) => apiRequest(`/users/${id}`, "PUT", u);
export const deleteUser      = (id) => apiRequest(`/users/${id}`, "DELETE");

// --- Helpers ---
function clearSession() {
    localStorage.removeItem("userId");
    localStorage.removeItem("email");
    localStorage.removeItem("roleId");
}

export default {
    loginUser,
    signupUser,
    getCurrentUser,
    getProducts,
    createProduct,
    updateProduct,
    deleteProduct,
    getProduct,
    getCategories,
    createCategory,
    updateCategory,
    deleteCategory,
    getCategory,
    getUsers,
    getUser,
    updateUser,
    deleteUser,
};
