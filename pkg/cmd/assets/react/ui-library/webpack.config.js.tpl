const path = require('path');

module.exports = {
  entry: './src/index.tsx', // Changed to .tsx
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'bundle.js',
    libraryTarget: 'umd',
    library: 'MyUILibrary',
    umdNamedDefine: true,
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js'], // Added .tsx and .ts
  },
  module: {
    rules: [
      {
        test: /\.(js|jsx|ts|tsx)$/, // Added ts and tsx
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: ['@babel/preset-env', '@babel/preset-react', '@babel/preset-typescript'], // Added typescript preset
          },
        },
      },
    ],
  },
  externals: {
    react: {
      commonjs: 'react',
      commonjs2: 'react',
      amd: 'react',
      root: 'React',
    },
    'react-dom': {
      commonjs: 'react-dom',
      commonjs2: 'react-dom',
      amd: 'react-dom',
      root: 'ReactDOM',
    },
  },
};
