from http_clients.backend import backend_client
from datetime import datetime
from config import logger
from exceptions import ExternalAPIException


def incident_resolver():
    logger.info("[A.I.M.S] Scanning for incidents to resolve")
    incidents = backend_client.get_incidents({"resolved": False})
    for inc in incidents:
        last_updated = datetime.fromisoformat(inc["createdAt"])
        if len(inc["comments"]) > 0:
            # Get the last comment - these are in descending order
            last_updated = datetime.fromisoformat(inc["comments"][0]["createdAt"])

        # Incidents that have not been updated in 21 days are resolved
        # TODO: This should be configurable
        # TODO: This should check that the same incident has not been recently seen on DynaTrace/Sentry etc
        if (datetime.now() - last_updated).days > 21:
            logger.info(f"[A.I.M.S] Resolving incident {inc['uuid']}")
            try:
                reason = "Incident automatically resolved due to inactivity"
                comment_url = backend_client.post_comment_on_incident(inc["uuid"], reason)
            except ExternalAPIException as e:
                logger.exception(e)
                logger.error(f"[A.I.M.S] Failed to post comment on incident {inc['uuid']}")
                continue

            try:
                body = {
                    "resolved": True,
                    "summary": inc["summary"],
                    "description": inc["description"],
                    "hostsAffected": [host["uuid"] for host in inc["hostsAffected"]],
                    "resolutionTeams": [team["uuid"] for team in inc["resolutionTeams"]],
                }
                backend_client.update_incident(inc["uuid"], body)
            except ExternalAPIException as e:
                logger.exception(e)
                logger.error(f"[A.I.M.S] Failed to resolve incident {inc['uuid']}")
                comment_uuid = comment_url.split("/")[-1]
                try:
                    backend_client.delete_comment_on_incident(inc["uuid"], comment_uuid)
                except ExternalAPIException as e:
                    logger.exception(e)
                    logger.error(f"[A.I.M.S] Failed to delete comment {comment_uuid} on incident {inc['uuid']}")
                continue

            logger.info(f"[A.I.M.S] Incident {inc['uuid']} resolved")
