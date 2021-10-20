const TsconfigPathsPlugin = require("tsconfig-paths-webpack-plugin")
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const { resolve } = require("path")

module.exports = {
    mode: "development",
    entry: "./src/main/ts/index.ts",
    output: {
        filename: "fate-diceroller.js",
        path: resolve(__dirname, 'target/classes/META-INF/resources'),
    },
    module: {
        rules: [
            {
                test: /\.ts$/,
                exclude: /node_modules/,
                use: "ts-loader"
            },
            {
                test: /\.css$/i,
                use: [
                    MiniCssExtractPlugin.loader,
                    // "style-loader",
                    "css-loader",
                    "postcss-loader",
                ]
            }
        ]
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: 'css/fate-diceroller.css'
        }),
    ],
    resolve: {
        extensions: [".ts", ".js", ".html"],
        plugins: [new TsconfigPathsPlugin()]
    },
    devServer: {
        static: "./src/main/ressources/META-INF/resources",
        compress: true,
        host: "0.0.0.0",
        port: 9999,
        open: false,
    }
};