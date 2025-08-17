FROM --platform=$BUILDPLATFORM hugomods/hugo:std-git AS static

ARG BUILD_DRAFTS=false

WORKDIR /app

COPY . .

RUN hugo build --source web --buildDrafts=${BUILD_DRAFTS}

#####

FROM --platform=$BUILDPLATFORM golang:alpine AS binary

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

COPY . .
COPY --from=static /app/web/public web/public

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o ./bin/k8s-workshop ./cmd/k8s-workshop

#####

FROM scratch

COPY --from=binary /app/bin/k8s-workshop /usr/local/bin/k8s-workshop

CMD ["k8s-workshop"]
