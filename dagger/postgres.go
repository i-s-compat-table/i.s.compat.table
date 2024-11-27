package main

import (
	"context"
	"dagger/i-s-compat-table/internal/dagger"
	"fmt"
	"strconv"
)

func (m *Main) PgService(ctx context.Context, version string) (service *dagger.Service, err error) {
	_, err = strconv.ParseFloat(version, 64) // Ensure version is parsable as a float
	if err != nil {
		return nil, err
	}
	ref := fmt.Sprintf("docker.io/library/postgres:%s-alpine", version)
	service = dag.Container().
		From(ref).
		WithExposedPort(5432).
		WithEnvVariable("POSTGRES_PASSWORD", "password").
		WithEnvVariable("PGPASSWORD", "password").
		AsService()

	return service, nil
}

func (m *Main) ObservePostgres(ctx context.Context, version string, src *dagger.Directory) (string, error) {
	const OBSERVER = "./cmd/postgres/observe"
	service, err := m.PgService(ctx, version)
	if err != nil {
		return "", err
	}
	observer := m.BuildAndCacheBin(OBSERVER, src)
	return dag.Container().
		From(builderImageRef).
		WithMountedFile("/go/bin/observer", observer).
		WithServiceBinding("db", service).
		WithExec([]string{"/go/bin/observer", version}).
		Stdout(ctx)
}

func (m *Main) BuildPgObserverBin(ctx context.Context, src *dagger.Directory) *dagger.File {
	const OBSERVER = "./cmd/postgres/observe"
	return m.BuildAndCacheBin(OBSERVER, src)
}

// func (m *Main) ScrapePostgres(ctx context.Context, version string, src *dagger.Directory) (string, error) {}
