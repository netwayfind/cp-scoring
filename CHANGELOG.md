# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.0] - 2019-06-02
### Added
- [#21](https://github.com/sumwonyuno/cp-scoring/issues/21) - added insight UI to show states, reports, and diffs for scenarios and teams. Various changes to support this.
- [#25](https://github.com/sumwonyuno/cp-scoring/issues/25) - include Windows features/roles to software list
- [#2](https://github.com/sumwonyuno/cp-scoring/issues/2) - added administrators page on admin UI, allowing password change.

### Changed
- [#23](https://github.com/sumwonyuno/cp-scoring/issues/23) - Windows processes don't collect user (permissions issue); may fix later. Use executable path where available, fallback to process name.
- [#28](https://github.com/sumwonyuno/cp-scoring/issues/28) - clickable links are not underlined. Selected links are bolded, underlined, and not clickable.

### Fixed
- [#22](https://github.com/sumwonyuno/cp-scoring/issues/22) - sort Linux users in order of /etc/passwd
- [#24](https://github.com/sumwonyuno/cp-scoring/issues/24) - skip collecting software with empty name on Windows
- [#27](https://github.com/sumwonyuno/cp-scoring/issues/27) - API and UI error handling improvements

## [0.3.0] - 2019-03-30
### Changed
- updated UI dependencies
- [#16](https://github.com/sumwonyuno/cp-scoring/issues/16) - use postgres for persistent backend, remove sqlite
- [#17](https://github.com/sumwonyuno/cp-scoring/issues/17) - use configuration file instead of command line arguments
- [#18](https://github.com/sumwonyuno/cp-scoring/issues/18) - include additional information with state, save state as JSONB
- [#19](https://github.com/sumwonyuno/cp-scoring/issues/19) - use uint64 for id variables
- [#20](https://github.com/sumwonyuno/cp-scoring/issues/20) - remove hostname from team_host_tokens

### Fixed
- [#15](https://github.com/sumwonyuno/cp-scoring/issues/15) - fix issue with UI not loading, use explicit react-plotly.js version

## [0.2.0] - 2019-02-24
### Added
- [#12](https://github.com/sumwonyuno/cp-scoring/issues/12) [agents] support collecting state of IPv6 connections
- add design document
- [#10](https://github.com/sumwonyuno/cp-scoring/issues/10) add time input to user PasswordLastSet

### Changed
- rename ScenarioLatestScore struct to TeamScore, rename ScenarioScore struct to ScenarioHostScore
- update UI dependencies
- [#1](https://github.com/sumwonyuno/cp-scoring/issues/1) improve admin UI, use tabbed/page interface

### Fixed
- use React production.min.js
- [#9](https://github.com/sumwonyuno/cp-scoring/issues/9) add README instructions to set report.html ownership on Windows
- [#14](https://github.com/sumwonyuno/cp-scoring/issues/14) only show enabled scenarios in report UI
- [#13](https://github.com/sumwonyuno/cp-scoring/issues/13) fix item selection bug for scenarios on admin UI
- reorder README instructions for setting up [agents]

## [0.1.0] - 2019-01-21
### Added
- adds model package for data structs, interfaces, and helper functions
- adds agent package for collecting data on hosts
- adds Windows and Linux agents to collect users, groups, software, processes, and network connections
- adds server package for server code
- adds sqlite as a persistent backing store
- adds processing package for transferring data from agents to server
- adds auditor package for auditing host state
- adds admin web application
- adds scoreboard web application
- adds report web application
- adds build script to build server and agents
- adds bundle script to bundle server, agents, and UI into a tar.gz
