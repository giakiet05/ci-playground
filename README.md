# ci-playground

A throwaway monorepo for practising with the shared pipeline in
[ci-templates](https://github.com/giakiet05/ci-templates). Break things here, not in real projects.

## Layout

```
compose.yaml                 the CI contract: two images, one check-only library
services/api/                Go HTTP service (standard library), imports libs/shared
services/worker/             Python background job (standard library, uv)
libs/shared/                 Go library with no image; api rebuilds when it changes
.github/workflows/pipeline.yml   calls the template; deploy: none
```

Each folder with a Makefile has `make lint` and `make test`.

```bash
docker compose up --build     # api on http://localhost:8080/hello/you
```

## Branch flow

PRs into `dev`; `dev` into `main` fast-forward only (`git push origin dev:main`); releases are
`v*` tags on `main`.

## Exercises

| # | Do this | Expect |
|---|---|---|
| 1 | PR into `dev` that edits only `services/worker/` | only `quality (worker)` and `build (worker)`; after merge, api is carried |
| 2 | Break a test in that PR | quality fails, build is skipped |
| 3 | Edit `libs/shared/` | `quality (api)`, `quality (libs/shared)`, `build (api)` |
| 4 | Edit `compose.yaml` | everything rebuilds |
| 5 | Actions tab, Run workflow on `dev` | every image rebuilt and pushed |
| 6 | PR `dev` into `main` | preflight checks `:sha-<head>` for every image |
| 7 | Fast-forward `main`, tag `v0.1.0`, push the tag | promote, then a GitHub Release |
| 8 | Switch a base image to an old version | Trivy fails the build; accept findings in `.trivyignore.yaml` |
| 9 | Push three commits to `dev` quickly | the concurrency queue in action |

Images: `ghcr.io/giakiet05/ci-playground-api` and `ghcr.io/giakiet05/ci-playground-worker`.
