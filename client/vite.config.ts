import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': '/src',
      components: '/src/components',
      pages: '/src/pages',
      api: '/src/api',
      hooks: '/src/hooks',
      assets: '/src/assets',
      lib: '/src/lib',
      providers: '/src/providers',
      global: '/src/global',
      types: '/src/types',
      messages: '/src/messages',
    },
  },
});
