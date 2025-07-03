document.addEventListener("DOMContentLoaded", () => {
    // Check authentication and role
    if (!utils.checkAuth()) return;
    if (!utils.checkRole(3)) return; // Only admins (role 3)

    console.log("Admin dashboard loaded successfully.");

    // Initialize admin dashboard
    initializeAdminDashboard();
});

function initializeAdminDashboard() {
    const roleId = parseInt(localStorage.getItem("roleId"), 10);

    if (roleId !== 3) {
        utils.showAlert("Access denied. Admin privileges required.", "error");
        setTimeout(() => {
            window.location.href = "dashboard.html";
        }, 2000);
        return;
    }

    // Display admin info
    displayAdminInfo();

    // Load initial data if functions exist
    if (typeof loadUsers === 'function') {
        loadUsers();
    }

    if (typeof loadProducts === 'function') {
        loadProducts();
    }

    if (typeof loadCategories === 'function') {
        loadCategories();
    }

    console.log("Admin dashboard initialized successfully.");
}

function displayAdminInfo() {
    const userStatusElement = document.getElementById("user-status");
    if (userStatusElement) {
        userStatusElement.textContent = "👤 Admin user authenticated";
    }

    // You can add more admin-specific UI updates here
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

// Helper function to handle admin-specific errors
function handleAdminError(error, operation) {
    console.error(`Admin operation failed (${operation}):`, error);
    utils.showAlert(`${operation} failed: ${error.message}`, "error");
}

// Helper function to confirm admin actions
function confirmAdminAction(action, callback) {
    const confirmed = confirm(`Are you sure you want to ${action}? This action cannot be undone.`);
    if (confirmed && typeof callback === 'function') {
        callback();
    }
    return confirmed;
}