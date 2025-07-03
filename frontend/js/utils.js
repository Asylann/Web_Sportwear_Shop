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

// Check if user has required role(s)
function checkRole(...allowedRoles) {
    const roleId = parseInt(localStorage.getItem("roleId"), 10);
    if (!roleId || !allowedRoles.includes(roleId)) {
        alert("You don't have permission for this page.");
        // Check if we're already in pages/ directory
        const currentPath = window.location.pathname;
        if (currentPath.includes('/pages/')) {
            window.location.href = "dashboard.html";
        } else {
            window.location.href = "pages/dashboard.html";
        }
        return false;
    }
    return true;
}

// Check if user is authenticated
function checkAuth() {
    const token = localStorage.getItem("token");
    if (!token) {
        alert("You must be logged in to access this page.");
        // Check if we're already in pages/ directory
        const currentPath = window.location.pathname;
        if (currentPath.includes('/pages/')) {
            window.location.href = "../index.html";
        } else {
            window.location.href = "index.html";
        }
        return false;
    }
    return true;
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
}

// Check if token exists
function isAuthenticated() {
    return !!getToken();
}

// Global logout function
function logout() {
    clearToken();
    // Check if we're in pages/ directory
    const currentPath = window.location.pathname;
    if (currentPath.includes('/pages/')) {
        window.location.href = "../index.html";
    } else {
        window.location.href = "index.html";
    }
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
    logout
};

// Also expose logout globally for easy access
window.logout = logout;