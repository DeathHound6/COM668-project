from src.utility import HTTPMethodEnum
from typing import Optional, Any
import requests


class APIClient:
    def make_api_request(self, url: str, method: HTTPMethodEnum, params: Optional[dict[str, Any]] = None,
                         headers: Optional[dict[str, str]] = None,
                         body: Optional[dict[str, Any]] = None) -> requests.Response:
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
