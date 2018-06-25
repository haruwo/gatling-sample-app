FROM golang:1.10-alpine3.7
ADD . /go/src/server
RUN go install server

FROM alpine:3.7
COPY --from=0 /go/bin/server ./server
RUN chmod +x ./server
ENV PORT 8000
EXPOSE 8000
CMD ./server $PORT
