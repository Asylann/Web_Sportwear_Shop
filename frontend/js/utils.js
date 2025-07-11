// frontend/js/utils.js

// Show alert messages to user
function showAlert(message, type = "success") {
    const alertContainer = document.createElement("div");
    alertContainer.className = `alert alert-${type}`;
    alertContainer.textContent = message;
    alertContainer.style.position = "fixed";
    alertContainer.style.top = "20px";
    alertContainer.style.right = "20px";
    alertContainer.style.zIndex = "1000";
    alertContainer.style.padding = "10px 20px";
    alertContainer.style.borderRadius = "4px";
    alertContainer.style.color = "white";
    alertContainer.style.backgroundColor = type === "success" ? "#4CAF50" : "#f44336";

    document.body.appendChild(alertContainer);

    setTimeout(() => {
        alertContainer.remove();
    }, 3000);
}

// Redirect to a different page
function redirectTo(page) {
    window.location.href = page;
}

// Check if user has required role(s) with better error handling
function checkRole(...allowedRoles) {
    console.log("Checking role access...");

    const token = localStorage.getItem("token");
    const roleId = parseInt(localStorage.getItem("roleId"), 10);

    console.log("Token exists:", !!token);
    console.log("Role ID:", roleId);
    console.log("Allowed roles:", allowedRoles);

    if (!token) {
        console.error("No token found");
        showAlert("Please log in to access this page.", "error");
        redirectToLogin();
        return false;
    }

    if (!roleId || isNaN(roleId)) {
        console.error("Invalid role ID");
        showAlert("Invalid role. Please log in again.", "error");
        clearToken();
        redirectToLogin();
        return false;
    }

    if (!allowedRoles.includes(roleId)) {
        console.error(`Role ${roleId} not in allowed roles:`, allowedRoles);
        showAlert("You don't have permission to access this page.", "error");
        setTimeout(() => {
            redirectToDashboard();
        }, 2000);
        return false;
    }

    console.log("Role check passed");
    return true;
}

// Check if user is authenticated with better token validation
function checkAuth() {
    console.log("Checking authentication...");

    const token = localStorage.getItem("token");

    if (!token) {
        console.error("No token found");
        showAlert("Please log in to access this page.", "error");
        redirectToLogin();
        return false;
    }

    // Check if token is expired
    try {
        const payload = JSON.parse(atob(token.split(".")[1]));
        const currentTime = Math.floor(Date.now() / 1000);

        if (payload.exp && payload.exp < currentTime) {
            console.error("Token expired");
            showAlert("Your session has expired. Please log in again.", "error");
            clearToken();
            redirectToLogin();
            return false;
        }
    } catch (error) {
        console.error("Error parsing token:", error);
        showAlert("Invalid token. Please log in again.", "error");
        clearToken();
        redirectToLogin();
        return false;
    }

    console.log("Authentication check passed");
    return true;
}

// Helper function to redirect to login
function redirectToLogin() {
    const currentPath = window.location.pathname;
    if (currentPath.includes('/pages/')) {
        window.location.href = "../index.html";
    } else {
        window.location.href = "index.html";
    }
}

// Helper function to redirect to dashboard
function redirectToDashboard() {
    const currentPath = window.location.pathname;
    if (currentPath.includes('/pages/')) {
        window.location.href = "dashboard.html";
    } else {
        window.location.href = "pages/dashboard.html";
    }
}

// Save token to localStorage
function saveToken(token) {
    localStorage.setItem("token", token);
}

// Get token from localStorage
function getToken() {
    return localStorage.getItem("token");
}

// Remove token from localStorage
function clearToken() {
    localStorage.removeItem("token");
    localStorage.removeItem("roleId");
    console.log("Tokens cleared");
}

// Check if token exists
function isAuthenticated() {
    return !!getToken();
}

// Global logout function
function logout() {
    console.log("Logging out...");
    localStorage.removeItem("token");
    localStorage.removeItem("role");
    localStorage.removeItem("email");
    localStorage.removeItem("userId");
    fetch("http://localhost:8080/logout", {
        method: "POST",
        credentials: "include"
    }).then(() => {
        showAlert("You have been logged out.", "success");
        setTimeout(() => {
            redirectToLogin();
        }, 1000);
    });
    window.location.href = "/index.html";
}


// Attach event listener helper
function onClick(selector, callback) {
    const element = document.querySelector(selector);
    if (element) {
        element.addEventListener("click", callback);
    }
}

// Attach form submit handler helper
function onSubmit(selector, callback) {
    const form = document.querySelector(selector);
    if (form) {
        form.addEventListener("submit", callback);
    }
}

// Enhanced API request with better error handling
function makeAuthenticatedRequest(endpoint, options = {}) {
    return fetch(`http://localhost:8080${endpoint}`, {
        ...options,
        credentials: "include" // <— just this
    }).then(response => {
        if (response.status === 401) {
            console.error("Unauthorized access - not logged in");
            utils.clearToken();
            utils.showAlert("Session expired. Please log in again.", "error");
            setTimeout(() => {
                redirectToLogin();
            }, 2000);
            throw new Error("Unauthorized");
        }

        if (response.status === 403) {
            console.error("Forbidden access - insufficient permissions");
            utils.showAlert("You don't have permission to perform this action.", "error");
            throw new Error("Forbidden");
        }

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}`);
        }

        return response.json();
    });
}


// Expose utilities globally
window.utils = {
    showAlert,
    redirectTo,
    saveToken,
    getToken,
    clearToken,
    checkRole,
    checkAuth,
    isAuthenticated,
    onClick,
    onSubmit,
    logout,
    makeAuthenticatedRequest,
    redirectToLogin,
    redirectToDashboard
};

// Also expose logout globally for easy access
window.logout = logout;