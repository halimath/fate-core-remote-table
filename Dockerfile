FROM node:16 as NODEJS

WORKDIR /build/app

COPY ./app/package.json ./
COPY ./app/package-lock.json ./

RUN npm i

COPY ./app/ ./
COPY ./docs/ ../docs

RUN npm run build

FROM golang:1.18-rc-bullseye as GOLANG

WORKDIR /backend

COPY ./backend/ ./

COPY --from=NODEJS /build/app/dist/ /backend/internal/boundary/public/

ENV CGO_ENABLED=0
RUN go build .

FROM scratch

COPY --from=GOLANG /backend/backend /fate-core-remote-table

EXPOSE 8080

ENTRYPOINT [ "/fate-core-remote-table" ]