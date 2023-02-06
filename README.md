[![Test and Build](https://github.com/SayHeyD/sops-age-manager/actions/workflows/test-and-build.yaml/badge.svg?branch=dev)](https://github.com/SayHeyD/sops-age-manager/actions/workflows/test-and-build.yaml) [![Lint](https://github.com/SayHeyD/sops-age-manager/actions/workflows/lint.yaml/badge.svg)](https://github.com/SayHeyD/sops-age-manager/actions/workflows/lint.yaml?branch=dev)

# sops-age-manager (sam)

sam is a tool to easily manage your sops configuration when using multiple age keys.
This is useful when f.ex. you have a k8s cluster where you have per-namespace decryption keys.

# Table of contents
- [sops-age-manager (sam)](#sops-age-manager-sam)
- [Table of contents](#table-of-contents)
- [Why isn't sops enough?](#why-isnt-sops-enough)
- [What exactly does sam do?](#what-exactly-does-sam-do)
- [User guide](#user-guide)
  - [Prerequisites](#prerequisites)
  - [General](#general)
  - [Installation](#installation)
  - [Commands](#commands)
  - [Configuration](#configuration)

# Why isn't sops enough?

With the tooling that sops provides currently, changing the configured age key required entering the public key
as an argument with every operation or defining an environment variable with the private key of the key to use.
Both options are rather cumbersome when having to change keys frequently.

# What exactly does sam do?

sam provides a configurable layer on top of sops. This means sam is basically a wrapper for sops when using age keys.
you can configure which key to use by name and execute sops commands with the configured key. In addition, sam also 
provides some small helper commands to manage and access your key data.

# User guide

## Prerequisites

sam requires [sops](https://github.com/mozilla/sops) to be installed before it can be used.
If sops is not installed everything still works as expected aside from the base command, which passes
its args to sops. sam also requires sops to be in the PATH.

## General

After installation add the age key files to the following path ```$HOME/.age/```. sam will detect age keys
in this directory automatically by default. The filename should follow the following format: ```<KEY_NAME>.txt```.

The default config file for sam will be created at ```$HOME/.sops-age-manager/config.yaml``` on first usage of sam
if it doesn't exist already.

## Installation

Download the binary for your OS from the releases page on GitHub.

## Commands

The base command of sam just calls sops with the passed arguments:

```bash
sam -- <SOPS_ARGS_HERE>
```

The ```--age``` argument is passed automatically by sam.

__COMMAND DOCUMENTATION:__

- [SAM](./docs/sam.md)
  - [Config](./docs/sam_config.md)
    - [Dump](./docs/sam_config_dump.md)
    - [Path](./docs/sam_config_path.md)
  - [Key](./docs/sam_key.md)
    - [List](./docs/sam_key_list.md)
    - [Use](./docs/sam_key_use.md)

## Configuration

Configuration is quite minimal and lets you configure the following values:

- [encryption-key](#encryption-key)
- [decryption-key](#decryption-key)
- [key-dir](#key-dir)

### Encryption Key

The name of the encryption key to use. This is passed to sops as the ```--age``` arg to sops.
Available key names can be listed with the [sam key list](./docs/sam_key_list.md) command.

### Decryption Key

The name of the decryption key to use. This is set as the value of the ```SOPS_AGE_KEY```
environment variable which is consumed by sops.
Available key names can be listed with the [sam key list](./docs/sam_key_list.md) command.

### Key dir

The directory where the age keys are stored. This has to be an absolute filepath. Environment variables are not parsed.

All keys that are not directly in the key-dir i.e. in subfolders will not be detected by sam.