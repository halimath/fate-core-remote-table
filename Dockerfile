FROM node:14-alpine

RUN apk add m4

WORKDIR /app

COPY package*.json ./
RUN npm i

ARG version
ARG vcs_ref
ARG build_number

RUN node -e "const p = require('./package.json'); p.version = '${version}'; p.versionLabel = '${version} (${build_number}; ${vcs_ref})'; console.log(JSON.stringify(p));" > package.json.mod
RUN mv package.json.mod package.json

COPY public ./public
COPY src ./src
COPY tsconfig.json ./
COPY webpack.config.js ./

RUN npm run dist

RUN m4 --define=VERSION="${version} (${vcs_ref})" dist/index.html > index.html
RUN sed -Ee "s/serviceworker\\./serviceworker.${build_number}./" index.html > /tmp/index.html
RUN sed -Ee "s/diceroller\\./diceroller.${build_number}./" /tmp/index.html > dist/index.html
RUN sed -Ee "s/^const CacheVersion =.*$/const CacheVersion = ${build_number};/" public/serviceworker.js > dist/serviceworker.${build_number}.js
RUN mv dist/css/diceroller.css dist/css/diceroller.${build_number}.css
RUN mv dist/diceroller.js dist/diceroller.${build_number}.js

FROM nginx:alpine

COPY --from=0 /app/dist /usr/share/nginx/html/

