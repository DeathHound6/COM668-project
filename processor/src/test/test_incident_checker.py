from unittest.mock import patch, MagicMock
from src.processors.incident_checker.handler import incident_checker
from src.test.fixtures import mock_api_request, SENTRY_HEADERS, SENTRY_EVENTS
from src.utility import HTTPMethodEnum
from copy import deepcopy


class TestIncidentCheckerProcessor:
    @patch("src.http_clients.base.APIClient.make_api_request")
    def test_sentry_handler(self, mock_request: MagicMock):
        def mock_events_node_modules(**_):
            events = deepcopy(SENTRY_EVENTS)
            for event in events:
                event["errors"][0]["data"]["url"] = "node_modules/test/index.js"
            return events

        def mock_events_node_modules_pnpm(**_):
            events = deepcopy(SENTRY_EVENTS)
            for event in events:
                event["errors"][0]["data"]["url"] = "node_modules/.pnpm/test@0.1.0"
            return events

        def mock_incidents(**_):
            return {"data": []}

        events = [lambda **_: deepcopy(SENTRY_EVENTS), mock_events_node_modules, mock_events_node_modules_pnpm]
        for event in events:
            mock_request.side_effect = mock_api_request(body={"GET events": event, "GET incidents": mock_incidents})

            incident_checker()

            # check that the whole processor has ran
            for call in mock_request.call_args_list:
                url = call.kwargs.get("url")
                method = call.kwargs.get("method")

                if "events" in url:
                    assert method == HTTPMethodEnum.GET
                elif "comments" in url:
                    assert method == HTTPMethodEnum.POST
                elif "incidents" in url:
                    assert method == HTTPMethodEnum.POST or method == HTTPMethodEnum.GET
                elif "providers" in url:
                    assert method == HTTPMethodEnum.GET
                elif "hosts" in url:
                    assert method == HTTPMethodEnum.GET
                elif "teams" in url:
                    assert method == HTTPMethodEnum.GET
                elif "conversations.create" in url:
                    assert method == HTTPMethodEnum.POST
                elif "conversations.invite" in url:
                    assert method == HTTPMethodEnum.POST
                elif "conversations.join" in url:
                    assert method == HTTPMethodEnum.POST
                elif "chat.postMessage" in url:
                    assert method == HTTPMethodEnum.POST
                elif "login" in url:
                    assert method == HTTPMethodEnum.POST

    @patch("src.http_clients.base.APIClient.make_api_request")
    @patch("src.processors.incident_checker.sentry.handle_event")
    def test_sentry_handler_all_pages(self, mock_handle_event: MagicMock, mock_request: MagicMock):
        def mock_events(**kwargs):
            offset = kwargs.get("params", {}).get("cursor", "0:0:0").split(":")[1]
            headers = deepcopy(SENTRY_HEADERS)
            if offset == "0":
                headers["Link"] = headers["Link"].replace("false", "true")
                headers["Link"] = headers["Link"].replace("0:0:0", "0:1:0")
            elif offset == "1":
                headers["Link"] = headers["Link"].replace("true", "false")
            return headers
        mock_request.side_effect = mock_api_request(headers={"GET events": mock_events})

        incident_checker()

        assert mock_handle_event.call_count == 2
