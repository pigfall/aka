{
  "name": "{{.Name}}",
  "version": "1.0.0",
  "main": "dist/bundle.js",
  "scripts": {
    "dev": "webpack server --mode development"
  },
  "keywords": [],
  "author": "",
  "license": "ISC",
  "description": "",
  "dependencies": {
    "react":"^19.0.0",
    "react-dom":"^19.0.0"
  },
  "devDependencies": {
    "@babel/core": "^7.15.0",
    "@babel/preset-env": "^7.15.0",
    "@babel/preset-react": "^7.14.5",
    "@babel/preset-typescript": "^7.26.0",
    "babel-loader": "^8.2.2",
    "html-webpack-plugin": "^5.6.3",
    "webpack":"^5.98.0",
    "webpack-cli":"^6.0.1",
    "webpack-dev-server":"^5.2.0"
  }
}
