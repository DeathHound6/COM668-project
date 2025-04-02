from src.config import api_host, logger
from src.exceptions import ExternalAPIException
from src.utility import HTTPMethodEnum, ThreadLocalStorage
from src.http_clients.base import APIClient
from os import getenv
from typing import Any


class BackendAPIClient(APIClient):
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
        resp = self.make_api_request(
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

    def get_providers(self, provider_type: str) -> list[dict[str, Any]]:
        self.handle_jwt()
        logger.info(f"[A.I.M.S] Getting providers of type '{provider_type}'")
        resp = self.make_api_request(
            url=f"{api_host}/providers",
            method=HTTPMethodEnum.GET,
            headers={
                "Authorization": self.jwt
            },
            params={
                "provider_type": provider_type
            }
        )
        if resp.status_code == 401:
            logger.info("[A.I.M.S] JWT expired. Getting new JWT and recalling get_providers")
            self.handle_jwt()
            return self.get_providers(provider_type)
        elif resp.status_code != 200:
            raise ExternalAPIException(resp.json()["error"])
        return resp.json()["data"]

    def get_hosts(self, filters: dict[str, Any] = {}) -> list[dict[str, Any]]:
        self.handle_jwt()
        logger.info("[A.I.M.S] Getting host machines")
        resp = self.make_api_request(
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

    def get_teams(self, filters: dict[str, Any]) -> list[dict[str, Any]]:
        self.handle_jwt()
        logger.info("[A.I.M.S] Getting teams")
        resp = self.make_api_request(
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
        resp = self.make_api_request(
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

    def get_incidents(self, params: dict[str, Any] = {}) -> list[dict[str, Any]]:
        self.handle_jwt()
        logger.info("[A.I.M.S] Getting incidents")
        resp = self.make_api_request(
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

    def update_incident(self, incident_id: str, body: dict[str, Any]) -> None:
        response = self.make_api_request(
            url=f"{api_host}/incidents/{incident_id}",
            method=HTTPMethodEnum.PUT,
            headers={
                "Authorization": self.jwt
            },
            body=body
        )
        if response.status_code == 401:
            logger.info("[A.I.M.S] JWT expired. Getting new JWT and recalling update_incident")
            self.handle_jwt()
            return self.update_incident(incident_id, body)
        elif response.status_code != 204:
            raise ExternalAPIException(response.json()["error"])

    def post_comment_on_incident(self, incident_id: str, comment: str) -> str:
        response = self.make_api_request(
            url=f"{api_host}/incidents/{incident_id}/comments",
            method=HTTPMethodEnum.POST,
            headers={
                "Authorization": self.jwt
            },
            body={
                "comment": comment
            }
        )
        if response.status_code == 401:
            logger.info("[A.I.M.S] JWT expired. Getting new JWT and recalling post_comment_on_incident")
            self.handle_jwt()
            return self.post_comment_on_incident(incident_id, comment)
        elif response.status_code != 201:
            raise ExternalAPIException(response.json()["error"])
        return response.headers["Location"]

    def delete_comment_on_incident(self, incident_id: str, comment_id: str) -> None:
        response = self.make_api_request(
            url=f"{api_host}/incidents/{incident_id}/comments/{comment_id}",
            method=HTTPMethodEnum.DELETE,
            headers={
                "Authorization": self.jwt
            }
        )
        if response.status_code == 401:
            logger.info("[A.I.M.S] JWT expired. Getting new JWT and recalling delete_comment_on_incident")
            self.handle_jwt()
            return self.delete_comment_on_incident(incident_id, comment_id)
        elif response.status_code != 204:
            raise ExternalAPIException(response.json()["error"])


backend_client = BackendAPIClient()
