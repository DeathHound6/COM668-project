from config import logger
from http_clients.sentry import get_issues
from typing import Any


def handle_sentry(log_provider: dict[str, Any], alert_providers: list[dict[str, Any]]):
    handled_all_pages = False
    offset = 0
    while not handled_all_pages:
        issue_events, link_headers = get_issues(log_provider["fields"], offset)
        for event in issue_events:
            pass

        for link in link_headers:
            if link["rel"] == "next":
                offset = link["cursor"]["offset"]
                handled_all_pages = link["results"]
                break