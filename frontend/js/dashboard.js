// dashboard.js

document.addEventListener("DOMContentLoaded", () => {
    // 1) Must be logged in
    const token = localStorage.getItem("token");
    const roleId = parseInt(localStorage.getItem("roleId"), 10);
    const email = localStorage.getItem("email")
    const userId = localStorage.getItem("userId")
    if (!token) {
        alert("Please log in first.");
        return window.location.href = "/index.html";
    }

    // 2) Display user role information
    displayUserInfo(roleId,email,userId);

    // 3) Always show products link
    document.getElementById("nav-products").style.display = "inline-block";

    // 4) Show/hide nav links based on role
    if (roleId === 2 || roleId === 3) {
        // sellers and admins can access seller panel
        document.getElementById("nav-seller").style.display = "inline-block";
    } else {
        document.getElementById("nav-seller").style.display = "none";
    }

    if (roleId === 3) {
        // only admins can access admin panel
        document.getElementById("nav-admin").style.display = "inline-block";
    } else {
        document.getElementById("nav-admin").style.display = "none";
    }

    // 5) Logout button
    const logoutBtn = document.getElementById("logoutBtn");
    if (logoutBtn) {
        logoutBtn.addEventListener("click", () => {
            localStorage.clear();
            window.location.href = "/index.html";
        });
    }
});

function displayUserInfo(roleId,email,userId) {
    const roleName = getRoleName(roleId);
    const userInfoContainer = document.querySelector(".dashboard-card p");

    if (userInfoContainer) {
        userInfoContainer.innerHTML = `
            <strong>Your Email:</strong> ${email}<br>
            <strong>Your Role:</strong> ${roleName}<br>
            <strong>Your ID:</strong> ${userId}<br>
            <strong>Available Actions:</strong><br>
            • View Products (all users)<br>
            ${roleId >= 2 ? '• Manage Products (sellers/admins)<br>' : ''}
            ${roleId === 3 ? '• User Management (admins only)<br>' : ''}
            Choose an action from the navigation above.
        `;
    }
}

function getRoleName(roleId) {
    switch (roleId) {
        case 1: return "Customer";
        case 2: return "Seller";
        case 3: return "Admin";
        default: return "Unknown";
    }
}