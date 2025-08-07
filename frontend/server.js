const https = require('https');
const fs = require('fs');
const path = require('path');
const express = require('express');

const app = express();

// Serve static files from current directory
app.use(express.static(__dirname));

// Handle all routes to serve index.html (if you have SPA routing)
app.get('*', (req, res) => {
    res.sendFile(path.join(__dirname, 'index.html'));
});

// HTTPS options
const options = {
    key: fs.readFileSync('../localhost+2-key.pem'),
    cert: fs.readFileSync('../localhost+2.pem')
};

// Create HTTPS server
https.createServer(options, app).listen(8081, () => {
    console.log('Frontend HTTPS server running on https://localhost:8081');
});