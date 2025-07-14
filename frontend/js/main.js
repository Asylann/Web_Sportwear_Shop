// frontend/js/main.js

// On page load, check token and redirect if logged in
window.addEventListener("DOMContentLoaded", () => {
    const token = localStorage.getItem("token");
    const role = localStorage.getItem("role");

    if (token) {
        // Redirect based on role
        if (role === "admin") {
            window.location.href = "pages/admin.html";
        } else if (role === "seller") {
            window.location.href = "pages/seller.html";
        } else {
            window.location.href = "pages/dashboard.html";
        }
    }
});

// Global logout function
function logout() {
    localStorage.removeItem("token");
    localStorage.removeItem("role");
    localStorage.removeItem("email");
    localStorage.removeItem("userId");
    fetch('http://localhost:8080/logout', {
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
