# kubectl-secret

A kubectl plugin for working with Kubernetes secrets more easily.

## Installation

```bash
go install github.com/felix-kaestner/kubectl-secret@latest
```

This installs the binary to `$GOBIN` where kubectl discovers plugins.

<details>
<summary>Build from source</summary>

```bash
make install
```

This builds the binary from source and copies it to `$GOBIN`.

</details>

## Usage

### View a secret (decoded)

```bash
kubectl secret view <name> [--namespace <ns>]
```

Fetches the named secret and prints it in a human-readable describe-style format with all data fields base64-decoded.

### Edit a secret

```bash
kubectl secret edit <name> [--namespace <ns>]
```

Opens the secret in `$EDITOR` (fallback: `vi`) with all data fields decoded.
On save, values are re-encoded and the secret is updated in the cluster.
If no changes are made, no update is issued.

## Development

```bash
make build    # build to bin/kubectl-secret
make test     # fmt + vet + test
make lint     # run golangci-lint
make install  # install to $GOBIN
```
