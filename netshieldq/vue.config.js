const { defineConfig } = require('@vue/cli-service')
// const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');
const webpack = require('webpack');

module.exports = defineConfig({
  publicPath: process.env.NODE_ENV === 'production' ? '../dist/' : '/',
  transpileDependencies: true,

  configureWebpack: {
    plugins: [
      new webpack.ProvidePlugin({
        process: 'process/browser'
      }),
    // new BundleAnalyzerPlugin()

    ],
    resolve: {
      fallback: {
        "os": require.resolve("os-browserify/browser"),
        "crypto": require.resolve("crypto-browserify"),
        "stream": require.resolve('stream-browserify'),
        "fs": false,
        "process": require.resolve("process/browser"),
        "path": require.resolve("path-browserify")
      }
    },
    experiments: {
      asyncWebAssembly: true,
    },
    module: {
      rules: [
        {
          test: /\.wasm$/,
          type: 'webassembly/async',
        },
      ],
    },
  },

  devServer: {
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
        pathRewrite: {
          '^/api': ''
        }
      }
    }
  }
});
