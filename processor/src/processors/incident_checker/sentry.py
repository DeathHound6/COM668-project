from config import logger
from http_clients.sentry import get_issues
from http_clients.backend import backend_client
from http_clients.slack import slack_client
from typing import Any
from exceptions import ExternalAPIException
from utility import generate_hash
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
    file = None

    if len(event["errors"]) > 0:
        file = event["errors"][0]["data"]["url"]

    # NOTE: future TODO: event["entries"][x]["data"]["type"] == "request" (request values), "breadcrumbs" (logs) for extra context
    stack_trace = ""
    for entry in event["entries"]:
        if entry["type"] == "exception":
            for value in entry["data"]["values"]:
                if value["type"] == "Error":
                    for frame in value["stacktrace"]["frames"]:
                        if file is not None and frame["filename"] == file or frame["absPath"] == file:
                            # NOTE: future TODO
                            # frame["lineNo"] is the line number of the error
                            # frame["colNo"] is the column number of the error
                            for ctx in frame["context"]:
                                # ctx[0] is the line number
                                stack_trace += ctx[1]
    incident_hash = generate_hash(stack_trace)

    try:
        incidents = backend_client.get_incidents({"hash": incident_hash})
        if not incidents:
            raise ExternalAPIException("Incident not found")
        logger.info(f"[SENTRY] Incident already exists for event: {event['id']}")
        return
    except ExternalAPIException as e:
        pass

    # NOTE: Only supporting JavaScript for now
    logger.info(f"[SENTRY] Determining root cause for event: {event['id']}")
    root_cause = ""
    if file is not None and "node_modules" in file:
        # NOTE: If error is from node_modules, it is most likely a dependency issue
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

        if len(modules) > 0:
            root_cause = f"Dependency issue with {modules[0]['name']}@{modules[0].get('version', 'latest')}"
        else:
            root_cause = "Unknown dependency issue"
    elif file is not None:
        # NOTE: If error is not from node_modules, it is most likely a code issue
        root_cause = f"Endpoint {endpoint}"
        pass
    else:
        # NOTE: If file not found
        logger.warning(f"[SENTRY] Could not find file path where error was raised for event: {event['id']}")
        root_cause = "Unknown issue"
        pass

    # Get affected servers list
    logger.info(f"[SENTRY] Determining affected servers for event: {event['id']}")
    server_names = [tag["value"] for tag in event["tags"] if tag["key"] == "server_name"]
    if not server_names:
        logger.warning(f"[SENTRY] Could not find hostnames for event: {event['id']}")
        return

    hosts = []
    try:
        hosts = backend_client.get_hosts({"hostnames": ",".join(server_names)})
    except ExternalAPIException as e:
        logger.exception(e)
    if not hosts:
        logger.warning(f"[SENTRY] Could not find host for event: {event['id']}")
        return

    teams = []
    try:
        teams = backend_client.get_teams({"pageSize": 1000})
    except ExternalAPIException as e:
        logger.exception(e)
    if not teams:
        logger.warning(f"[SENTRY] Could not find teams")
        return

    logger.info(f"[SENTRY] Determining resolution teams for event: {event['id']}")
    resolution_teams = []
    # TODO: AI/ML used to determine resolution teams based on existing team names and keywords in error message
    if "node_modules" not in str(file):
        resolution_teams += [host["team"]["uuid"] for host in hosts]
    if "net" in str(message).lower():
        resolution_teams += [team["uuid"] for team in teams if team["name"] == "NetOps"]

    if not resolution_teams:
        resolution_teams = [host["team"]["uuid"] for host in hosts]

    incident_body = {
        "summary": message[:100],
        "description": f"{message}\n{endpoint}\n{root_cause}"[:500],
        "hostsAffected": [host["uuid"] for host in hosts],
        "resolutionTeams": resolution_teams,
        "hash": incident_hash
    }

    incident_url = None
    try:
        incident_url = backend_client.create_incident(incident_body)
    except ExternalAPIException as e:
        logger.exception(e)
    if not incident_url:
        logger.warning(f"[SENTRY] Could not create incident for event: {event['id']}")
        return

    users = set()
    users.update([u["slackID"] for host in hosts for u in host["team"]["users"] if u["slackID"]])
    users.update([u["slackID"] for team in teams for u in team["users"] if u["slackID"]])

    if len(users) > 0:
        # TODO: make the url host configurable - this url points to the frontend
        incident_uuid = incident_url.split("/")[-1]
        incident_url = f"http://localhost:3000/incidents/{incident_uuid}"
        channel_name = f"incident-{incident_uuid}"
        for provider in alert_providers:
            try:
                if any([field["value"] for field in provider["fields"] if str(field["name"]).lower() == "enabled"]) and \
                   str(provider["name"]).lower() == "slack":
                    logger.info(f"[SENTRY] Sending incident to Slack channel: {channel_name}")
                    channel = slack_client.create_conversation(channel_name)
                    slack_client.join_conversation(channel["channel"]["id"])
                    slack_client.invite_to_conversation(channel["channel"]["id"], users)
                    slack_client.send_message(channel["channel"]["id"], f"New incident: {incident_url}")
            except ExternalAPIException as e:
                logger.exception(e)
