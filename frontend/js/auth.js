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

// Helper: hide error message
function hideError() {
    const errorElem = document.getElementById("error-message");
    if (errorElem) {
        errorElem.style.display = "none";
    }
}

// After login/signup, redirect to dashboard
function redirectToDashboard() {
    window.location.href = "pages/dashboard.html";
}

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) {
        return parts.pop().split(';').shift();
    }
    return null;
}

// Handle login form submission
document.addEventListener("DOMContentLoaded", () => {
    const loginForm = document.getElementById("login-form");
    if (loginForm) {
        loginForm.addEventListener("submit", async (e) => {
            e.preventDefault();
            hideError(); // clear previous errors

            const email = document.getElementById("email").value.trim();
            const password = document.getElementById("password").value;

            if (!email || !password) {
                showError("Please fill in all fields");
                return;
            }
            if (password=="nullByGoogle"){
                showError("Please try other passwords");
                return
            }
            if (password=="nullByGithub"){
                showError("Please try other passwords");
                return
            }

            try {
                const res = await fetch(`${API_BASE}/login`, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ email, password }),
                    credentials: "include"
                });


                // successful login: store token
                const token = getCookie("auth_token");
                if (!token) {
                    showError("UnAuthorized, please sing up firstly!");
                    return;
                }
                // Decode JWT to get role
                try {
                    const payload = JSON.parse(atob(token.split(".")[1]));
                    const roleId = parseInt(payload.role_id, 10);
                    const userId = parseInt(payload.sub, 10);

                    localStorage.setItem("token", token);
                    localStorage.setItem("roleId", roleId);
                    localStorage.setItem("email", email);
                    localStorage.setItem("userId", userId);

                    redirectToDashboard();
                } catch (jwtError) {
                    console.error("JWT decode error:", jwtError);
                    showError("Invalid token received");
                }
            } catch (err) {
                console.error("Login error:", err);
                showError("Network error. Please try again.");
            }
        });
    }

    // Handle signup form if present
    const signupForm = document.getElementById("signup-form");
    if (signupForm) {
        signupForm.addEventListener("submit", async (e) => {
            e.preventDefault();
            hideError();

            const email = document.getElementById("signup-email").value.trim();
            const password = document.getElementById("signup-password").value;
            const roleId = parseInt(document.getElementById("signup-role").value, 10);

            if (!email || !password || !roleId) {
                showError("Please fill in all fields");
                return;
            }

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
                showError("Network error. Please try again.");
            }
        });
    }
});

function loginWithGoogle() {

    window.location.href = "http://localhost:8080/auth/google/login";
}

function loginWithGithub() {

    window.location.href = "http://localhost:8080/auth/github/login";
}