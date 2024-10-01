#!/bin/bash

python3 -m pip install -r /app/requirements.txt
python3 -m pip install -r /app/requirements-dev.txt

python3 /app/src/main.py