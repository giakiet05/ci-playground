# CI contract (github.com/giakiet05/ci-templates): `lint` and `test`.

# Pinned so local runs and CI lint with the same rules; uvx fetches it on first use.
RUFF := ruff@0.17.5

.PHONY: lint test

## lint: ruff rules and formatting check
lint:
	uvx $(RUFF) check .
	uvx $(RUFF) format --check .

## test: unit tests (standard library unittest)
test:
	uv run python -m unittest discover -s tests -t . -v
