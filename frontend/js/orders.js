// orders.js - Handle orders page functionality

const API_BASE = "https://localhost:8081/api";

let orders = [];
let expandedOrders = new Set(); // Track which orders are expanded

document.addEventListener("DOMContentLoaded", () => {
    if (!utils.checkAuth()) return;
    if (!utils.checkRole(1, 2, 3)) return;

    loadOrders();
});

function loadOrders() {
    showLoading();

    fetch(`${API_BASE}/orders`, {
        method: 'GET',
        credentials: "include"
    })
        .then(res => {
            if (!res.ok) {
                throw new Error(`HTTP ${res.status}`);
            }
            return res.json();
        })
        .then(data => {
            if (data.data && Array.isArray(data.data)) {
                orders = data.data;
                displayOrders();
            } else if (data.data && data.data.length === 0) {
                showEmptyOrders();
            } else {
                showEmptyOrders();
            }
        })
        .catch(err => {
            console.error('Error loading orders:', err);
            showError(err.message);
        })
        .finally(() => {
            hideLoading();
        });
}

function showLoading() {
    document.getElementById('loading').style.display = 'block';
    document.getElementById('orders-content').style.display = 'none';
    document.getElementById('empty-orders').style.display = 'none';
    document.getElementById('error-message').style.display = 'none';
}

function hideLoading() {
    document.getElementById('loading').style.display = 'none';
}

function displayOrders() {
    if (!orders || orders.length === 0) {
        showEmptyOrders();
        return;
    }

    document.getElementById('orders-content').style.display = 'block';
    document.getElementById('empty-orders').style.display = 'none';
    document.getElementById('error-message').style.display = 'none';

    // Update orders count
    const ordersCount = orders.length;
    document.getElementById('orders-count').textContent = `${ordersCount} order${ordersCount !== 1 ? 's' : ''}`;

    // Sort orders by date (newest first)
    orders.sort((a, b) => new Date(b.create_at) - new Date(a.create_at));

    // Display orders
    const ordersList = document.getElementById('orders-list');
    ordersList.innerHTML = '';

    orders.forEach(order => {
        const orderElement = createOrderElement(order);
        ordersList.appendChild(orderElement);
    });
}

function createOrderElement(order) {
    const orderCard = document.createElement('div');
    orderCard.className = 'order-card';
    orderCard.setAttribute('data-order-id', order.id);

    const isExpanded = expandedOrders.has(order.id);
    if (isExpanded) {
        orderCard.classList.add('expanded');
    }

    // Format dates
    const createDate = formatDate(order.create_at);
    const deliveryDate = formatDate(order.delivered_at);

    // Determine order status based on delivery date
    const now = new Date();
    const deliveryDateTime = new Date(order.delivered_at);
    const isDelivered = deliveryDateTime <= now;

    orderCard.innerHTML = `
        <div class="order-header" onclick="toggleOrderDetails(${order.id})">
            <div class="order-header-content">
                <div class="order-info-left">
                    <div class="order-id">Order #${order.id}</div>
                    <div class="order-date">Placed: ${createDate}</div>
                </div>
                <div class="order-info-right">
                    <div class="order-status">
                        <span class="status-badge ${isDelivered ? 'delivered' : 'processing'}">
                            ${isDelivered ? 'Delivered' : 'Processing'}
                        </span>
                        <span class="transport-badge ${order.transport_type.toLowerCase()}">
                            ${order.transport_type}
                        </span>
                    </div>
                    <div class="expand-toggle">
                        <span>${isExpanded ? 'Hide' : 'Show'} Details</span>
                        <span class="expand-icon">${isExpanded ? '▲' : '▼'}</span>
                    </div>
                </div>
            </div>
        </div>
        <div class="order-details" style="display: ${isExpanded ? 'block' : 'none'};">
            <div class="order-meta">
                <div class="meta-item">
                    <span class="meta-label">Delivery Address</span>
                    <span class="meta-value">${order.address}</span>
                </div>
                <div class="meta-item">
                    <span class="meta-label">Expected Delivery</span>
                    <span class="meta-value">${deliveryDate}</span>
                </div>
                <div class="meta-item">
                    <span class="meta-label">Transport Type</span>
                    <span class="meta-value">${order.transport_type}</span>
                </div>
                <div class="meta-item">
                    <span class="meta-label">Cart ID</span>
                    <span class="meta-value">#${order.cart_id}</span>
                </div>
            </div>
            <div class="order-items">
                <div class="order-items-header">Order Items</div>
                <div id="order-items-${order.id}" class="order-items-list">
                    ${isExpanded ? '' : '<div class="loading-items">Loading items...</div>'}
                </div>
            </div>
        </div>
    `;

    // Load items if already expanded
    if (isExpanded) {
        loadOrderItems(order.id);
    }

    return orderCard;
}

