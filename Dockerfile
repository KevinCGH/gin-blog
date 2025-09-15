FROM ubuntu:latest
LABEL authors="kevinchen"

ENTRYPOINT ["top", "-b"]