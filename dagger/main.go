// A generated module for ISCompatTable functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"dagger/i-s-compat-table/internal/dagger"
)

// TODO: BuildGoBinary
// TODO: ObservePostgres
// TODO: ObserveMySQL
// TODO: ObserveMariaDB
// TODO: ObserveTidb
// TODO: ObserveCockroachDB
// TODO: ObserveClickhouse

// TODO: ScrapePostgres
// TODO: ScrapeMySQL
// TODO: ScrapeMariaDB
// TODO: ScrapeTidb
// TODO: ScrapeCockroachDB

type Main struct{}

const builderImageRef = "docker.io/library/golang:1.22-bookworm"

// Returns a container that echoes whatever string argument is provided
func (m *Main) ContainerEcho(stringArg string) *dagger.Container {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg})
}
func (m *Main) BuildAndCacheBin(
	//+required
	main string,
	//+required
	src *dagger.Directory,
) *dagger.File {
	goBuildCache := dag.CacheVolume("GOCACHE") // from running `go env GOCACHE` inside the container
	goModCache := dag.CacheVolume("/go/pkg")
	cacheOpts := dagger.ContainerWithMountedCacheOpts{Sharing: dagger.Shared}

	return dag.Container().
		From(builderImageRef).
		WithWorkdir("/go/src").
		WithMountedCache("/root/.cache/go-build", goBuildCache, cacheOpts).
		WithMountedCache("/go/pkg", goModCache, cacheOpts).
		WithMountedDirectory("/go/src", src). // <-- ??? how does this invalidate caches?
		WithExec([]string{"go", "build", "-o", "/go/bin/out", main}).
		File("/go/bin/out")
}
