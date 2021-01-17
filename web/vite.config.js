import vue from '@vitejs/plugin-vue'

export default {
  plugins: [vue()],
  optimizeDeps: {
    include: [
      'monaco-editor/esm/vs/language/json/json.worker.js',
      'monaco-editor/esm/vs/editor/editor.worker.js'
    ]
  }
}
