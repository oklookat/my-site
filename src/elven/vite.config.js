import { defineConfig, loadEnv } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import path from 'path'
import fs from 'fs'

const baseConf = {
  plugins: [svelte()],
  optimizeDeps: { exclude: ["svelte-router-spa"] },
  resolve: {
    alias: [
      {
        find: '@', replacement: path.resolve('./src')
      },
    ],
  }
}

const withHTTPS = (certPath, keyPath) => {
  return {
    server: {
      open: true,
      https: {
        cert: fs.readFileSync(certPath, 'utf8'),
        key: fs.readFileSync(keyPath, 'utf8'),
      }
    }
  }
}

export default defineConfig(({ command, mode }) => {
  // build
  if (command !== 'serve') {
    return baseConf
  }
  // serve
  // env trick: https://github.com/vitejs/vite/issues/1930#issuecomment-986088066
  Object.assign(process.env, loadEnv(mode, process.cwd(), ''))
  const certPath = process.env.VITE_CERT_PATH
  const keyPath = process.env.VITE_KEY_PATH
  const https = withHTTPS(certPath, keyPath)
  return Object.assign(baseConf, https)
})