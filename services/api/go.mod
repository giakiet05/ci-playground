module github.com/giakiet05/ci-playground/services/api

go 1.27

require github.com/giakiet05/ci-playground/libs/shared v0.0.0

// The shared library lives in this repository; the Docker build context is the repo root so the
// Dockerfile can copy it next to the service.
replace github.com/giakiet05/ci-playground/libs/shared => ../../libs/shared
