if ! go test ./... -cover -race; then
  echo "Tests failed" >&2
  exit 1
fi

echo "Tests passed"