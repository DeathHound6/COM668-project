from config import api_host, logger
from exceptions import ExternalAPIException
from utility import make_api_request, HTTPMethodEnum
from os import getenv


def get_jwt() -> str:
    logger.info("[A.I.M.S] Getting JWT")
    url = f"{api_host}/users/login"
    resp = make_api_request(
        url=url,
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


def get_providers(provider_type: str, jwt: str):
    logger.info(f"[A.I.M.S] Getting providers of type '{provider_type}'")
    url = f"{api_host}/providers"
    params = {
        "provider_type": provider_type
    }
    resp = make_api_request(
        url=url,
        method=HTTPMethodEnum.GET,
        headers={
            "Authorization": jwt
        },
        params=params
    )
    # If we get a 401, we need to get a new JWT
    if resp.status_code == 401:
        logger.info("[A.I.M.S] JWT expired. Getting new JWT and recalling get_providers")
        jwt = get_jwt()
        return get_providers(provider_type, jwt)
    elif resp.status_code != 200:
        raise ExternalAPIException(resp.json()["error"])
    return resp.json()["providers"]
