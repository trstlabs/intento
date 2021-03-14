const VuetifyLoaderPlugin = require('vuetify-loader/lib/plugin')

module.exports = {
  /*devServer: {
    disableHostCheck: true,
  },*/
  configureWebpack: {
    plugins: [
      new VuetifyLoaderPlugin(),
      
    ],
    /*optimization: {
      runtimeChunk: 'single',
      splitChunks: {
        chunks: 'all',
        maxInitialRequests: Infinity,
        minSize: 0,
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
    },*/
  },
  
};
