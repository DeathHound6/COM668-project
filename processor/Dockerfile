FROM python:3.9

# Setup the FS
RUN mkdir /app
COPY requirements.txt /app
COPY requirements-dev.txt /app
COPY bin/ /app/bin
COPY setup.py /app
RUN chmod +x -R /app/bin

WORKDIR /app
ENTRYPOINT [ "/app/bin/entrypoint.sh" ]