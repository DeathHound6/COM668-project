FROM denoland/deno:ubuntu

# Install Node.Js (for NPM)
RUN apt update
RUN apt install -y ca-certificates curl gnupg
RUN mkdir -p /etc/apt/keyrings
RUN curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
RUN echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_20.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list
RUN apt update
RUN apt install -y nodejs
# Verify npm install
RUN npm --version

# Setup the FS
COPY bin/ /app/bin
RUN chmod +rx -R /app/bin
COPY ./deno.json /app
COPY ./package.json /app
COPY ./angular.json /app
COPY ./pnpm-lock.yaml /app

WORKDIR /app
ENTRYPOINT [ "/app/bin/entrypoint.sh" ]