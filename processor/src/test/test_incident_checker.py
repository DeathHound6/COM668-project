from mock import patch, MagicMock
from src.processors.incident_checker.handler import incident_checker
from src.test.fixtures import mock_api_request
from src.utility import HTTPMethodEnum


class TestIncidentCheckerProcessor:
    @patch("src.http_clients.backend.make_api_request")
    def test_sentry_handler(self, mock_request: MagicMock):
        mock_request.side_effect = mock_api_request()

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
