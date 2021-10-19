const TsconfigPathsPlugin = require("tsconfig-paths-webpack-plugin");
const MiniCssExtractPlugin = require('mini-css-extract-plugin')

module.exports = {
    mode: "development",
    entry: "./src/script/index.ts",
    output: {
        filename: "fate-diceroller.js"
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
        static: "./public",
        compress: true,
        host: "0.0.0.0",
        port: 9999,
        open: false,
    }
};