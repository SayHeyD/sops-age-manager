[![Test and Build](https://github.com/SayHeyD/sops-age-manager/actions/workflows/test-and-build.yaml/badge.svg)](https://github.com/SayHeyD/sops-age-manager/actions/workflows/test-and-build.yaml) [![Lint](https://github.com/SayHeyD/sops-age-manager/actions/workflows/lint.yaml/badge.svg)](https://github.com/SayHeyD/sops-age-manager/actions/workflows/lint.yaml)

# sops-age-manager (sam)

sam is a tool to easily manage your sops configuration when using multiple age keys.
This is useful when f.ex. you have a k8s cluster where you have per-namespace decryption keys.

# Table of contents
- [sops-age-manager (sam)]()
- [Table of contents]()
- [Why isn't sops enough?]()
- [What exactly does sam do?]()
- [User guide]()
  - [General]()
  - [Commands]()
  - [Configuration]()

# Why isn't sops enough?

With the tooling that sops provides currently, changing the configured age key required entering the public key
as an argument with every operation or defining an environment variable with the private key of the key to use.
Both options are rather cumbersome when having to change keys frequently.

# What exactly does sam do?

sam provides a configurable layer on top of sops. This means sam is basically a wrapper for sops when using age keys.
you can configure which key to use by name and execute sops commands with the configured key. In addition, sam also 
provides some small helper commands to manage and access your key data.

# User guide

> TODO

## General

> TODO

## Commands

> TODO

## Configuration

> TODO