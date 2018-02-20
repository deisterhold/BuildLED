FROM resin/raspberry-pi-golang:slim as build

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

# Download any 3rd party libraries
RUN go get

# Build the executable
RUN go build

#FROM resin/raspberry-pi-debian:latest

ENV INITSYSTEM on

# Copy the executable over
#WORKDIR /go/bin
#COPY --from=build /go/src/github.com/deisterhold/BuildLED/BuildLED ./

# Run the executable
#CMD /go/bin/BuildLED
CMD /go/src/github.com/deisterhold/BuildLED/BuildLED
