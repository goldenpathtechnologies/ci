## [1.0.0-dev.16](https://github.com/goldenpathtechnologies/ci/compare/v1.0.0-dev.15...v1.0.0-dev.16) (2021-12-31)


### Features

* **coverage:** implemented code coverage to an acceptable threshold ([fdcd78d](https://github.com/goldenpathtechnologies/ci/commit/fdcd78d0c9bdcc0a203481e8bc88d8c1c6259331)), closes [#9](https://github.com/goldenpathtechnologies/ci/issues/9) [#2](https://github.com/goldenpathtechnologies/ci/issues/2)

## [1.0.0-dev.15](https://github.com/GoldenPathTechnologies/ci/compare/v1.0.0-dev.14...v1.0.0-dev.15) (2021-12-02)


### Features

* **documentation:** added a README to the project ([6ae455d](https://github.com/GoldenPathTechnologies/ci/commit/6ae455d1293abd10ad7f2fa207249dd50d9aa70c)), closes [#6](https://github.com/GoldenPathTechnologies/ci/issues/6)

## [1.0.0-dev.14](https://github.com/GoldenPathTechnologies/ci/compare/v1.0.0-dev.13...v1.0.0-dev.14) (2021-12-01)


### Features

* **release:** digitally signed the Windows executable ([93d525b](https://github.com/GoldenPathTechnologies/ci/commit/93d525bbdfa98e7a014772db876c3a0af0749d43)), closes [#1](https://github.com/GoldenPathTechnologies/ci/issues/1)

## [1.0.0-dev.13](https://github.com/GoldenPathTechnologies/ci/compare/v1.0.0-dev.12...v1.0.0-dev.13) (2021-12-01)


### Bug Fixes

* **release:** added LICENSE and CHANGELOG.md files to the release package ([b5c59d9](https://github.com/GoldenPathTechnologies/ci/commit/b5c59d9261783c5a1e1edcc75da0d0a277c6b4ae)), closes [#5](https://github.com/GoldenPathTechnologies/ci/issues/5)

## [1.0.0-dev.12](https://github.com/GoldenPathTechnologies/ci/compare/v1.0.0-dev.11...v1.0.0-dev.12) (2021-11-30)


### Bug Fixes

* **release:** ensured module version in PowerShell manifest file is in correct format ([34c49cc](https://github.com/GoldenPathTechnologies/ci/commit/34c49ccdc769d2ad06f8dc9705f6f3915fb9e401))

## [1.0.0-dev.11](https://github.com/GoldenPathTechnologies/ci/compare/v1.0.0-dev.10...v1.0.0-dev.11) (2021-11-29)


### Bug Fixes

* **release:** ensured changelog text is generated to an existent directory ([dd33ba1](https://github.com/GoldenPathTechnologies/ci/commit/dd33ba1478836261fc90daec50effff052fab900))

## [1.0.0-dev.10](https://github.com/GoldenPathTechnologies/ci/compare/v1.0.0-dev.9...v1.0.0-dev.10) (2021-11-29)


### Bug Fixes

* **release:** output changelog text to file due to envvar issues with markdown format ([39deadd](https://github.com/GoldenPathTechnologies/ci/commit/39deadde553beafd5c4b64d108bce7bf5a6ec90c))

## [1.0.0-dev.9](https://github.com/GoldenPathTechnologies/ci/compare/v1.0.0-dev.8...v1.0.0-dev.9) (2021-11-28)


### Features

* **release:** ensured each release contains changelog text for the current version only ([c2a595d](https://github.com/GoldenPathTechnologies/ci/commit/c2a595defcad9f640a05dd1434033e94d85ceccb))

## [1.0.0-dev.8](https://github.com/GoldenPathTechnologies/ci/compare/v1.0.0-dev.7...v1.0.0-dev.8) (2021-11-28)


### Features

* **release:** skipped Version workflow for branches when tagging releases ([8bdef97](https://github.com/GoldenPathTechnologies/ci/commit/8bdef977e50611d5a1d9fbc919f333d52adaa900))

## [1.0.0-dev.7](https://github.com/GoldenPathTechnologies/ci/compare/v1.0.0-dev.6...v1.0.0-dev.7) (2021-11-28)


### Bug Fixes

* **release:** ensured that the checksum file is named and formatted correctly ([0ce7a07](https://github.com/GoldenPathTechnologies/ci/commit/0ce7a07b4e7597405e961eea5052fac5121745fa))

## [1.0.0-dev.6](https://github.com/GoldenPathTechnologies/ci/compare/v1.0.0-dev.5...v1.0.0-dev.6) (2021-11-27)


### Features

* **release:** generated SHA-256 checksums of release package files ([906adc2](https://github.com/GoldenPathTechnologies/ci/commit/906adc203aa4160bbe8d13c5d1a0d7a899ba04aa))

## [1.0.0-dev.5](https://github.com/GoldenPathTechnologies/ci/compare/v1.0.0-dev.4...v1.0.0-dev.5) (2021-11-27)


### Features

* **release:** digitally signing all PowerShell scripts during release ([bcfff6f](https://github.com/GoldenPathTechnologies/ci/commit/bcfff6f1a9a2c5ee8d41449a2f0dbb35fe8e0226))

## [1.0.0-dev.4](https://github.com/GoldenPathTechnologies/ci/compare/v1.0.0-dev.3...v1.0.0-dev.4) (2021-11-23)


### Features

* **installation:** added PowerShell module manifest file ([03125a0](https://github.com/GoldenPathTechnologies/ci/commit/03125a0d8edfeb4fc7f9ed4889c5d7ba5a24b2ba))


### Bug Fixes

* **release:** copied correct files to Windows release package during build ([0d32392](https://github.com/GoldenPathTechnologies/ci/commit/0d32392e0738a212522b9f9375acf783429115c3))

## [1.0.0-dev.3](https://github.com/GoldenPathTechnologies/ci/compare/v1.0.0-dev.2...v1.0.0-dev.3) (2021-11-18)


### Bug Fixes

* **installation:** corrected handling of version numbers and their comparison ([e5ee2ef](https://github.com/GoldenPathTechnologies/ci/commit/e5ee2ef1e1b136d1aa626bbc8f39959ab711be1e))

## [1.0.0-dev.2](https://github.com/GoldenPathTechnologies/ci/compare/v1.0.0-dev.1...v1.0.0-dev.2) (2021-11-15)


### Bug Fixes

* **logging:** changed the location of the log file to that of the executable directory ([9e3b533](https://github.com/GoldenPathTechnologies/ci/commit/9e3b53332b515a0aa38933f61aa05a7ec688a25e))

## 1.0.0-dev.1 (2021-11-12)


### Features

* **installation:** completed installation scripts for Linux ([55154bd](https://github.com/GoldenPathTechnologies/ci/commit/55154bd6db48f663dec334706f5ce80811d2fe31))
* **installation:** completed installation scripts for Windows ([d59bcdf](https://github.com/GoldenPathTechnologies/ci/commit/d59bcdf83ce94190751278c85d0c5ae712047816))
