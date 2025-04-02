import hashlib
from logging import Formatter, DEBUG, INFO, WARNING, ERROR, CRITICAL
from typing import Any
from enum import Enum


class HTTPMethodEnum(Enum):
    GET = "GET"
    PUT = "PUT"
    DELETE = "DELETE"
    POST = "POST"
    PATCH = "PATCH"


class LoggerFormatter(Formatter):
    blue = "\x1b[1;34;40m"
    green = "\x1b[1;32;40m"
    yellow = "\x1b[33;20m"
    red = "\x1b[31;20m"
    bold_red = "\x1b[31;1m"
    reset = "\x1b[0m"
    format_template = "%(message)s"
    FORMATS = {
        DEBUG: f"{blue}{format_template}{reset}",
        INFO: f"{green}{format_template}{reset}",
        WARNING: f"{yellow}{format_template}{reset}",
        ERROR: f"{red}{format_template}{reset}",
        CRITICAL: f"{bold_red}{format_template}{reset}",
    }

    def format(self, record):
        log_fmt = self.FORMATS.get(record.levelno)
        formatter = Formatter(log_fmt)
        return formatter.format(record)


class ThreadLocalStorage:
    def __init__(self):
        self._storage = {}

    def set(self, key: str, value: Any):
        self._storage[key] = value

    def get(self, key: str) -> Any:
        return self._storage.get(key)


def generate_hash(data: str) -> str:
    return hashlib.sha1(data.encode()).hexdigest()
