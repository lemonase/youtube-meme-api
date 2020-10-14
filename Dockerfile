# container for building
FROM golang:1.13 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o youtube-meme-api .

# container for running
FROM scratch
COPY --from=builder /build/youtube-meme-api /app/bin/

ENV YT_API_KEY my_key_value
ENV PORT 8000

ENTRYPOINT /app/bin/youtube-meme-api

CMD ["--key", "${YT_API_KEY}"]
CMD ["--port", "${PORT}"]

EXPOSE 8000
