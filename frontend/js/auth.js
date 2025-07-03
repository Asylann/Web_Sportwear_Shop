// frontend/js/auth.js

// Base URL for your API
const API_BASE = "http://localhost:8080";

// Helper: show error message under the form
function showError(message) {
    const errorElem = document.getElementById("error-message");
    if (errorElem) {
        errorElem.textContent = message;
        errorElem.style.display = "block";
    } else {
        alert(message);
    }
}

// After login/signup, redirect based on role
function redirectByRole(roleId) {
    window.location.href = "/pages/dashboard.html";
}

// Handle login form submission
document.addEventListener("DOMContentLoaded", () => {
    const loginForm = document.getElementById("login-form");
    if (loginForm) {
        loginForm.addEventListener("submit", async (e) => {
            e.preventDefault();
            showError(""); // clear

            const email = document.getElementById("email").value.trim();
            const password = document.getElementById("password").value;

            try {
                const res = await fetch(`${API_BASE}/login`, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ email, password }),
                });
                const result = await res.json();

                if (!res.ok) {
                    showError(result.error || "Login failed");
                    return;
                }

                // successful login: store token + roleId
                const token = result.data.token;
                const payload = JSON.parse(atob(token.split(".")[1]));
                const roleId = parseInt(payload.role_id, 10);

                localStorage.setItem("token", token);
                localStorage.setItem("roleId", roleId);

                redirectByRole(roleId);
            } catch (err) {
                console.error("Login error:", err);
                showError("Something went wrong.");
            }
        });
    }

    // Handle signup form if present
    const signupForm = document.getElementById("signup-form");
    if (signupForm) {
        signupForm.addEventListener("submit", async (e) => {
            e.preventDefault();
            showError("");

            const email = document.getElementById("signup-email").value.trim();
            const password = document.getElementById("signup-password").value;
            const roleId = parseInt(document.getElementById("signup-role").value, 10);

            try {
                const res = await fetch(`${API_BASE}/signup`, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ email, password, roleId }),
                });
                const result = await res.json();

                if (!res.ok) {
                    showError(result.error || "Signup failed");
                    return;
                }

                alert("Signup successful! Please log in.");
                window.location.href = "/index.html";
            } catch (err) {
                console.error("Signup error:", err);
                showError("Something went wrong.");
            }
        });
    }
});
