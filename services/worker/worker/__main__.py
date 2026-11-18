"""Entry point: run the job on a fixed interval until SIGTERM or SIGINT."""

from __future__ import annotations

import json
import logging
import os
import signal
import sys
import threading
from types import FrameType

from worker.jobs import top_words

SAMPLE_TEXT = "build once, promote the same image, build once more only when the code changes"


class JsonFormatter(logging.Formatter):
    """Format log records as one JSON object per line, including any ``extra`` fields."""

    _RESERVED = set(logging.makeLogRecord({}).__dict__)

    def format(self, record: logging.LogRecord) -> str:
        """Render a record as JSON.

        Args:
            record: The log record.

        Returns:
            A single-line JSON string.
        """
        payload = {
            "time": self.formatTime(record, "%Y-%m-%dT%H:%M:%S%z"),
            "level": record.levelname,
            "msg": record.getMessage(),
        }
        payload.update({k: v for k, v in record.__dict__.items() if k not in self._RESERVED})
        if record.exc_info:
            payload["exc"] = self.formatException(record.exc_info)
        return json.dumps(payload)


def _interval_seconds() -> float:
    """Read WORKER_INTERVAL_SECONDS (default 30).

    Returns:
        The interval between job runs.

    Raises:
        ValueError: If the variable is not a positive number.
    """
    raw = os.environ.get("WORKER_INTERVAL_SECONDS", "30")
    interval = float(raw)
    if interval <= 0:
        raise ValueError(f"WORKER_INTERVAL_SECONDS must be positive, got {raw!r}")
    return interval


def main() -> int:
    """Run the job until a stop signal arrives.

    Returns:
        Process exit code: 0 on a clean stop, 1 on bad configuration.
    """
    handler = logging.StreamHandler(sys.stdout)
    handler.setFormatter(JsonFormatter())
    logging.basicConfig(level=logging.INFO, handlers=[handler])
    logger = logging.getLogger("worker")

    try:
        interval = _interval_seconds()
    except ValueError:
        logger.exception("invalid configuration")
        return 1

    stop = threading.Event()

    def _on_signal(signum: int, _frame: FrameType | None) -> None:
        logger.info("stop signal received", extra={"signal": signal.Signals(signum).name})
        stop.set()

    signal.signal(signal.SIGTERM, _on_signal)
    signal.signal(signal.SIGINT, _on_signal)

    logger.info("worker started", extra={"interval_seconds": interval})
    while not stop.is_set():
        logger.info("job finished", extra={"top_words": top_words(SAMPLE_TEXT)})
        stop.wait(interval)
    logger.info("worker stopped")
    return 0


if __name__ == "__main__":
    sys.exit(main())
