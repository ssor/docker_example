FROM alpine
COPY hello /
EXPOSE 8001 8002 8003
ENTRYPOINT ["/hello"]
