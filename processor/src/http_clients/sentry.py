from src.exceptions import ExternalAPIException
from src.utility import HTTPMethodEnum
from src.http_clients.base import APIClient
from src.config import sentry_token, logger
from typing import Optional, Any


class SentryAPIClient(APIClient):
    def parse_link_header(self, link_header: str) -> list[dict[str, str]]:
        parsed_links = list()
        # Example value
        # <https://sentry.io/api/0/projects/sentry/sentry/events/?cursor=0:0:0>; rel="previous"; results="true"; cursor="0:0:0",
        for link in link_header.split(", "):
            parts = link.split("; ")
            rel = None
            url = None
            results = None
            cursor = None
            for part in parts:
                if part.startswith("rel="):
                    rel = part.split("=")[1].replace('"', "")
                elif part.startswith("results="):
                    results = part.split("=")[1].replace('"', "")
                elif part.startswith("cursor="):
                    cursor = part.split("=")[1].replace('"', "")
                elif part.startswith("<") and part.endswith(">"):
                    url = part.replace("<", "").replace(">", "")

            if rel is None or url is None or results is None or cursor is None:
                logger.warning(f"[SENTRY] Could not parse link header part: '{part}'")
                continue
            parsed_links.append({
                "url": url,
                "rel": rel,
                "results": results.lower() == "true",
                "cursor": {
                    "identifier": int(cursor.split(":")[0]),  # is usually 0
                    "offset": int(cursor.split(":")[1]),
                    "isPrevious": bool(int(cursor.split(":")[2]))
                }
            })
        return parsed_links

    def get_issues(
        self, settings_fields: list[dict[str, Any]], offset: Optional[int] = None
    ) -> tuple[list[dict[str, Any]], list[dict[str, Any]]]:
        logger.info("[SENTRY] Fetching issues")
        org_slug_field = [field for field in settings_fields if field["key"] == "orgSlug"]
        proj_slug_field = [field for field in settings_fields if field["key"] == "projSlug"]
        if not org_slug_field or not proj_slug_field:
            logger.warning("[SENTRY] Fields 'orgSlug' or 'projSlug' are not present")
            return list()
        org_slug = org_slug_field[0]["value"]
        proj_slug = proj_slug_field[0]["value"]
        # https://docs.sentry.io/api/pagination/
        if offset is None:
            offset = 0
        response = self.make_api_request(
            url=f"https://sentry.io/api/0/projects/{org_slug}/{proj_slug}/events/",
            method=HTTPMethodEnum.GET,
            headers={
                "Authorization": f"Bearer {sentry_token}"
            },
            params={
                "full": True,
                "cursor": f"0:{offset}:0"
            }
        )
        if response.status_code == 401:
            raise ExternalAPIException("[SENTRY] Invalid Bearer token")
        elif response.status_code == 403:
            raise ExternalAPIException("[SENTRY] Insufficient Bearer token scopes")
        elif response.status_code != 200:
            raise ExternalAPIException(f"[SENTRY] Failed to fetch issues. Status '{response.status_code}'")
        link_headers = self.parse_link_header(response.headers["Link"])
        return response.json(), link_headers


sentry_client = SentryAPIClient()
