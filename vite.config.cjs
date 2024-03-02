const { resolve } = require('path');
const { defineConfig } = require('vite');
const legacy = require('@vitejs/plugin-legacy');

module.exports = defineConfig({
  plugins: [
    legacy({
      targets: ['defaults', 'not IE 11'],
    }),
  ],
  build: {
    rollupOptions: {
      input: {
        about: resolve(__dirname, '/about.html'),
        main: resolve(__dirname, 'index.html'),

        // aboutjs: resolve(__dirname, 'about.js')
      },
    },
  },
});
