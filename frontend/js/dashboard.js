// dashboard.js

const API_BASE = "https://localhost:8080";

document.addEventListener("DOMContentLoaded", () => {
    initDashboard();
});

async function initDashboard() {
    try {
        // 1) Validate session & get user info
        const meRes = await fetch(`${API_BASE}/me`, {
            method: "GET",
            credentials: "include"
        });
        if (!meRes.ok) {
            // not logged in or session expired
            console.error("Failed to fetch /me:", meRes.status);
            alert("Please log in first.");
            clearSession();
            return window.location.href = "/index.html";
        }

        const raw = await meRes.json();
        console.log("Raw `/me` response:", raw);

        if (raw.err) {
            console.error("Backend /me error:", raw.err);
            alert("Please log in first.");
            clearSession();
            return window.location.href = "/index.html";
        }

        const { id, email, roleId } = raw.data;
        console.log({ id, email, roleId });

        // 3) Store what you need in localStorage
        localStorage.setItem("userId", id);
        localStorage.setItem("email", email);
        localStorage.setItem("roleId", roleId);

        // 3) Display user info
        displayUserInfo(roleId, email);

        // 4) Show universal nav items
        document.getElementById("nav-products").style.display = "inline-block";
        document.getElementById("nav-orders").style.display = "inline-block";
        document.getElementById("nav-cart").style.display = "inline-block";

        // 5) Role-based nav
        document.getElementById("nav-seller").style.display =
            (roleId === 2 || roleId === 3) ? "inline-block" : "none";

        document.getElementById("nav-admin").style.display =
            (roleId === 3) ? "inline-block" : "none";

        // 6) Wire up logout
        const logoutBtn = document.getElementById("logoutBtn");
        if (logoutBtn) {
            logoutBtn.addEventListener("click", () => {
                clearSession();
                window.location.href = "/index.html";
            });
        }

    } catch (err) {
        console.error("Network error in initDashboard:", err);
        alert("Network error. Please try again.");
        window.location.href = "/index.html";
    }
}

function displayUserInfo(roleId, email) {
    const roleName = getRoleName(roleId);
    const container = document.querySelector(".dashboard-card p");
    if (!container) return;

    container.innerHTML = `
    <strong>Your Email:</strong> ${email}<br>
    <strong>Your Role:</strong> ${roleName}<br>
    <strong>Available Actions:</strong><br>
    • View Products (all users)<br>
    ${roleId >= 2 ? '• Manage Products (sellers/admins)<br>' : ''}
    ${roleId === 3 ? '• User Management (admins only)<br>' : ''}
    Choose an action from the navigation above.
  `;
}

function getRoleName(roleId) {
    switch (roleId) {
        case 1: return "Customer";
        case 2: return "Seller";
        case 3: return "Admin";
        default: return "Unknown";
    }
}

function clearSession() {
    // clear any client-side stored data
    localStorage.removeItem("userId");
    localStorage.removeItem("email");
    localStorage.removeItem("roleId");
    // auth_token cookie expires server-side or on next login
}
