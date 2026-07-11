import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { createRequire } from 'module'
const require = createRequire(import.meta.url)
const Database = require('better-sqlite3')

function recalculatePlugin() {
  return {
    name: 'recalculate-plugin',
    configureServer(server) {
      server.middlewares.use('/api/entries/recalculate', async (req, res) => {
        if (req.method !== 'POST') { res.statusCode = 405; res.end(); return }
        try {
          const url = new URL(req.url, 'http://localhost')
          const type = url.searchParams.get('type')
          if (!type) { res.statusCode = 400; res.end(JSON.stringify({code:1,message:'type required'})); return }
          const db = new Database('./server/data/time.db')
          const rows = db.prepare('SELECT id, recorded_at FROM entries WHERE type = ? ORDER BY recorded_at ASC').all(type)
          if (rows.length > 0) {
            const update = db.prepare('UPDATE entries SET description = ? WHERE id = ?')
            for (let i = rows.length - 1; i >= 0; i--) {
              let count = 1
              for (let j = i - 1; j >= 0; j--) {
                const d1 = new Date(rows[j+1].recorded_at.substring(0,10))
                const d2 = new Date(rows[j].recorded_at.substring(0,10))
                if ((d1 - d2) / 86400000 === 1) count++; else break
              }
              update.run(String(count), rows[i].id)
            }
          }
          db.close()
          res.statusCode = 200; res.setHeader('Content-Type','application/json'); res.end(JSON.stringify({code:0,message:'ok'}))
        } catch(e) { res.statusCode = 500; res.end(JSON.stringify({code:1,message:e.message})) }
      })
    }
  }
}

export default defineConfig({
  plugins: [recalculatePlugin(), vue()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})