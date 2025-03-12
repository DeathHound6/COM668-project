from config import api_host, logger
from exceptions import ExternalAPIException
from utility import make_api_request, HTTPMethodEnum, ThreadLocalStorage
from os import getenv
from typing import Any


class BackendAPIClient:
    jwt = None
    tls = ThreadLocalStorage()


    def __init__(self):
        self.jwt = None


    def handle_jwt(self):
        self.jwt = self.tls.get("jwt")
        if self.jwt is None:
            logger.warning("[A.I.M.S] JWT not found in TLS. Getting new JWT")
            self.jwt = self.get_jwt()
            self.tls.set("jwt", self.jwt)


    def get_jwt(self) -> str:
        logger.info("[A.I.M.S] Getting JWT")
        resp = make_api_request(
            url=f"{api_host}/users/login",
            method=HTTPMethodEnum.POST,
            body={
                "email": getenv("API_USER_EMAIL"),
                "password": getenv("API_USER_PW")
            }
        )
        if resp.status_code != 204:
            raise ExternalAPIException(resp.json()["error"])
        jwt = resp.headers["Authorization"]
        return jwt


    def get_providers(self, provider_type: str):
        self.handle_jwt()
        logger.info(f"[A.I.M.S] Getting providers of type '{provider_type}'")
        resp = make_api_request(
            url=f"{api_host}/providers",
            method=HTTPMethodEnum.GET,
            headers={
                "Authorization": self.jwt
            },
            params={
                "provider_type": provider_type
            }
        )
        # If we get a 401, we need to get a new JWT
        if resp.status_code == 401:
            logger.info("[A.I.M.S] JWT expired. Getting new JWT and recalling get_providers")
            self.handle_jwt()
            return self.get_providers(provider_type)
        elif resp.status_code != 200:
            raise ExternalAPIException(resp.json()["error"])
        return resp.json()["data"]


    def get_hosts(self, filters: dict[str, Any] = {}):
        self.handle_jwt()
        logger.info("[A.I.M.S] Getting host machines")
        resp = make_api_request(
            url=f"{api_host}/hosts",
            params=filters,
            method=HTTPMethodEnum.GET,
            headers={
                "Authorization": self.jwt
            }
        )
        if resp.status_code == 401:
            logger.info("[A.I.M.S] JWT expired. Getting new JWT and recalling get_hosts")
            self.handle_jwt()
            return self.get_hosts(filters)
        elif resp.status_code != 200:
            raise ExternalAPIException(resp.json()["error"])
        return resp.json()["data"]


    def get_me(self):
        self.handle_jwt()
        logger.info("[A.I.M.S] Getting user")
        resp = make_api_request(
            url=f"{api_host}/me",
            method=HTTPMethodEnum.GET,
            headers={
                "Authorization": self.jwt
            }
        )
        if resp.status_code == 401:
            logger.info("[A.I.M.S] JWT expired. Getting new JWT and recalling get_me")
            self.handle_jwt()
            return self.get_me()
        elif resp.status_code != 200:
            raise ExternalAPIException(resp.json()["error"])
        return resp.json()["data"]


    def get_teams(self, filters: dict[str, Any]):
        self.handle_jwt()
        logger.info("[A.I.M.S] Getting teams")
        resp = make_api_request(
            url=f"{api_host}/teams",
            method=HTTPMethodEnum.GET,
            headers={
                "Authorization": self.jwt
            },
            params=filters
        )
        if resp.status_code == 401:
            logger.info("[A.I.M.S] JWT expired. Getting new JWT and recalling get_teams")
            self.handle_jwt()
            return self.get_teams(filters)
        elif resp.status_code != 200:
            raise ExternalAPIException(resp.json()["error"])
        return resp.json()["data"]


    def create_incident(self, incident_data: dict[str, Any]) -> str:
        self.handle_jwt()
        logger.info("[A.I.M.S] Creating incident")
        resp = make_api_request(
            url=f"{api_host}/incidents",
            method=HTTPMethodEnum.POST,
            headers={
                "Authorization": self.jwt
            },
            body=incident_data
        )
        if resp.status_code == 401:
            logger.info("[A.I.M.S] JWT expired. Getting new JWT and recalling create_incident")
            self.handle_jwt()
            return self.create_incident(incident_data)
        elif resp.status_code != 201:
            raise ExternalAPIException(resp.json()["error"])
        return resp.headers["Location"]


    def get_incidents(self, params: dict[str, Any]) -> dict[str, Any]:
        self.handle_jwt()
        logger.info("[A.I.M.S] Getting incidents")
        resp = make_api_request(
            url=f"{api_host}/incidents",
            params=params,
            method=HTTPMethodEnum.GET,
            headers={
                "Authorization": self.jwt
            }
        )

        if resp.status_code == 401:
            logger.info("[A.I.M.S] JWT expired. Getting new JWT and recalling get_incidents")
            self.handle_jwt()
            return self.get_incidents(params)
        elif resp.status_code != 200:
            raise ExternalAPIException(resp.json()["error"])
        return resp.json()["data"]


backend_client = BackendAPIClient()
