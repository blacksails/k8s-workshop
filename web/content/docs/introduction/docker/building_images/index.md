---
title: Building Images
draft: false
---

# Building Images

## Dockerfile

A `Dockerfile` is a special file which acts as a recipe for how to build an
image. The `Dockerfile` has some specific keywords which have semantic meaning
for how docker should build the container image. We will briefly take a look at
some common ones.

### `FROM`

Is used to tell docker where to start this is also reffered to as the base of
the docker image. Usually this points to an existing image on Dockerhub or
another contrainer registry.

### `WORKDIR`

Is used to specify what should be the working directory for the following
commands. Think of it as using the `cd` command to go to a directory before
executing commands.

### `COPY`

Is used to copy files from your filesystem to the image filesystem. This can be
used to copy source code or compiled binaries into the image filesystem,
depending on wether you build your code outside or within your `Dockerfile`.

### `RUN`

Is used to run commands as part of building your image. The commands used in
`RUN` statements must be available within the docker image. Usually `RUN` is
used to install and or compile software in your docker image.

### `CMD`

Is used to specify what command should be run when your container is starting.

### Other `Dockerfile` keywords

The `Dockerfile` specification is long and has way more keywords than the
above. For a full reference see [the official `Dockerfile`
reference](https://docs.docker.com/reference/dockerfile/).

### Examples

The following are a few simple examples of what `Dockerfile`'s look like.

#### Java Application

You build a java application and now you have a `.jar` file which you would
like to run in a container. In order to run a `.jar` file we need a java
runtime. A popular one is [Eclipse Temurin](https://adoptium.net/temurin),
luckily they have an official repository on Dockerhub which we can use. The
resulting `Dockerfile` would look something like this:

```Dockerfile
FROM eclipse-temurin:21
WORKDIR /opt/app
COPY app.jar .
CMD ["java", "-jar", "/opt/app/app.jar"]
```

#### Serving static website

You have a static website consisting of HTML and javascript which you would
like to have served by a webserver. There is a lot of different web server
implementations, but a popular one is [nginx](https://nginx.org/). In this
example we see how a static website is copied into the nginx image.

```Dockerfile
FROM nginx
COPY static-html-directory /usr/share/nginx/html
```

## Exercises

### Build a custom nginx image

In the [Docker introduction](../) we ran the nginx image. This time let's build
our own nginx image containing a custom HTML file.

- Find the exercise files in the git repo under
  `exercises/introduction/docker/building_images/nginx`.
- Open the Dockerfile and use the `COPY` keyword to copy `index.html` into
  `/usr/share/nginx/html`.
- Build the docker image using `docker build`. Use the `-t` flag to give your
  image a name. For more info about the command run `docker build --help`.
{{% hint info %}}
It is common to release new versions of an image and to be able to refer to
different versions of the same image. This is done by including a tag in the
image name eg. `docker build -t my-nginx:v1.0.0`. If you don't specify a tag,
the tag `latest` is implied.
{{% /hint %}}
- Run the docker image and open the page in the browser.

You should see the following

<img src="hello_world.png" width="250px"/>

### Run the nginx container with a volume mount

We actually don't need to build our own nginx image in order to have the nginx
application show our HTML. Instead we can use docker volume mounts to mount
something from our host machine into the docker container.

To mount a host volume to the container you will need to use the `-v` flag when
you start your container.

Instead of building a custom image, run the nginx image with
`exercises/introduction/docker/building_images/nginx` mounted at
`/usr/share/nginx/html`. When you open the browser you should see the same hello
world message as above.

Documentation for the nginx image:
[https://hub.docker.com/_/nginx](https://hub.docker.com/_/nginx)

### Build a go application

Go is a very simple programming language. I have prepared a simple hello world
webserver which we will compile and pack in a docker image.

- Find the exercise files in the git repo at
  `exercises/introduction/docker/building_images/go`.
- Open the Dockerfile
  - Use the `COPY` command to copy the source files into the image
  - Use the `RUN` command to compile the source files. The go command for
    compiling a binary is `go build -o <output path> <path to source files>`
  - Use the `CMD` command to run the compiled code.
- Build and run docker image.
- Open the page in the browser.

You should see this:

<img src="./go_hello_world.png" width="250px"/>

### Build a go application as a multi stage build

Dockerfiles can be multi stage This is useful when you want to exclude source
code, compilers etc. from your final image. If a `Dockerfile` contains multiple
`FROM` statements, it is considered multi stage. The lines after the last
`FROM` statement defines the final image.

Prior stages can be named and then referenced in later stages. This is done by
using the `AS` keyword in a `FROM` statement.

```Dockerfile
FROM golang:alpine AS builder
```

Here we named the stage "builder". We can then reference this stage with the
`COPY` statement in order to copy something from one stage to another.

```Dockerfile
COPY --from=builder /app/my-binary /usr/local/bin/my-binary
```

- Reuse the exercise files from the last exercise.
- Modify the Dockerfile to become a multistage build.
  - Give the existing `FROM` statement a name using the `AS` keyword.
  - Add a new `FROM scratch` statement at the bottom of the file.
  - Add a new `COPY` statement which copies the binary from the first stage
    into the final stage.
- Build the image with a new `-t` tag.
- Use the `docker images` command and notice that the image size is now smaller.

{{% hint info %}}
`scratch` is a special empty base image. It doesn't contain any OS, package
manager or anything else. There are two great advantages of using a stripped
image like this.

- Image size. The image will be way smaller as it only holds what is needed for
  your app to run.
- Security. When an image contains a lot of software other than your
  application, the attack surface of the container is larger. You can have CVE's
  in your image which aren't even related to your application code, but they
  might still be there because some of the other software in the image has a CVE.
  By using a stripped down image, you will only need to patch your application
  and not all of the other software in the image.

Sometimes it can be a bit of a hassle to make your application work when based
on the scratch image. Google made some stripped down base images which might
help in these cases, they are called
[distroless](https://github.com/GoogleContainerTools/distroless)
{{% /hint %}}

### Build a go application as a multi architecture image

{{% hint warning %}}
Multi architecture images is an advanced feature of docker and not strictly
nessesary to know. Only do this exercise if you want to have a deeper learning.
{{% /hint %}}

Docker supports building images for multiple architectures. This means that you
can reference a docker image and the one for your computer architecture will
automatically be chosen. In order to build these images we need to do a few
special things in our `Dockerfile`. I won't go into much detail here, instead go
have a look at the example in the docker documentation: [Cross-compiling a Go
application](https://docs.docker.com/build/building/multi-platform/#cross-compiling-a-go-application)

Modify the `Dockerfile` from the previous exercises and build a multi
architecture image.
