const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = (env) => {
  const example = env.example; // e.g., "basic-usage", "form-example", etc.

  return {
    mode: 'development',
    entry: `./src/examples/${example}/main.tsx`,
    output: {
      path: path.resolve(__dirname, 'dist', example),
      filename: 'bundle.js',
      clean: true, // Cleans the output folder before building
    },
    resolve: {
      extensions: ['.ts', '.tsx', '.js', '.jsx'],
    },
    module: {
      rules: [
        {
          test: /\.(ts|tsx)$/,
          exclude: /node_modules/,
          use: {
            loader: 'babel-loader',
            options: {
              presets: [
                '@babel/preset-env',
                '@babel/preset-react',
                '@babel/preset-typescript',
              ],
            },
          },
        },
      ],
    },
    plugins: [
      new HtmlWebpackPlugin({
        template: `./src/examples/index.html`,
        filename: 'index.html',
      }),
    ],
    devServer: {
      static: path.resolve(__dirname, 'dist', example),
      port: 3000 + (['basic', 'form-example', 'styled-card-example'].indexOf(example) + 1), // Unique ports: 3001, 3002, 3003
      open: true,
    },
  };
};
