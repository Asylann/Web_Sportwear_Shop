// dashboard.js

document.addEventListener("DOMContentLoaded", () => {
    // 1) Must be logged in
    const token = localStorage.getItem("token");
    const roleId = parseInt(localStorage.getItem("roleId"), 10);

    if (!token) {
        alert("Please log in first.");
        return window.location.href = "/index.html";
    }

    // 2) Always land here — now show/hide nav links
    document.getElementById("nav-products").style.display = "inline-block";

    if (roleId === 2 || roleId === 3) {
        // sellers and admins
        document.getElementById("nav-seller").style.display = "inline-block";
    } else {
        document.getElementById("nav-seller").style.display = "none";
    }

    if (roleId === 3) {
        // only admins
        document.getElementById("nav-admin").style.display = "inline-block";
    } else {
        document.getElementById("nav-admin").style.display = "none";
    }

    // 3) Logout button
    document.getElementById("logoutBtn").addEventListener("click", () => {
        localStorage.clear();
        window.location.href = "/index.html";
    });
});
