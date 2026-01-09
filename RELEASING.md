# Releasing

This document describes the steps to release a new version of digcaa.

## Prerequisites

- You have commit access to the repository
- You have push access to the repository

## Release process

1. **Determine the new version**

   ```shell
   VERSION=0.X.0
   ```

   Use simple incremental versioning:

   - Current version: `v0.5.0`
   - Next version: `v0.6.0`, `v0.7.0`, etc.

2. **Run tests** and confirm they pass

   ```shell
   go test ./...
   ```

3. **Build** and confirm it compiles

   ```shell
   go build ./...
   ```

4. **Commit any pending changes**

   ```shell
   git add .
   git commit -m "Description of changes"
   ```

5. **Create a tag**

   ```shell
   git tag -a v$VERSION -s -m "Release $VERSION"
   ```

6. **Push the changes and tag**

   ```shell
   git push origin main
   git push origin v$VERSION
   ```

## Post-release

- Verify the new version appears on [GitHub releases](https://github.com/weppos/digcaa/releases)
- Wait a few minutes for the Go module proxy to index the new version
- Test installation:

  ```shell
  go install github.com/weppos/digcaa/cmd/digcaa@v$VERSION
  ```

- Users can now install with `@latest` to get the new version
