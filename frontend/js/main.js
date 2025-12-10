// frontend/js/main.js

// On page load, check token and redirect if logged in
window.addEventListener("DOMContentLoaded", () => {
    const role = localStorage.getItem("role");

    if (role === "admin") {
        window.location.href = "pages/admin.html";
    } else if (role === "seller") {
        window.location.href = "pages/seller.html";
    } else {
        window.location.href = "pages/dashboard.html";
    }
});

// Global logout function
function logout() {
    localStorage.removeItem("role");
    localStorage.removeItem("email");
    localStorage.removeItem("userId");
    fetch('https://localhost:8081/api/logout', {
        credentials: "include"
    })
        .then(res => res.json())
        .catch(err => {
            console.error('Error logging out', err);
            utils.showAlert('Error logging out', 'error');
        });
    window.location.href = "/index.html";
}

// Expose logout to global scope
window.logout = logout;
