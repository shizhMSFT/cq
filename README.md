# cq - Command-line CBOR processor

> [!IMPORTANT]
> This tool is for demonstration purpose. DO NOT use in production.
>
> For production use, the following libraries are recommended:
> - CBOR: [github.com/fxamacker/cbor](https://github.com/fxamacker/cbor)
> - COSE: [github.com/veraison/go-cose](https://github.com/veraison/go-cose)

`cq` is a simple command-line CBOR parser developed from scratch.

## Build and Install

> **Note**
> Make sure `go 1.23.1` or above is installed before `make`.

To build and install `cq` to `~/bin` on Linux, simply run

```bash
make install
```

## Examples

### Usage

Get the usage of `cq`.

```console
$ cq
NAME:
   cq - Command-line CBOR processor

USAGE:
   cq [options] [file]

VERSION:
   v0.2.0

GLOBAL OPTIONS:
   --cose-payload  decode the payload of a COSE message (default: false)
   --help, -h      show help
   --version, -v   print the version
```

### Print a CBOR Object

Inspect a CBOR object in a file.

```shell
cq file.cbor
```

Alternatively, `stdin` is also accepted.

```bash
cq < file.cbor
cat file.cbor | cq
```

### Extract the Payload of a COSE Signature

Extract the payload of a `COSE_Sign1` object.

```shell
cq --cose-payload cose.sig
```
