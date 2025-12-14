// Configuration for different environments
const config = {
    // For local development
    local: {
        backendURL: 'localhost:8080',
        useWSS: false
    },
    // For production (Render)
    production: {
        backendURL: window.ENV_BACKEND_URL || 'your-backend-url.onrender.com',
        useWSS: true
    }
};

// Auto-detect environment
const isDevelopment = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';
const currentConfig = isDevelopment ? config.local : config.production;

window.APP_CONFIG = currentConfig;
