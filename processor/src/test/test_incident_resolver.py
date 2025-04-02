from unittest.mock import MagicMock, patch
from src.test.fixtures import mock_api_request, INCIDENTS
from src.processors.incident_resolver.handler import incident_resolver
from src.http_clients.base import HTTPMethodEnum
from datetime import datetime, timedelta
from copy import deepcopy
from src.exceptions import ExternalAPIException


class TestIncidentResolver:
    @patch("src.http_clients.base.APIClient.make_api_request")
    def test_resolve_incident_not_old_enough(self, mock_request: MagicMock):
        mock_request.side_effect = mock_api_request()

        incident_resolver()

        # Assert that it has not tried to resolve the incident
        assert not any(["incidents" in call.kwargs.get("url") and "comments" not in call.kwargs.get("url")
                        and call.kwargs.get("method") is HTTPMethodEnum.PUT for call in mock_request.call_args_list])

    @patch("src.http_clients.base.APIClient.make_api_request")
    def test_resolve_incident(self, mock_request: MagicMock):
        def mock_incidents(**_):
            incidents = deepcopy(INCIDENTS)
            for inc in incidents:
                inc["createdAt"] = (datetime.now() - timedelta(days=22)).isoformat()
                inc["comments"] = []
            return {"data": incidents}
        mock_request.side_effect = mock_api_request(body={"GET incidents": mock_incidents})

        incident_resolver()

        # Assert that it has tried to update the incident and post a comment
        assert any(["incidents" in call.kwargs.get("url") and "comments" not in call.kwargs.get("url")
                    and call.kwargs.get("method") is HTTPMethodEnum.PUT for call in mock_request.call_args_list])
        assert any(["comments" in call.kwargs.get("url") and call.kwargs.get("method") is HTTPMethodEnum.POST
                   for call in mock_request.call_args_list])

        for call in mock_request.call_args_list:
            url = call.kwargs.get("url")
            if "incidents" in url and "comments" not in url:
                if call.kwargs.get("method") == HTTPMethodEnum.PUT:
                    assert call.kwargs.get("body", {}).get("resolved") is True
            elif "comments" in url:
                assert call.kwargs.get("method") == HTTPMethodEnum.POST
                assert call.kwargs.get("body") == {"comment": "Incident automatically resolved due to inactivity"}

    @patch("src.http_clients.base.APIClient.make_api_request")
    def test_resolve_incident_delete_comment_on_update_fail(self, mock_request: MagicMock):
        def mock_put_incident(**_):
            raise ExternalAPIException("test")

        def mock_incidents(**_):
            incidents = deepcopy(INCIDENTS)
            for inc in incidents:
                inc["createdAt"] = (datetime.now() - timedelta(days=22)).isoformat()
                inc["comments"] = []
            return {"data": incidents}
        mock_request.side_effect = mock_api_request(body={
            "GET incidents": mock_incidents,
            "PUT incidents": mock_put_incident
        })

        incident_resolver()

        # Assert that it has tried to update the incident and post a comment
        assert any(["incidents" in call.kwargs.get("url") and "comments" not in call.kwargs.get("url")
                    and call.kwargs.get("method") is HTTPMethodEnum.PUT for call in mock_request.call_args_list])
        assert any(["comments" in call.kwargs.get("url") and call.kwargs.get("method") is HTTPMethodEnum.POST
                    for call in mock_request.call_args_list])
        # Assert that it has tried to delete a comment
        assert any(["comments" in call.kwargs.get("url") and call.kwargs.get("method") is HTTPMethodEnum.DELETE
                    for call in mock_request.call_args_list])
