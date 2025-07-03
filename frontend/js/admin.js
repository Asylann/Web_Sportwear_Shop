document.addEventListener("DOMContentLoaded", () => {
    utils.checkAuth();
    utils.checkRole(3);

    console.log("Admin dashboard loaded.");
    // You can now load admin-specific content here.
});
