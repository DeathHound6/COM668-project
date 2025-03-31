from logging import getLogger, StreamHandler, DEBUG, INFO
from src.utility import LoggerFormatter
from sys import stdout
import os


env = os.getenv("ENV", "prod").lower()
assert env in ("dev", "prod"), "Env var `ENV` must be either `dev` or `prod`"
api_host = os.getenv("API_HOST")
assert api_host, "Env var `API_HOST` is required"
sentry_token = os.getenv("SENTRY_TOKEN")
assert sentry_token.startswith("sntryu_"), "Invalid env var `SENTRY_TOKEN`"

logger = getLogger(__name__)
handler = StreamHandler(stream=stdout)
handler.setFormatter(LoggerFormatter())
logger.addHandler(handler)

if env == "dev":
    logger.setLevel(DEBUG)
elif env == "prod":
    logger.setLevel(INFO)
