import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  //base: "http://localhost:4080/public/",
  base: "",
  build: {
    outDir: '../../client/default',
    emptyOutDir: true,

  }
})
