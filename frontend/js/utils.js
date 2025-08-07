// frontend/js/utils.js


const API_BASEE = "https://localhost:8080";
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
// Call this at the top of any page or route that needs role-based guard
async function checkRole(...allowedRoles) {
    console.log("Checking role access via /me…");

    try {
        // 1) Call /me to validate session
        const res = await fetch(`${API_BASEE}/me`, {
            method: "GET",
            credentials: "include"
        });

        if (!res.ok) {
            console.error("Auth check failed while fetching role:", res.status);
            showAlert("Please log in to access this page.", "error");
            clearSession();
            redirectToLogin();
            return false;
        }

        // 2) Parse full payload
        const raw = await res.json();
        console.log("Raw /me payload:", raw);

        // 3) Extract roleId (handle both camelCase and snake_case)
        const returnedData = raw.data || {};
        const roleId = returnedData.roleId ?? returnedData.role_id;
        console.log("Current RoleId:", roleId);
        console.log("Allowed roles:", allowedRoles);

        if (typeof roleId !== "number") {
            console.error("Unable to determine roleId from /me response");
            showAlert("Invalid session data. Please log in again.", "error");
            clearSession();
            redirectToLogin();
            return false;
        }

        // 4) Enforce access
        if (!allowedRoles.includes(roleId)) {
            console.error(`Role ${roleId} not permitted. Allowed:`, allowedRoles);
            showAlert("You don't have permission to access this page.", "error");
            setTimeout(redirectToDashboard, 2000);
            return false;
        }

        console.log("Role check passed");
        return true;

    } catch (err) {
        console.error("Network error during role check:", err);
        showAlert("Network error. Please try again.", "error");
        /*clearSession();
        redirectToLogin();*/
        return false;
    }
}


// Check if user is authenticated with better token validation
async function checkAuth() {
    console.log("Checking authentication via /me…");

    try {
        const res = await fetch(`${API_BASEE}/me`, {
            method: "GET",
            credentials: "include",
            headers: { "Content-Type": "application/json" }
        });

        if (!res.ok) {
            // either not logged in or token expired/invalid
            console.error("Auth check failed:", res.status);
            showAlert("Please log in to access this page.", "error");
            clearSession();
            redirectToLogin();
            return false;
        }

        // you could optionally use the returned user info here
        // const userInfo = await res.json();
        console.log("Authentication check passed");
        return true;

    } catch (err) {
        console.error("Network error during auth check:", err);
        showAlert("Network error. Please try again.", "error");
        /*clearSession();
        redirectToLogin();*/
        return false;
    }
}

// Helper to clear any client‐stored session data
function clearSession() {
    localStorage.removeItem("userId");
    localStorage.removeItem("email");
    localStorage.removeItem("roleId");
    // (you no longer have a token in localStorage)
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



// Global logout function
function logout() {
    console.log("Logging out...");
    localStorage.removeItem("role");
    localStorage.removeItem("email");
    localStorage.removeItem("userId");
    fetch("https://localhost:8080/logout", {
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
    return fetch(`https://localhost:8080${endpoint}`, {
        ...options,
        credentials: "include" // <— just this
    }).then(response => {
        if (response.status === 401) {
            console.error("Unauthorized access - not logged in");
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
            throw new Error(`HTTPs ${response.status}`);
        }

        return response.json();
    });
}


// Expose utilities globally
window.utils = {
    showAlert,
    redirectTo,
    checkRole,
    checkAuth,
    onClick,
    onSubmit,
    logout,
    makeAuthenticatedRequest,
    redirectToLogin,
    redirectToDashboard
};

// Also expose logout globally for easy access
window.logout = logout;