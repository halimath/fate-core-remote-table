FROM node:20 as NODEJS

WORKDIR /build/app

COPY ./app/package.json ./
COPY ./app/package-lock.json ./

RUN npm i

COPY ./app/ ./
COPY ./docs/ ../docs

RUN npm run build

FROM golang:1.22-alpine as GOLANG

ARG version=0.11.0
ARG commit=local

WORKDIR /backend

COPY ./backend/ ./

COPY --from=NODEJS /build/app/dist/ /backend/internal/boundary/public/

ENV CGO_ENABLED=0
RUN go build -ldflags "-X main.Version=${version} -X main.Commit=${commit}" .

FROM scratch

ARG version=0.10.1
ARG commit=local

LABEL maintainer="Alexander Metzner <alexander.metzner@gmail.com>" \
    version=${version} \
    commit=${commit} \
    url="https://github.com/halimath/fate-core-remote-table" \
    vcs-uri="https://github.com/halimath/fate-core-remote-table.git"

COPY --from=GOLANG /backend/backend /fate-core-remote-table

EXPOSE 8080

ENTRYPOINT [ "/fate-core-remote-table" ]