function toggleOrderDetails(orderId) {
    const orderCard = document.querySelector(`[data-order-id="${orderId}"]`);
    const orderDetails = orderCard.querySelector('.order-details');
    const expandIcon = orderCard.querySelector('.expand-icon');
    const expandText = orderCard.querySelector('.expand-toggle span:first-child');

    const isCurrentlyExpanded = expandedOrders.has(orderId);

    if (isCurrentlyExpanded) {
        // Collapse
        orderDetails.style.display = 'none';
        expandIcon.textContent = '▼';
        expandText.textContent = 'Show Details';
        orderCard.classList.remove('expanded');
        expandedOrders.delete(orderId);
    } else {
        // Expand
        orderDetails.style.display = 'block';
        expandIcon.textContent = '▲';
        expandText.textContent = 'Hide Details';
        orderCard.classList.add('expanded');
        expandedOrders.add(orderId);

        // Load order items
        loadOrderItems(orderId);
    }
}

function loadOrderItems(orderId) {
    const itemsContainer = document.getElementById(`order-items-${orderId}`);
    if (!itemsContainer) return;

    // Show loading state
    itemsContainer.innerHTML = '<div class="loading-items">Loading items...</div>';

    fetch(`${API_BASE}/orders/${orderId}`, {
        method: 'GET',
        credentials: "include"
    })
        .then(res => {
            if (!res.ok) {
                throw new Error(`HTTP ${res.status}`);
            }
            return res.json();
        })
        .then(data => {
            if (data.data && Array.isArray(data.data)) {
                displayOrderItems(orderId, data.data);
            } else {
                itemsContainer.innerHTML = '<div class="loading-items">No items found</div>';
            }
        })
        .catch(err => {
            console.error('Error loading order items:', err);
            itemsContainer.innerHTML = '<div class="loading-items">Error loading items</div>';
        });
}

function displayOrderItems(orderId, items) {
    const itemsContainer = document.getElementById(`order-items-${orderId}`);
    if (!itemsContainer) return;

    if (!items || items.length === 0) {
        itemsContainer.innerHTML = '<div class="loading-items">No items in this order</div>';
        return;
    }

    itemsContainer.innerHTML = '';

    items.forEach(item => {
        const itemElement = createOrderItemElement(item);
        itemsContainer.appendChild(itemElement);
    });
}

function createOrderItemElement(item) {
    const orderItem = document.createElement('div');
    orderItem.className = 'order-item';

    orderItem.innerHTML = `
        ${item.imageURL ?
        `<img src="${item.imageURL}" alt="${item.name}" class="item-image" onclick="viewProduct(${item.id})">` :
        `<div class="placeholder-image">No Image</div>`
    }
        <div class="item-details">
            <h3 class="item-name" onclick="viewProduct(${item.id})">${item.name}</h3>
            <p class="item-description">${item.description}</p>
            <div class="item-meta">
                <span>Size: ${item.size}</span>
                <span>Category ID: ${item.category_id}</span>
            </div>
        </div>
        <div class="item-price">${item.price} ₸</div>
    `;

    return orderItem;
}

function viewProduct(productId) {
    window.location.href = `product-detail.html?id=${productId}`;
}

function formatDate(dateString) {
    try {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    } catch (err) {
        return dateString; // Return original if parsing fails
    }
}

function showEmptyOrders() {
    document.getElementById('orders-content').style.display = 'none';
    document.getElementById('empty-orders').style.display = 'block';
    document.getElementById('error-message').style.display = 'none';
}

function showError(errorMessage) {
    document.getElementById('orders-content').style.display = 'none';
    document.getElementById('empty-orders').style.display = 'none';
    document.getElementById('error-message').style.display = 'block';

    const errorText = document.getElementById('error-text');
    if (errorText) {
        errorText.textContent = errorMessage || 'Unable to load your orders. Please try again.';
    }
}

// Make functions available globally
window.loadOrders = loadOrders;
window.toggleOrderDetails = toggleOrderDetails;
window.viewProduct = viewProduct;
