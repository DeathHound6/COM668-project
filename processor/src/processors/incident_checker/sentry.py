from config import logger
from http_clients.sentry import get_issues
from typing import Any
import re

# Sentry API reqs
# https://sentry.io/api/0/projects/testing-77/test_app/events/?full=1 -> body[index]["contexts"]["trace"]["trace_id"]
# https://de.sentry.io/api/0/organizations/testing-77/events-trace/[trace_id]/?limit=10000&timestamp=[nowMS]&useSpans=1

def handle_sentry(log_provider: dict[str, Any], alert_providers: list[dict[str, Any]]):
    handled_all_pages = False
    offset = 0
    while not handled_all_pages:
        issue_events, link_headers = get_issues(log_provider["fields"], offset)
        for event in issue_events:
            # If it is an unhanded error event
            if any([tag["key"] == "handled" and tag["value"] == "no" for tag in event["tags"]]):
                handle_event(event, alert_providers)

        for link in link_headers:
            if link["rel"] == "next":
                offset = link["cursor"]["offset"]
                handled_all_pages = link["results"]
                break


def handle_event(event: dict[str, Any], alert_providers: list[dict[str, Any]]):
    message = event["title"]
    endpoint = event["culprit"]
    date = event["dateCreated"]
    file = None

    if len(event["errors"]) > 0:
        file = event["errors"][0]["data"]["url"]

    root_cause = ""
    if file is not None and "node_modules" in file:
        # NOTE: Only matching UNIX paths for now
        modules = []
        # If using PNPM
        if "node_modules/.pnpm/" in file:
            matches = re.findall(r"(?P<name>\@?[^@]+)\@(?P<version>[.0-9]+)", file, re.IGNORECASE)
            if matches:
                modules.append({"name": matches[0][0], "version": matches[0][1]})
        else:
            # If using NPM
            matches = re.findall(r"node_modules/(?P<name>[^/]+)/", file, re.IGNORECASE)
            # There shoould only be 1 match for NPM
            if matches:
                modules.append({"name": matches[0], "version": None})
    elif file is not None:
        pass
    else:
        pass