FROM golang:1.13

RUN mkdir /app
ADD . /app
WORKDIR /app

ENV YT_API_KEY my_key_value
ENV PORT 8000

RUN go build -o bin/youtube-meme-api

# CMD ["/app/bin/youtube-meme-api"]
ENTRYPOINT /app/bin/youtube-meme-api

CMD ["--key", "${YT_API_KEY}"]
CMD ["--port", "${PORT}"]

EXPOSE 8000
