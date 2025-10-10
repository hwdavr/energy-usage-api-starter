#!/usr/bin/env bash
set -eu

OUT=.ai_context
mkdir -p "$OUT"

# 1) Trimmed tree (skip vendor, binaries, testdata dumps)
git ls-files | egrep -v '(^vendor/|^bin/|^dist/|\.png$|\.jpg$|\.pdf$|\.md$|^testdata/.*\.json$)' \
  | head -n 400 > "$OUT/repo_tree.txt"

# 2) Go module + package graph
go env GOPATH >/dev/null  # warm modules
go list -m -json all > "$OUT/go_mods.json" || true
go list ./... > "$OUT/packages.txt" || true

# 3) Entrypoints & wiring
rg -n --no-heading "func main\(|Provide|NewServer\(|Route\(|Register\(|Wire" \
  --glob '!vendor' --glob '!**/*_test.go' > "$OUT/entrypoints_grep.txt" || true

# 4) HTTP routes (chi/echo/gin examples)
rg -n --no-heading "Route|HandleFunc|GET\(|POST\(|PUT\(|DELETE\(" \
  internal cmd pkg --glob '!**/*_test.go' > "$OUT/routes_grep.txt" || true

# 5) DB layer (sqlc / gorm / raw SQL)
rg -n --no-heading "sqlc|gorm|DB\.Query|DB\.Exec|BEGIN|COMMIT|ROLLBACK" \
  --glob '!**/*_test.go' > "$OUT/db_grep.txt" || true

# 6) Concurrency & context usage
rg -n --no-heading "go\s+|context\.With(Timeout|Cancel|Deadline)" \
  --glob '!**/*_test.go' > "$OUT/concurrency_grep.txt" || true

# 7) Error/observability patterns
rg -n --no-heading "errors\.|fmt\.Errorf|Wrap|logger|zap|slog|metrics|prom" \
  --glob '!**/*_test.go' > "$OUT/err_obs_grep.txt" || true

echo "Wrote AI context to $OUT/"
