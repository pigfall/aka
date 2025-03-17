{
  "name": "{{.Name}}",
  "version": "1.0.0",
  "main": "dist/bundle.js",
  "types": "dist/index.d.ts",
  "scripts": {
    "build": "webpack --mode production && tsc --emitDeclarationOnly",
    "build-types": "tsc --emitDeclarationOnly",
    "example-basic": "webpack server --config=./webpack.example.config.js --mode development --env  example=basic"
  },
  "keywords": [],
  "author": "",
  "license": "ISC",
  "description": "",
  "peerDependencies": {
    "react":"^19.0.0",
    "react-dom":"^19.0.0"
  },
  "devDependencies": {
    "typescript": "^5.8.2",
    "@types/react":"^19.0.0",
    "@types/react-dom":"^19.0.0",
    "webpack": "^5.98.0",
    "webpack-cli": "^6.0.1",
    "webpack-dev-server": "^5.2.0",
    "html-webpack-plugin": "^5.6.3",
    "@babel/core": "^7.15.0",
    "@babel/preset-env": "^7.15.0",
    "@babel/preset-react": "^7.14.5",
    "@babel/preset-typescript": "^7.26.0",
    "babel-loader": "^8.2.2"
  }
}
