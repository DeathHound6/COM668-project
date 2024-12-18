from logging import Formatter, DEBUG, INFO, WARNING, ERROR, CRITICAL
from typing import Optional, Any
from enum import Enum
import requests


class HTTPMethodEnum(Enum):
    GET = "GET"
    PUT = "PUT"
    DELETE = "DELETE"
    POST = "POST"
    PATCH = "PATCH"


def make_api_request(url: str, method: HTTPMethodEnum, params: Optional[dict[str, Any]]=None, headers: Optional[dict[str, str]]=None, body: Optional[dict[str, Any]]=None) -> requests.Response:
    verify_ssl = False
    if headers is not None:
        if "Content-Type" not in list(headers.keys()):
            headers["Content-Type"] = "application/json"
    else:
        headers = {
            "Content-Type": "application/json"
        }

    kwargs = {
        "url": url,
        "params": params,
        "headers": headers,
        "verify": verify_ssl
    }
    # If we have a body, we need to send it as JSON
    if method in (HTTPMethodEnum.PUT, HTTPMethodEnum.POST, HTTPMethodEnum.PATCH) and body is not None:
        if headers["Content-Type"] == "application/json":
            kwargs["json"] = body
        else:
            kwargs["data"] = body

    if method is HTTPMethodEnum.GET:
        resp = requests.get(**kwargs)
    elif method is HTTPMethodEnum.PUT:
        resp = requests.put(**kwargs)
    elif method is HTTPMethodEnum.DELETE:
        resp = requests.delete(**kwargs)
    elif method is HTTPMethodEnum.POST:
        resp = requests.post(**kwargs)
    elif method is HTTPMethodEnum.PATCH:
        resp = requests.patch(**kwargs)

    return resp


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