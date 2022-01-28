## 1.0.0 (2022-01-28)


### Features

* **documentation:** added a README to the project ([6ae455d](https://github.com/goldenpathtechnologies/ci/commit/6ae455d1293abd10ad7f2fa207249dd50d9aa70c)), closes [#6](https://github.com/goldenpathtechnologies/ci/issues/6)
* **help:** added a second build owner field and a repository field to help text ([60dd857](https://github.com/goldenpathtechnologies/ci/commit/60dd857a5983419a17e7ad0c69726052a20b5f9a))
* **installation:** added PowerShell module manifest file ([03125a0](https://github.com/goldenpathtechnologies/ci/commit/03125a0d8edfeb4fc7f9ed4889c5d7ba5a24b2ba))
* **help:** changed title of details view to 'Help' when the help text is displayed ([bbef895](https://github.com/goldenpathtechnologies/ci/commit/bbef895ef73af5eea639bb285aaf5a25abee10d2))
* **installation:** completed installation scripts for Linux ([55154bd](https://github.com/goldenpathtechnologies/ci/commit/55154bd6db48f663dec334706f5ce80811d2fe31))
* **installation:** completed installation scripts for Windows ([d59bcdf](https://github.com/goldenpathtechnologies/ci/commit/d59bcdf83ce94190751278c85d0c5ae712047816))
* **help:** created help text that displays in the details pane ([adb4ec5](https://github.com/goldenpathtechnologies/ci/commit/adb4ec54d09838e70d580eeaaf5d15af672d29d3)), closes [#11](https://github.com/goldenpathtechnologies/ci/issues/11)
* **release:** digitally signed the Windows executable ([93d525b](https://github.com/goldenpathtechnologies/ci/commit/93d525bbdfa98e7a014772db876c3a0af0749d43)), closes [#1](https://github.com/goldenpathtechnologies/ci/issues/1)
* **release:** digitally signing all PowerShell scripts during release ([bcfff6f](https://github.com/goldenpathtechnologies/ci/commit/bcfff6f1a9a2c5ee8d41449a2f0dbb35fe8e0226))
* disabled glob characters in filter input unless in manual glob entry mode ([7e38380](https://github.com/goldenpathtechnologies/ci/commit/7e383802e59b69b19fafcbb046fb0c9c184f9c8b))
* **release:** ensured each release contains changelog text for the current version only ([c2a595d](https://github.com/goldenpathtechnologies/ci/commit/c2a595defcad9f640a05dd1434033e94d85ceccb))
* **installation:** ensured that forceful exits of the application are indicated pleasantly ([55432ba](https://github.com/goldenpathtechnologies/ci/commit/55432ba18ad58fd5196b4d5156076c39c01b06b9))
* **release:** generated SHA-256 checksums of release package files ([906adc2](https://github.com/goldenpathtechnologies/ci/commit/906adc203aa4160bbe8d13c5d1a0d7a899ba04aa))
* **coverage:** implemented code coverage to an acceptable threshold ([fdcd78d](https://github.com/goldenpathtechnologies/ci/commit/fdcd78d0c9bdcc0a203481e8bc88d8c1c6259331)), closes [#9](https://github.com/goldenpathtechnologies/ci/issues/9) [#2](https://github.com/goldenpathtechnologies/ci/issues/2)
* introduced different filtering modes for directory list ([fccfa28](https://github.com/goldenpathtechnologies/ci/commit/fccfa285ae28c4a4cafbab0dea4becca1ddc52e4)), closes [#12](https://github.com/goldenpathtechnologies/ci/issues/12)
* set styles for filter component and handled Esc key presses ([31eea84](https://github.com/goldenpathtechnologies/ci/commit/31eea84ed4c0f82b7512af8e1635ce5306649edd))
* **release:** skipped Version workflow for branches when tagging releases ([8bdef97](https://github.com/goldenpathtechnologies/ci/commit/8bdef977e50611d5a1d9fbc919f333d52adaa900))


### Bug Fixes

* **release:** added LICENSE and CHANGELOG.md files to the release package ([b5c59d9](https://github.com/goldenpathtechnologies/ci/commit/b5c59d9261783c5a1e1edcc75da0d0a277c6b4ae)), closes [#5](https://github.com/goldenpathtechnologies/ci/issues/5)
* **installation:** changed the encoding of the PS module manifest to UTF-8 ([ad2586d](https://github.com/goldenpathtechnologies/ci/commit/ad2586dc43e479b80660d96787bc7c3d3fe848f2)), closes [#27](https://github.com/goldenpathtechnologies/ci/issues/27)
* **logging:** changed the location of the log file to that of the executable directory ([9e3b533](https://github.com/goldenpathtechnologies/ci/commit/9e3b53332b515a0aa38933f61aa05a7ec688a25e))
* **release:** copied correct files to Windows release package during build ([0d32392](https://github.com/goldenpathtechnologies/ci/commit/0d32392e0738a212522b9f9375acf783429115c3))
* **installation:** corrected handling of version numbers and their comparison ([e5ee2ef](https://github.com/goldenpathtechnologies/ci/commit/e5ee2ef1e1b136d1aa626bbc8f39959ab711be1e))
* **release:** ensured changelog text is generated to an existent directory ([dd33ba1](https://github.com/goldenpathtechnologies/ci/commit/dd33ba1478836261fc90daec50effff052fab900))
* ensured consistent directory and file sorting across components and platforms ([10627b6](https://github.com/goldenpathtechnologies/ci/commit/10627b6c5f41b6d2a2427d9f1e42fba007662d0b)), closes [#24](https://github.com/goldenpathtechnologies/ci/issues/24)
* **release:** ensured module version in PowerShell manifest file is in correct format ([34c49cc](https://github.com/goldenpathtechnologies/ci/commit/34c49ccdc769d2ad06f8dc9705f6f3915fb9e401))
* ensured symbolic links are navigable in the directory list ([d830ddd](https://github.com/goldenpathtechnologies/ci/commit/d830ddd4eaebb52bd39f3ea1b361ec5f90a490ee)), closes [#15](https://github.com/goldenpathtechnologies/ci/issues/15)
* **release:** ensured that the checksum file is named and formatted correctly ([0ce7a07](https://github.com/goldenpathtechnologies/ci/commit/0ce7a07b4e7597405e961eea5052fac5121745fa))
* **release:** output changelog text to file due to envvar issues with markdown format ([39deadd](https://github.com/goldenpathtechnologies/ci/commit/39deadde553beafd5c4b64d108bce7bf5a6ec90c))
