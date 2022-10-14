# Frequency Planner
### Configuration
1. Add files in `sites/` to define sites. See `sites/example.json`
2. Edit `config.json` to define your platforms and set the webserver `(address):port`
3. Edit the theme as necessary in `index.html`

### Execution
1. [Install go](https://go.dev/doc/install)

2. `go run .`

Note: changes in `config.yml` require a restart. Changes in `index.html` and `sites/` do not require a restart.

### Compiling into a binary
`GOOS=linux GOARCH=amd64 go build .`

See [supported OS and ARCH](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63)

An executable will be created under the name `frequencyplanner`