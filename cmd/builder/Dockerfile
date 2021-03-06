ARG IMAGE=archlinux
FROM golang:1.15 as builder

WORKDIR /build

# Force modules
ENV GO111MODULE=on

# Cache dependencies
COPY go.* ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o builder /build/cmd/builder

FROM $IMAGE
USER root

RUN useradd -m -G wheel runner || true

WORKDIR /home/runner

RUN mkdir -p /etc/sudoers.d && chmod 750 /etc/sudoers.d && echo 'runner ALL=(ALL) NOPASSWD: ALL' > /etc/sudoers.d/yeet

RUN curl "https://www.archlinux.org/mirrorlist/?country=NL&protocol=http&protocol=https&ip_version=4&use_mirror_status=on" | sed "s/#S/S/g" > /etc/pacman.d/mirrorlist

COPY --from=builder /build/builder /usr/bin/builder

RUN pacman -Syu --noconfirm base-devel git

USER runner

CMD ["/usr/bin/builder"]
