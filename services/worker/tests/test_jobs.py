"""Unit tests for worker.jobs."""

from __future__ import annotations

import unittest

from worker.jobs import top_words


class TopWordsTest(unittest.TestCase):
    """top_words counts case-insensitively and respects the limit."""

    def test_most_frequent_first(self) -> None:
        self.assertEqual(top_words("b a b c b a", limit=2), [("b", 3), ("a", 2)])

    def test_case_and_punctuation_are_ignored(self) -> None:
        self.assertEqual(top_words("Go, go! GO.", limit=1), [("go", 3)])

    def test_empty_text(self) -> None:
        self.assertEqual(top_words(""), [])

    def test_non_positive_limit_is_rejected(self) -> None:
        with self.assertRaises(ValueError):
            top_words("anything", limit=0)


if __name__ == "__main__":
    unittest.main()
