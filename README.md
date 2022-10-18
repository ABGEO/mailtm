<div align="center">
<img width="600px" src="./.assets/logo.png" />

<p>CLI client for <a href="https://mail.tm">Mail.tm</a> disposable mail service</p>

<p>
    <a href="https://github.com/abgeo/mailtm/releases">
        <img alt="GitHub release" src="https://img.shields.io/github/release/abgeo/mailtm.svg">
    </a>
    <a href="https://github.com/abgeo/mailtm">
        <img alt="GitHub go.mod Go version of a Go module" src="https://img.shields.io/github/go-mod/go-version/abgeo/mailtm.svg">
    </a>
    <a href="https://github.com/abgeo/mailtm/actions/workflows/ci.yaml">
        <img alt="CI" src="https://github.com/abgeo/mailtm/actions/workflows/ci.yaml/badge.svg">
    </a>
    <a href="https://github.com/abgeo/mailtm/actions/workflows/cd.yaml">
        <img alt="CD" src="https://github.com/abgeo/mailtm/actions/workflows/cd.yaml/badge.svg">
    </a>
    <a href="https://app.fossa.com/projects/custom%2B30026%2Fgithub.com%2FABGEO%2Fmailtm?ref=badge_shield" alt="FOSSA Status">
        <img src="https://app.fossa.com/api/projects/custom%2B30026%2Fgithub.com%2FABGEO%2Fmailtm.svg?type=shield"/>
    </a>
    <a href="https://sonarcloud.io/project/overview?id=ABGEO_mailtm">
        <img alt="Quality Gate Status" src="https://sonarcloud.io/api/project_badges/measure?project=ABGEO_mailtm&metric=alert_status"/>
    </a>
    <a href="https://sonarcloud.io/project/overview?id=ABGEO_mailtm">
        <img alt="Maintainability Rating" src="https://sonarcloud.io/api/project_badges/measure?project=ABGEO_mailtm&metric=sqale_rating"/>
    </a>
    <a href="https://sonarcloud.io/project/overview?id=ABGEO_mailtm">
        <img alt="Reliability Rating" src="https://sonarcloud.io/api/project_badges/measure?project=ABGEO_mailtm&metric=reliability_rating"/>
    </a>
    <a href="https://sonarcloud.io/project/overview?id=ABGEO_mailtm">
        <img alt="Security Rating" src="https://sonarcloud.io/api/project_badges/measure?project=ABGEO_mailtm&metric=security_rating"/>
    </a>
    <a href="https://codecov.io/gh/ABGEO/mailtm">
     <img src="https://codecov.io/gh/ABGEO/mailtm/branch/main/graph/badge.svg?token=TC7WWTT2A5"/>
     </a>
    <a href="https://goreportcard.com/report/github.com/ABGEO/mailtm">
        <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/ABGEO/mailtm">
    </a>
</p>

<p><img width="1000" src="./.assets/demo.gif" /></p>
</div>

## Installation

- Go to the [Releases](https://github.com/ABGEO/mailtm/releases) page and download the version suitable for your OS.
- Extract `mailtm` binary file from the archive: `tar -xzf mailtm_*.tar.gz mailtm`
- Make `mailtm` file executable: `chmod +x mailtm`
- Move `mailtm` to a location in your `PATH`: `sudo mv mailtm /usr/local/bin/`

### Verify Installation

`mailtm` releases are signed using PGP key (rsa4096) with fingerprint 
`5B8D 6B31 D430 43AD 711C  7C10 0E28 CC94 816E 5E0C`. Our key can be retrieved from common keyservers.

```shell
# Import key.
curl -s 'https://keys.openpgp.org/vks/v1/by-fingerprint/5B8D6B31D43043AD711C7C100E28CC94816E5E0C' | gpg --import

# Verify signature.
gpg --verify mailtm_*_checksums.txt.sig mailtm_*_checksums.txt

# Verify checksum.
sha256sum --ignore-missing -c mailtm_*_checksums.txt
```

## Usage

Get available commands by running `mailtm --help`

## Authors

- [Temuri Takalandze](https://abgeo.dev) - *Maintainer*

## License

Copyright (c) 2022 [Temuri Takalandze](https://abgeo.dev).  
Released under the [GPL-3.0](LICENSE) license.

[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B30026%2Fgithub.com%2FABGEO%2Fmailtm.svg?type=large)](https://app.fossa.com/projects/custom%2B30026%2Fgithub.com%2FABGEO%2Fmailtm?ref=badge_large)
