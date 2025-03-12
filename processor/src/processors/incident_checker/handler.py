from config import logger
from .sentry import handle_sentry
from http_clients.backend import backend_client


def incident_checker():
    log_providers = backend_client.get_providers("log")
    alert_providers = backend_client.get_providers("alert")

    for provider in log_providers:
        enabled_field = [field for field in provider["fields"] if field["key"] == "enabled"]
        if not enabled_field:
            logger.warning(f"Field 'enabled' not found in provider '{provider['name']}'. Skipping...")
            continue
        if not bool(enabled_field[0]["value"]):
            logger.info(f"Provider '{provider['name']}' is disabled. Skipping...")
            continue

        try:
            if str(provider["name"]).lower() == "sentry":
                logger.info("Begin scanning Sentry")
                handle_sentry(provider, alert_providers)
        except Exception as e:
            logger.exception(e, stack_info=True)
