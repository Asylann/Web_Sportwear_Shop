// dashboard.js

document.addEventListener("DOMContentLoaded", () => {

    const token = getCookie("auth_token")
    if (!token) {
        showError("No token received from server");
        return;
    }


    const payload = JSON.parse(atob(token.split(".")[1]));
    const CookieroleId = parseInt(payload.role_id, 10);
    const CookieuserId = parseInt(payload.sub, 10);
    const Cookieemail = payload.email

    localStorage.setItem("roleId", CookieroleId);
    localStorage.setItem("email", Cookieemail);
    localStorage.setItem("userId", CookieuserId);

    // 1) Must be logged in
    const roleId = parseInt(localStorage.getItem("roleId"), 10);
    const email = localStorage.getItem("email")
    if (!token) {
        alert("Please log in first.");
        return window.location.href = "/index.html";
    }

    localStorage.setItem("token", token)

    // 2) Display user role information
    displayUserInfo(roleId,email);

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

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) {
        return parts.pop().split(';').shift();
    }
    return null;
}

function displayUserInfo(roleId,email,userId) {
    const roleName = getRoleName(roleId);
    const userInfoContainer = document.querySelector(".dashboard-card p");

    if (userInfoContainer) {
        userInfoContainer.innerHTML = `
            <strong>Your Email:</strong> ${email}<br>
            <strong>Your Role:</strong> ${roleName}<br>
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