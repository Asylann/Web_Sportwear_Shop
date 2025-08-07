// frontend/js/auth.js

// Base URL for your API
const API_BASE = "https://localhost:8080";

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

            const Email = document.getElementById("email").value.trim();
            const password = document.getElementById("password").value;

            if (!Email || !password) {
                showError("Please fill in all fields");
                return;
            }
            if (password === "nullByGoogle" || password === "nullByGithub") {
                showError("Please try other passwords");
                return;
            }

            try {
                // 1) POST to /login; server will set HttpsOnly, Secure cookie
                const res = await fetch(`${API_BASE}/login`, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    credentials: "include",
                    body: JSON.stringify({Email,password})
                });
                if (!res.ok) {
                    const err = await res.json();
                    showError(err.error || "Login failed");
                    return;
                }

                // 2) Immediately GET /me with the cookie, to retrieve user info
                const meRes = await fetch(`${API_BASE}/me`, {
                    method: "GET",
                    credentials: "include"
                });
                if (!meRes.ok) {
                    showError("Failed to fetch user info");
                    return;
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

                // 4) Redirect on success
                redirectToDashboard();

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

    window.location.href = "https://localhost:8080/auth/google/login";
}

function loginWithGithub() {

    window.location.href = "https://localhost:8080/auth/github/login";
}