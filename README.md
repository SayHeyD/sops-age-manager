[![Test and Build](https://github.com/SayHeyD/sops-age-manager/actions/workflows/test-and-build.yaml/badge.svg?branch=dev)](https://github.com/SayHeyD/sops-age-manager/actions/workflows/test-and-build.yaml) [![Lint](https://github.com/SayHeyD/sops-age-manager/actions/workflows/lint.yaml/badge.svg)](https://github.com/SayHeyD/sops-age-manager/actions/workflows/lint.yaml?branch=dev)

# sops-age-manager (sam) _in-development_

sam is a tool to easily manage your sops configuration when using multiple age keys.
This is useful when f.ex. you have a k8s cluster where you have per-namespace decryption keys.

# Table of contents
- [sops-age-manager (sam)](#sops-age-manager-sam-in-development)
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

sam provides a configurable layer on top of sops. This means sam is a 
wrapper for sops itself and other applications that use sops under the hood. f.ex. 
the sops terraform provider. 

# User guide

## Prerequisites

sam doesn't directly require [sops](https://github.com/mozilla/sops) to be installed
before it can be used but without it, sam is kinda useless.

[age](https://github.com/FiloSottile/age) isn't per se a requirement, 
but you will already need to have age keys to use sam. Sam will not create age keys for you.

## General

After installation, add the age key files to the following path ```$HOME/.age/```. sam will detect age keys
in this directory automatically by default. The filename should follow the following format: ```<KEY_NAME>.txt```.

The default config file for sam will be created at ```$HOME/.sops-age-manager/config.yaml``` on first usage of sam
if it doesn't exist already.

## Installation

Download the binary for your OS from the releases page on GitHub.

Make sure to set the active key before using sam, 
if not sops will return an error and sam will return the following error.

```
Could not find decryption key ""
Could not find encryption key ""
```

## Commands

The base command of sam does nothing by itself without a ```--``` separator after which you can 
execute whatever you want. The base command simply sets the ```SOPS_AGE_KEY``` environment variable to 
the correct value. For sops commands the ```--age``` argument will be injected automatically to the selected key.

### Examples

```bash
sam key use private-helm-manifest
sam -- sops -d super-secret.enc.yaml
```

```bash
sam key use private-helm-manifest
sam -- sops -e super-secret.dec.yaml
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