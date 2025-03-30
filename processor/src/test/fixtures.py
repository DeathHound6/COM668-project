from datetime import datetime, timedelta
from uuid import uuid4
from src.config import logger


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


SENTRY_HEADERS = [
    {
        "rel": "next",
        "cursor": {
            "offset": 1
        },
        "results": False
    }
]


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
    "channel": {
        "id": uuid4()
    }
}


class MockHTTPResponse:
    def __init__(self, status_code, headers, json):
        self.status_code = status_code
        self.data = json
        self.headers = headers

    def json(self):
        return self.data


def mock_api_request(status_code={}, headers={}, body={}):
    def wrapper(**kwargs):
        method = kwargs.get("method", None) or ""
        url = kwargs.get("url", None) or ""
        params = kwargs.get("params", None) or {}
        # Method may be an enum
        if method != "":
            method = method.value

        default_status = {
            "POST login": 204,
            "POST incidents": 201,
            "PUT incidents": 204,
            "POST comments": 201,
            "DELETE comments": 204,
        }
        default_headers = {
            "GET events": SENTRY_HEADERS.copy(),
            "POST login": {"Authorization": "Bearer test"}
        }
        default_body = {
            "GET events": SENTRY_EVENTS.copy(),
            "GET incidents": {"data": INCIDENTS.copy()},
            "GET providers": {"data": LOG_PROVIDERS.copy()
                              if params.get("provider_type") == "log"
                              else ALERT_PROVIDERS.copy()},
            "GET hosts": {"data": HOSTS.copy()},
            "GET teams": {"data": TEAMS.copy()},
            "POST incidents": {"data": INCIDENTS.copy()},
            "POST conversations.create": SLACK_CONVERSATION.copy()
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
        elif "postMessage" in url:
            endpoint = "postMessage"

        if not endpoint:
            return MockHTTPResponse(404, {}, {})

        method_endpoint = f"{method} {endpoint}"
        status = status_code.get(method_endpoint, None) or default_status.get(method_endpoint, None) or 200
        header = headers.get(method_endpoint, None) or default_headers.get(method_endpoint, None) or {}
        json = body.get(method_endpoint, None) or default_body.get(method_endpoint, None) or {}
        return MockHTTPResponse(status, header, json)
    return wrapper
