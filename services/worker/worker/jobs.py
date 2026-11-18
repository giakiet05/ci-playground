"""The worker's unit of work, kept free of I/O so it is easy to test."""

from __future__ import annotations

import re
from collections import Counter

_WORD = re.compile(r"[a-z0-9']+")


def top_words(text: str, limit: int = 3) -> list[tuple[str, int]]:
    """Count words in a text and return the most frequent ones.

    Words are lowercase runs of letters, digits and apostrophes; ties keep first-seen order.

    Args:
        text: Input text.
        limit: Maximum number of words to return.

    Returns:
        Up to ``limit`` (word, count) pairs, most frequent first.

    Raises:
        ValueError: If ``limit`` is not positive.
    """
    if limit <= 0:
        raise ValueError(f"limit must be positive, got {limit}")
    return Counter(_WORD.findall(text.lower())).most_common(limit)
