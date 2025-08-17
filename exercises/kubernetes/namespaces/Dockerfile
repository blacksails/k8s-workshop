FROM --platform=$BUILDPLATFORM golang:alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o ./bin/echosvc

FROM scratch

COPY --from=builder /app/bin/echosvc /usr/local/bin/echosvc

ENTRYPOINT ["echosvc"]
