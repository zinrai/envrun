# envrun

A command-line tool that applies named environment variable profiles to a command execution.

## Configuration

See [config.yaml.example](./config.yaml.example) for a full example.

A profile must specify exactly one of `vars` or `file`.

- `vars`: Inline KEY=VALUE definitions with `$VAR` expansion.
- `file`: Path to a `.env` file.

By default, envrun reads `envrun.yaml` from the current directory. Use `-config` to specify a different path.

## Usage

Apply a profile to a command:

```
$ envrun -profile <name> -- <command> [args...]
```

Specify a custom configuration file:

```
$ envrun -profile <name> -config /path/to/envrun.yaml -- <command> [args...]
```

Load a `.env` file directly without a YAML configuration:

```
$ ENVRUN_FILE=".env.prod" envrun -- <command> [args...]
```

Specifying both `-profile` and `ENVRUN_FILE` is an error.

### Examples

Use tfenv without adding it to your shell profile:

```
$ envrun -profile tfenv -- terraform plan
```

Route commands through a proxy:

```
$ envrun -profile proxy -- kubectl get pods
```

Use a remote Docker daemon:

```
$ envrun -profile docker -- docker ps
```

Load environment variables from a file:

```
$ ENVRUN_FILE=".env.staging" envrun -- ./deploy.sh
```

## Environment File Format

```
# Comment lines are ignored
KEY=value
SECRET="quoted value"
PASSWORD='also quoted'
COMPOUND="$HOME/path"
```

- Empty lines and lines starting with `#` are ignored.
- Values are split on the first `=`.
- Surrounding double or single quotes on values are stripped.
- `$VAR` and `${VAR}` are expanded.

## License

This project is licensed under the [MIT License](./LICENSE).
