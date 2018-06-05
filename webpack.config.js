var webpack = require("webpack");
var glob = require("glob");
var CopyWebpackPlugin = require("copy-webpack-plugin");
var ExtractTextPlugin = require("extract-text-webpack-plugin");
var ManifestPlugin = require("webpack-manifest-plugin");
var PROD = process.env.NODE_ENV || "development";
var CleanWebpackPlugin = require("clean-webpack-plugin");

var entries = {
  application: [
    'babel-polyfill',
    './assets/css/application.scss',
  ],
}

glob.sync("./assets/**/*.*").reduce((_, entry) => {
  let key = entry.replace(/(\.\/assets\/(js|css|go)\/)|\.(jsx?|s[ac]ss|go)/g, '')
  if(key.startsWith("_") || (/(jsx?|s[ac]ss|go)$/i).test(entry) == false) {
    return
  }
  
  if( entries[key] == null) {
    entries[key] = [entry]
    return
  } 
  
  entries[key].push(entry)
})

var mode = 'development';
if (PROD != "development") {
    mode = 'production';
}

module.exports = {
  mode: mode,
  entry: entries,
  output: {
    filename: "[name].[hash].js",
    path: `${__dirname}/public/assets`
  },
  plugins: [
    new CleanWebpackPlugin([
      "public/assets"
    ], {
      verbose: false,
    }),
    new ExtractTextPlugin("[name].[hash].css"),
    new CopyWebpackPlugin(
      [{
        from: "./assets",
        to: ""
      }], {
        copyUnmodified: true,
        ignore: ["css/**", "js/**"]
      }
    ),
    new webpack.LoaderOptionsPlugin({
      minimize: true,
      debug: false
    }),
    new ManifestPlugin({
      fileName: "manifest.json"
    })
  ],
  module: {
    rules: [{
      test: /\.jsx?$/,
      exclude: /node_modules/,
      use: [
          {
              loader: 'babel-loader',
              options: {
                  presets: ['react']
              }
          }
      ],
    },
      {
        test: /\.s[ac]ss$/,
        use: ExtractTextPlugin.extract({
          fallback: "style-loader",
          use: [{
            loader: "css-loader",
            options: {
              sourceMap: true
            }
          },
            {
              loader: "sass-loader",
              options: {
                sourceMap: true
              }
            }
          ]
        })
      },
      { test: /\.(woff|woff2|ttf|svg)(\?v=\d+\.\d+\.\d+)?$/,use: "url-loader"},
      { test: /\.eot(\?v=\d+\.\d+\.\d+)?$/,use: "file-loader" },
      {
        test: /\.go$/,
        use: "gopherjs-loader"
      }
    ]
  }
};

if (PROD != "development") {
  module.exports.plugins.push(
    new webpack.optimize.UglifyJsPlugin({
      beautify: false,
      mangle: {
        screw_ie8: true,
        keep_fnames: true
      },
      compress: {
        screw_ie8: true
      },
      comments: false
    })
  );
}
