FROM resin/raspberry-pi-alpine-golang:slim as build

ENV GOOS    linux
ENV GOARCH  arm

# Copy required header files
WORKDIR /usr/local/include
COPY include/*.h ./

# Copy static library
WORKDIR /usr/local/lib
COPY lib/libws2811.a ./

# Copy go library
WORKDIR /go/src/ws2811
COPY ./ws2811 ./

# Specify directory for go source
WORKDIR /go/src/github.com/deisterhold/BuildLED

# Copy over the source code
COPY . ./

# Build the executable
RUN go build

FROM resin/resin/raspberry-pi-alpine-golang:slim

ENV INITSYSTEM on

# Copy the executable over
COPY --from=build /go/src/github.com/deisterhold/BuildLED/BuildLED ./

# Run the executable
CMD ./BuildLED
