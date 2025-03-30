#!/bin/bash

cd /app
pip3 install -r requirements.txt
pip3 install -r requirements-dev.txt
pip3 install -e .
python3 /app/src/main.py