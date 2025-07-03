// frontend/js/components.js

// Reusable UI components

const components = {
    // Create a styled button
    button: (text, id, extraClass = "") => {
        const btn = document.createElement("button");
        btn.id = id;
        btn.className = `btn ${extraClass}`;
        btn.textContent = text;
        return btn;
    },

    // Create an alert message
    alert: (message, type = "info") => {
        const div = document.createElement("div");
        div.className = `alert alert-${type}`;
        div.textContent = message;
        setTimeout(() => div.remove(), 3000);
        return div;
    },

    // Create a product card
    productCard: (product) => {
        const card = document.createElement("div");
        card.className = "product-card";

        card.innerHTML = `
      <img src="${product.imageURL}" alt="${product.name}" class="product-image">
      <h3>${product.name}</h3>
      <p>${product.description}</p>
      <span class="price">${product.price} ₸</span>
    `;

        return card;
    },

    // Create an input field
    inputField: (type, id, placeholder = "", value = "") => {
        const input = document.createElement("input");
        input.type = type;
        input.id = id;
        input.placeholder = placeholder;
        input.value = value;
        return input;
    },
};

window.components = components;
