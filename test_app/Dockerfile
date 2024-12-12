FROM node:20.18

# Setup the FS
RUN mkdir /app
COPY bin/ /app/bin
RUN chmod +x -R /app/bin

RUN npm install -g pnpm

WORKDIR /app
ENTRYPOINT [ "/app/bin/entrypoint.sh" ]