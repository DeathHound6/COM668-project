FROM golang:1.21

# Setup the FS
COPY bin/ /app/bin
RUN chmod +x -R /app/bin

WORKDIR /app
ENTRYPOINT [ "/app/bin/entrypoint.sh" ]