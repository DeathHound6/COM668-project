from utility import make_api_request, HTTPMethodEnum
from exceptions import ExternalAPIException
from os import getenv
from typing import Any


class SlackAPIClient:
    def create_conversation(self, channel_name: str) -> dict[str, Any]:
        response = make_api_request(
            url="https://slack.com/api/conversations.create",
            params={
                "name": channel_name,
                "is_private": "true"
            },
            method=HTTPMethodEnum.POST,
            headers={
                "Authorization": getenv('SLACK_TOKEN'),
                "Content-Type": "application/json"
            }
        )

        if response.status_code != 200 or not response.json()["ok"]:
            raise ExternalAPIException(response.json()["error"])
        return response.json()


    def invite_to_conversation(self, conversation_id: str, user_ids: set[str]) -> dict[str, Any]:
        response = make_api_request(
            url="https://slack.com/api/conversations.invite",
            params={
                "channel": conversation_id,
                "users": ",".join(user_ids),
                "force": "true"
            },
            method=HTTPMethodEnum.POST,
            headers={
                "Authorization": getenv('SLACK_TOKEN'),
                "Content-Type": "application/json"
            }
        )

        if response.status_code != 200 or not response.json()["ok"]:
            raise ExternalAPIException(response.json()["error"])
        return response.json()


    def join_conversation(self, conversation_id: str) -> dict[str, Any]:
        response = make_api_request(
            url="https://slack.com/api/conversations.join",
            params={
                "channel": conversation_id
            },
            method=HTTPMethodEnum.POST,
            headers={
                "Authorization": getenv('SLACK_TOKEN'),
                "Content-Type": "application/json"
            }
        )

        if response.status_code != 200 or not response.json()["ok"]:
            raise ExternalAPIException(response.json()["error"])
        return response.json()


    def send_message(self, conversation_id: str, message: str) -> dict[str, Any]:
        response = make_api_request(
            url="https://slack.com/api/chat.postMessage",
            params={
                "channel": conversation_id,
                "text": message
            },
            method=HTTPMethodEnum.POST,
            headers={
                "Authorization": getenv('SLACK_TOKEN'),
                "Content-Type": "application/json"
            }
        )

        if response.status_code != 200 or not response.json()["ok"]:
            raise ExternalAPIException(response.json()["error"])
        return response.json()


slack_client = SlackAPIClient()
