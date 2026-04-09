const path = require("path");
const webpack = require("webpack");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const { VueLoaderPlugin } = require("vue-loader");

module.exports = (_, argv) => {
  const isProduction = argv.mode === "production";

  return {
    entry: path.resolve(__dirname, "src/main.js"),
    output: {
      path: path.resolve(__dirname, "dist"),
      filename: isProduction ? "js/[name].[contenthash:8].js" : "js/[name].js",
      chunkFilename: isProduction ? "js/[name].[contenthash:8].chunk.js" : "js/[name].chunk.js",
      clean: true,
      publicPath: "/"
    },
    resolve: {
      extensions: [".js", ".vue", ".json"],
      alias: {
        "@": path.resolve(__dirname, "src"),
        "vue$": "vue/dist/vue.esm-bundler.js"
      }
    },
    devtool: isProduction ? "source-map" : "eval-cheap-module-source-map",
    module: {
      rules: [
        {
          test: /\.vue$/,
          loader: "vue-loader"
        },
        {
          test: /\.js$/,
          exclude: /node_modules/,
          use: {
            loader: "babel-loader"
          }
        },
        {
          test: /\.css$/,
          use: ["style-loader", "css-loader"]
        },
        {
          test: /\.(png|jpe?g|gif|svg|webp)$/i,
          type: "asset",
          parser: {
            dataUrlCondition: {
              maxSize: 8 * 1024
            }
          },
          generator: {
            filename: "images/[name].[hash:8][ext]"
          }
        }
      ]
    },
    plugins: [
      new webpack.DefinePlugin({
        __VUE_OPTIONS_API__: JSON.stringify(true),
        __VUE_PROD_DEVTOOLS__: JSON.stringify(false),
        __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: JSON.stringify(false)
      }),
      new VueLoaderPlugin(),
      new HtmlWebpackPlugin({
        template: path.resolve(__dirname, "public/index.html"),
        inject: "body"
      })
    ],
    devServer: {
      host: "0.0.0.0",
      port: 9000,
      hot: true,
      historyApiFallback: true,
      proxy: [
        {
          context: ["/api"],
          target: "http://127.0.0.1:8080",
          changeOrigin: true
        }
      ]
    }
  };
};
