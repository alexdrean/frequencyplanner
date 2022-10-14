# Frequency Planner
### Configuration
Add files in `sites/` to define sites. See `sites/example.json`

Edit `config.json` to define your platforms and set the webserver `(address):port`

### Execution
`go run .`

### Compiling into a binary
`GOOS=linux GOARCH=amd64 go build .`

See [supported OS and ARCH](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63)

An executable will be created under the name `frequencyplanner`