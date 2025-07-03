// frontend/js/utils.js

// Show alert messages to user
function showAlert(message, type = "success") {
    const alertContainer = document.createElement("div");
    alertContainer.className = `alert alert-${type}`;
    alertContainer.textContent = message;

    document.body.appendChild(alertContainer);

    setTimeout(() => {
        alertContainer.remove();
    }, 3000);
}

// Redirect to a different page
function redirectTo(page) {
    window.location.href = page;
}

function checkRole(...allowedRoles) {
    const roleId = parseInt(localStorage.getItem("roleId"), 10);
    if (!allowedRoles.includes(roleId)) {
        alert("You don't have permission for this page.");
        window.location.href = "../pages/dashboard.html";
    }
}

function checkAuth() {
    const token = localStorage.getItem("token");
    if (!token) {
        alert("You must be logged in to access this page.");
        window.location.href = "../index.html";
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
}

// Check if token exists
function isAuthenticated() {
    return !!getToken();
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
    isAuthenticated,
    onClick,
    onSubmit,
};
