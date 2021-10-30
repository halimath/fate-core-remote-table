# fate-core-remote-table

[![CI Status](https://github.com/halimath/fate-core-remote-table/workflows/CI/badge.svg)](https://github.com/halimath/fate-core-remote-table/actions/workflows/ci.yml)
[![CD Status](https://github.com/halimath/fate-core-remote-table/workflows/CD/badge.svg)](https://github.com/halimath/fate-core-remote-table/actions/workflows/cd.yml)

A virtual remote table for playing [Fate Core](https://www.evilhat.com/home/fate-core/) role playing
games.

## About

This repo contains a web application that supports people playing any kind of _Fate Core_ based role playing
game remotely using whatever video conferencing tool they like. This app adds support for
* rolling fate dice - in case you don't have physical ones
* manage fate points
* share _aspects_ with all players

# Architecture

This application is build from two parts:
1. a backend service managing data and distributing data update among players
1. a web frontend consiting of a single page application that can be used with different devices

The parts communicate using [websockets](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API) and
communication utilizes the [CQRS](https://en.wikipedia.org/wiki/Command%E2%80%93query_separation#Command_query_responsibility_segregation)
paradigm. The communication protocol is documented in 
[`docs/fate-core-remote-table.asyncapi.yml`](./docs/fate-core-remote-table.asyncapi.yml) which is an 
[AsyncAPI](https://www.asyncapi.com/) spec.


![System Architecture](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/halimath/fate-core-remote-table/master/docs/sysarch.puml)

The frontend application uses an internal architecture modeled after the well known _model, view, control_
pattern. A _controller_ receives _messages_ that describe updates to the _model_ and executes them, returning
a new _model_ value to be rendered by the _view_ functions. The weccoframework provides excellent support for
implementing this kind of architectures.

The backend applies the 
[_entity-control-boundary-pattern_ (ECB)](http://www.cs.sjsu.edu/~pearce/modules/patterns/enterprise/ecb/ecb.htm).
Currently, all state is only stored _in memory_ but plans are to integrate a database to store the table
state so that tables can be revisited.

The entity model is shown by the below UML class diagramm:

![Entities](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/halimath/fate-core-remote-table/master/docs/entities.puml)

# Development

The backend is implemented using Java 17 and [Quarkus](https://quarkus.io/). The frontend is implemented using
TypeScript and the [wecco framework](https://github.com/weccoframework/core). Almost all CSS is coming from
[Tailwind](https://tailwindcss.com/) with minimal CSS being written to embed the Fate Core font for displaying
dice results.

## Local Environment

For being able to develop the app, you should have a local install of
* Java 17
* Node v14 (v16 should work as well)
* NPM (>=6.14)

You should also have an IDE which supports Java and TypeScript. VSCode works perfectly, IntelliJ IDEA works,
too. I haven't tried other IDEs, but the should work the same.

The application is build using [Apache Maven](https://maven.apache.org/). The frontend part uses webpack,
postcss and typescript and the build is configured using the respective files (i.e. `webpack.config.js`) and
run via `npm run build`. This task is also wrapped in a maven module to control execution from a single build
process. 

To get working locally, you should open two terminal windows (or tabs or whatever you use). In the first,
run 

```
$ ./mvnw -pl service compile quarkus:dev
```

This will start the quarkus development tools and bring up the backend app on `localhost:8080`.

In the second terminal, run

```
$ cd app
$ npm i
$ npm start
```

(You need to install dependencies with `npm i` only once or when the `package.json` file changes). This will
start the webpack dev server to bring up the frontend on `localhost:9999`. 

Now point your browser to [http://localhost:9999](http://localhost:9999) and you can use the app.

## CI/CD

Both parts of the application are wrapped in a single OCI container build with `podman` (but you can use 
Docker as well). The container build uses multiple stages and builds the whole app as part of the container 
build. The final container will only contain the compiled application, though.

We use [Github Actions](https://github.com/features/actions) to build the application, run the tests, build
the container image and publish it to [https://ghcr.io](https://github.com/features/packages).

# License

Copyright 2021 Alexander Metzner.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
