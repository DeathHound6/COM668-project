from src.config import logger
from multiprocessing import Process
from time import sleep
from src.processors.incident_checker.handler import incident_checker
from src.processors.incident_resolver.handler import incident_resolver
from typing import Any


def run(**kwargs):
    func = kwargs.get("func", lambda: None)
    wait_time = kwargs.get("wait_time", 30)
    while True:
        func()
        sleep(int(wait_time))


def main():
    # Give the API time to start
    logger.info("Waiting for the API to start")
    sleep(5)

    processes = [
        {
            "name": "Incident Checker",
            "args": {
                "func": incident_checker,
                "wait_time": 30
            }
        },
        {
            "name": "Incident Resolver",
            "args": {
                "func": incident_resolver,
                "wait_time": 60
            }
        }
    ]
    running_processes: list[tuple[Process, dict[str, Any]]] = []

    # Startup the processes
    for proc in processes:
        p = Process(target=run, name=proc["name"], kwargs=proc["args"])
        p.start()
        running_processes.append((p, proc))

    # Keep-alive
    for proc, conf in running_processes:
        if not proc.is_alive():
            logger.warning(f"Process '{proc.name}' is not running. Restarting...")
            running_processes.remove((proc, conf))
            p = Process(target=run, name=proc.name, kwargs=conf["args"])
            p.start()
            running_processes.append((p, conf))
        sleep(2)


if __name__ == "__main__":
    main()
