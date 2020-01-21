FROM balenalib/raspberrypi3-golang:latest-build AS build_base

WORKDIR /go/src/github.com/wisegrowth/balena-wisebot

#COPY . .

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download

#######################

FROM build_base AS build
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARM=7 GOARCH=arm go install
RUN go build

#######################

FROM balenalib/raspberrypi3-debian:stretch

COPY --from=build /go/bin/balena-wisebot .

CMD ./balena-wisebot