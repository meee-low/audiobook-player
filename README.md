# Audiobook Player (name pending)

A TUI audiobook player written in Go.

## Contributing:
Check the [PLANNING.md](PLANNING.md) file for the current planned features/roadmap.

If you are using WSL, keep in mind that WSL doesn't support playing audio directly. You can cross compile for Windows by using the enviroment variables: `GOOS=windows GOARCH=amd64`

## Dev Requirements:
- Go 1.24 (may work on earlier version, but not tested)
- [SQLC](https://sqlc.dev/) 1.29

## Building:
```bash
sqlc generate
go mod tidy
go build -o bin/ ./cmd
```

Currently, we are doing experiments, so there's nothing interesting in the `cmd` directory yet. Try the `experiments` directory for some small experiments that have not been integrated yet.
