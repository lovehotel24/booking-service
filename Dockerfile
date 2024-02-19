FROM golang:1.21.6 as build_boomking-service
ENV CGO_ENABLED 0
ARG BUILD_REF

COPY . /booking-service

WORKDIR /booking-service

RUN go build -ldflags "-extldflags \"-static\" -X main.build=${BUILD_REF}"

FROM scratch
ARG BUILD_DATE
ARG BUILD_REF
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build_booking-service /booking-service/booking-service /service/booking-service

WORKDIR /service
CMD ["./booking-service"]
EXPOSE 8080

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="booking-service" \
      org.opencontainers.image.authors="Dther <dtherhtun.cw@gmail.com>" \
      org.opencontainers.image.source="https://github.com/lovehotel24/booking-service" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="Love hotel"