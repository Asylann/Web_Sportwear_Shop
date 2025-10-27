const https = require('https');
const fs = require('fs');
const path = require('path');
const express = require('express');

const app = express();

app.use((req, res, next) => {
    // ставим CSP только для запросов, где ожидается HTML (чтобы не ломать загрузку шрифтов/картинок)
    if (req.accepts('html')) {
        res.setHeader("Content-Security-Policy",
            "default-src 'self' 'unsafe-inline';"+
                    "script-src 'self' 'unsafe-inline';"+
                    "img-src 'self' https://cdn.jsdelivr.net;"+
                    "connect-src 'self' https://localhost:8080;"+
                    "style-src 'self' https://fonts.googleapis.com 'unsafe-inline';"+
                    "font-src 'self' https://fonts.googleapis.com https://fonts.gstatic.com 'unsafe-inline'"
        );
    }
    next();
});

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
https.createServer(options, app).listen(8082, () => {
    console.log('Frontend HTTPS server running on https://localhost:8082');
});