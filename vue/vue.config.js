const VuetifyLoaderPlugin = require('vuetify-loader/lib/plugin')

module.exports = {
  /*devServer: {
    disableHostCheck: true,
  },*/
  configureWebpack: {
    plugins: [
      new VuetifyLoaderPlugin(),

    ],
    optimization: {
      removeAvailableModules: true,
      splitChunks: {
        chunks: 'async',
        minSize: 20000,
      
        maxSize: 0,
        minChunks: 1,
        maxAsyncRequests: 30,
        maxInitialRequests: 30,
        enforceSizeThreshold: 50000,
        cacheGroups: {
          vendor: {
            test: /[\\/]node_modules[\\/]/,
            name(module) {
              const packageName = module.context.match(/[\\/]node_modules[\\/](.*?)([\\/]|$)/)[1];
              return `npm.${packageName.replace('@', '')}`;
            },
            
          },
        },
      },
    },
  },
};
