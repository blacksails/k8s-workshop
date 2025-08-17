FROM --platform=$BUILDPLATFORM golang:alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o ./bin/httplogger

FROM scratch

COPY --from=builder /app/bin/httplogger /usr/local/bin/httplogger

CMD ["httplogger"]
