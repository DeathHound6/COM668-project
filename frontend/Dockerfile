FROM node:20.18

# Setup the FS
COPY bin/ /app/bin
RUN chmod +rx -R /app/bin
COPY ./package.json /app
COPY ./pnpm-lock.yaml /app
COPY ./eslint.config.mjs /app
COPY ./next-env.d.ts /app
COPY ./tailwind.config.ts /app
COPY ./postcss.config.mjs /app
COPY ./next.config.ts /app

RUN npm install -g pnpm

WORKDIR /app
ENTRYPOINT [ "/app/bin/entrypoint.sh" ]