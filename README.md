# Audiobook Player (name pending)

A TUI audiobook player written in Go.

## Contributing:
Check the [PLANNING.md](PLANNING.md) file for the current planned features/roadmap.

If you are using WSL, keep in mind that WSL doesn't support playing audio directly. You can cross compile for Windows by using the enviroment variables: 
```bash
GOOS=windows GOARCH=amd64 go build -o bin/abs.exe ./cmd
```

## Dev Requirements:
- Go 1.24 (may work on earlier version, but not tested)
- [SQLC](https://sqlc.dev/) 1.29

## Building:
### Preparing:
```bash
sqlc generate
go mod tidy
```

### Building:
```bash
go build -o bin/abs ./cmd
```

### Running 
```bash
go run ./cmd
``` 

Currently, we are doing experiments, so there's nothing interesting in the `cmd` directory yet. Try the `experiments` directory for some small experiments that have not been integrated yet.
