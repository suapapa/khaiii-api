FROM ubuntu:20.04 AS khaiii-builder

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update && apt install -y \
    build-essential \
    cmake \
    wget \
    libboost-all-dev \
 && apt clean \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /src
RUN wget https://github.com/kakao/khaiii/archive/refs/tags/v0.4.tar.gz && \
    tar -xzf v0.4.tar.gz && \
    rm v0.4.tar.gz

WORKDIR /src/khaiii-0.4
RUN mkdir build
WORKDIR /src/khaiii-0.4/build

RUN mkdir /usr_local
RUN cmake -E env CXXFLAGS="-w" cmake -DCMAKE_INSTALL_PREFIX=/usr_local -DCMAKE_BUILD_TYPE=Release  ..
RUN make -j$(nproc)
RUN make large_resource
# RUN make resource
RUN make install

# ---

FROM golang:1.27 AS go-builder

COPY --from=khaiii-builder /usr_local/ /usr/local/
RUN ldconfig

WORKDIR /app

COPY internal ./internal
COPY pkg ./pkg
COPY main.go .
COPY go.mod .
COPY go.sum .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

RUN go build -o app

RUN apt update && apt install -y locales
RUN locale-gen en_US.UTF-8
RUN localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8

# ---

FROM ubuntu:24.04

RUN apt update && apt install -y locales \
 && apt clean \
 && rm -rf /var/lib/apt/lists/*
RUN locale-gen en_US.UTF-8
RUN localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8

COPY --from=khaiii-builder /usr_local/ /usr/local/
RUN ldconfig

WORKDIR /app
COPY --from=go-builder /app/app .

CMD ["/app/app"]