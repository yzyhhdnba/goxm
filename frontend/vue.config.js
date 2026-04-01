module.exports = {
  assetsDir: 'static',
  parallel: false,
  publicPath: './',
  devServer: {
    proxy: {
      '/api/v1': {
        target: process.env.VUE_APP_API_PROXY_TARGET || 'http://127.0.0.1:8080',
        changeOrigin: true,
      }
    }
  }
}
