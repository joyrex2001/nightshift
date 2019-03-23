module.exports = {
    devServer: {
      host: '0.0.0.0',
      port: 5000,
      https: false,
      hotOnly: false,
      proxy: 'http://localhost:8080',
    },
    publicPath: '/public/'
}
