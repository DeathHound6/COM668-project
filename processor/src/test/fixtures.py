from datetime import datetime, timedelta
from uuid import uuid4
from typing import Callable, Any, Union
from copy import deepcopy
import json


LOG_PROVIDERS = [
    {
        "name": "Sentry",
        "fields": [
            {
                "key": "enabled",
                "value": True
            },
            {
                "key": "orgSlug",
                "value": "test"
            },
            {
                "key": "projSlug",
                "value": "test"
            }
        ]
    }
]


ALERT_PROVIDERS = [
    {
        "name": "Slack",
        "fields": [
            {
                "key": "enabled",
                "value": True
            }
        ]
    }
]


SENTRY_EVENTS = [
    {
        "id": "1",
        "title": "Test event",
        "culprit": "GET /test",
        "tags": [
            {
                "key": "handled",
                "value": "no"
            },
            {
                "key": "server_name",
                "value": "test"
            }
        ],
        "errors": [
            {
                "data": {
                    "url": "test.js"
                }
            }
        ],
        "entries": [
            {
                "type": "exception",
                "data": {
                    "values": [
                        {
                            "type": "Error",
                            "stacktrace": {
                                "frames": [
                                    {
                                        "filename": "test.py",
                                        "absPath": "test.py",
                                        "context": [
                                            [1, "test"],
                                            [2, "test"],
                                            [3, "test"]
                                        ]
                                    }
                                ]
                            }
                        }
                    ]
                }
            }
        ]
    }
]


SENTRY_HEADERS = {
    "Link": "<https://sentry.io/api/0/projects/sentry/sentry/events/?cursor=0:0:0>; rel=\"next\"; results=\"true\"; cursor=\"0:0:0\""  # noqa: E501
}


INCIDENTS = [
    {
        "createdAt": (datetime.now() - timedelta(days=3)).isoformat(),
        "comments": [
            {
                "createdAt": datetime.now().isoformat()
            }
        ],
        "uuid": "e3bfb580-de90-4969-9bdf-e1fed5467fb7",
        "summary": "Test incident",
        "description": "Test incident",
        "hostsAffected": [
            {
                "uuid": "c2d9e7dc-a94e-47da-9f92-e3c541b1d56c"
            }
        ],
        "resolutionTeams": [
            {
                "uuid": "e6435916-fab8-4af0-93a7-762c9f399245"
            }
        ]
    }
]


HOSTS = [
    {
        "uuid": "c2d9e7dc-a94e-47da-9f92-e3c541b1d56c",
        "name": "test",
        "team": {
            "uuid": "e6435916-fab8-4af0-93a7-762c9f399245",
            "users": [
                {
                    "slackID": "test"
                }
            ]
        }
    }
]


TEAMS = [
    {
        "uuid": "e6435916-fab8-4af0-93a7-762c9f399245",
        "name": "test",
        "users": [
            {
                "slackID": "test"
            }
        ]
    }
]


SLACK_CONVERSATION = {
    "ok": True,
    "channel": {
        "id": uuid4().hex
    }
}


class MockHTTPResponse:
    def __init__(self, status_code, headers, body):
        self.status_code = status_code
        self.text = json.dumps(body)
        self.headers = headers

    def json(self):
        return json.loads(self.text)


def mock_api_request(status_code: dict[str, Callable[[dict[str, Any]], int]] = {},
                     headers: dict[str, Callable[[dict[str, Any]], dict[str, Any]]] = {},
                     body: dict[str, Callable[[dict[str, Any]], Union[dict[str, Any], list[dict[str, Any]]]]] = {}):
    def wrapper(**kwargs):
        method = kwargs.get("method", None) or ""
        url = kwargs.get("url", None) or ""
        params = kwargs.get("params", None) or {}
        # Method may be an enum
        if method != "":
            method = method.value

        default_status: dict[str, Callable[[dict[str, Any]], int]] = {
            "POST login": lambda **_: 204,
            "POST incidents": lambda **_: 201,
            "PUT incidents": lambda **_: 204,
            "POST comments": lambda **_: 201,
            "DELETE comments": lambda **_: 204,
        }
        default_headers: dict[str, Callable[[dict[str, Any]], dict[str, Any]]] = {
            "GET events": lambda **_: deepcopy(SENTRY_HEADERS),
            "POST login": lambda **_: {"Authorization": "Bearer test"},
            "POST incidents": lambda **_: {"Location": "http://test.com/incidents/1"},
            "POST comments": lambda **_: {"Location": "http://test.com/incidents/1/comments/1"},
            "POST providers": lambda **_: {"Location": "http://test.com/providers/1"},
            "POST hosts": lambda **_: {"Location": "http://test.com/hosts/1"}
        }
        default_body: dict[str, Callable[[dict[str, Any]], Union[dict[str, Any], list[dict[str, Any]]]]] = {
            "GET events": lambda **_: deepcopy(SENTRY_EVENTS),
            "GET incidents": lambda **_: {"data": deepcopy(INCIDENTS)},
            "GET providers": lambda **_: {"data": deepcopy(LOG_PROVIDERS)
                                          if params.get("provider_type") == "log"
                                          else deepcopy(ALERT_PROVIDERS)},
            "GET hosts": lambda **_: {"data": deepcopy(HOSTS)},
            "GET teams": lambda **_: {"data": deepcopy(TEAMS)},
            "POST incidents": lambda **_: {"data": deepcopy(INCIDENTS)},
            "POST conversations.create": lambda **_: deepcopy(SLACK_CONVERSATION),
            "POST conversations.join": lambda **_: {"ok": True},
            "POST conversations.invite": lambda **_: {"ok": True},
            "POST chat.postMessage": lambda **_: {"ok": True}
        }

        endpoint = None

        if "events" in url:
            endpoint = "events"
        elif "comments" in url:
            endpoint = "comments"
        elif "incidents" in url:
            endpoint = "incidents"
        elif "providers" in url:
            endpoint = "providers"
        elif "hosts" in url:
            endpoint = "hosts"
        elif "teams" in url:
            endpoint = "teams"
        elif "login" in url:
            endpoint = "login"
        elif "conversations.create" in url:
            endpoint = "conversations.create"
        elif "conversations.invite" in url:
            endpoint = "conversations.invite"
        elif "conversations.join" in url:
            endpoint = "conversations.join"
        elif "chat.postMessage" in url:
            endpoint = "chat.postMessage"

        if not endpoint:
            return MockHTTPResponse(404, {}, {})

        method_endpoint = f"{method} {endpoint}"
        status = status_code.get(method_endpoint, default_status.get(method_endpoint, lambda **_: 200))
        header = headers.get(method_endpoint, default_headers.get(method_endpoint, lambda **_: {}))
        json = body.get(method_endpoint, default_body.get(method_endpoint, lambda **_: {}))
        return MockHTTPResponse(status(**kwargs), header(**kwargs), json(**kwargs))
    return wrapper
