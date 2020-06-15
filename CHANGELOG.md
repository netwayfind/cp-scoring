# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.6.0] - 2020-06-14
### Added
- [#37](https://github.com/netwayfind/cp-scoring/issues/37) - insight UI, add button to show/hide findings
- [#4](https://github.com/netwayfind/cp-scoring/issues/4) - initial support for scheduled tasks in Windows
- [#3](https://github.com/netwayfind/cp-scoring/issues/3) - initial support for Windows Firewall
- [#5](https://github.com/netwayfind/cp-scoring/issues/5) - initial support for Windows settings
- [#33](https://github.com/netwayfind/cp-scoring/issues/33) - adds clone button for existing hosts, templates, and scenarios
- [#36](https://github.com/netwayfind/cp-scoring/issues/36) - adds shortcut on Desktop to team key registration
- [#11](https://github.com/netwayfind/cp-scoring/issues/11) - adds shortcut on Desktop to report and scoreboard

### Changed
- generated team key is only upper case letters and numbers
- Windows agent runs under SYSTEM account to allow more actions
- updated repository URLs
- updated npm dependencies

### Fixed
- [#23](https://github.com/netwayfind/cp-scoring/issues/23) - able to retrieve user for Windows processes, due to running with SYSTEM account privileges
- [#32](https://github.com/netwayfind/cp-scoring/issues/32) - adds behavior to press Enter key before exiting team key registration. This allows seeing success or any errors.
- [#41](https://github.com/netwayfind/cp-scoring/issues/41) - fixes parsing PowerShell versino

## [0.5.0] - 2019-09-01
### Added
- [#31](https://github.com/netwayfind/cp-scoring/issues/31) - create template from state ID
- [#8](https://github.com/netwayfind/cp-scoring/issues/8) - adds scenario description page
- [#35](https://github.com/netwayfind/cp-scoring/issues/35) - added refresh button to scoreboard and report UI

### Changed
- updated README
- updated UI dependencies
- minor insight UI layout improvements
- [#30](https://github.com/netwayfind/cp-scoring/issues/30) - updated diff code to have report/state diff response streamed to reduce response times
- [#32](https://github.com/netwayfind/cp-scoring/issues/32) - moved team/host registration from report UI to agent

### Fixed
- fixed random seed
- fixed missing await keywords in report UI
- [#29](https://github.com/netwayfind/cp-scoring/issues/29) - prevent deleting current or last admin user
- [#34](https://github.com/netwayfind/cp-scoring/issues/34) - Windows features are only collected on Windows Server. Without this fix, software could not be collected on non-Server Windows hosts.

## [0.4.0] - 2019-06-02
### Added
- [#21](https://github.com/netwayfind/cp-scoring/issues/21) - added insight UI to show states, reports, and diffs for scenarios and teams. Various changes to support this.
- [#25](https://github.com/netwayfind/cp-scoring/issues/25) - include Windows features/roles to software list
- [#2](https://github.com/netwayfind/cp-scoring/issues/2) - added administrators page on admin UI, allowing password change.

### Changed
- [#23](https://github.com/netwayfind/cp-scoring/issues/23) - Windows processes don't collect user (permissions issue); may fix later. Use executable path where available, fallback to process name.
- [#28](https://github.com/netwayfind/cp-scoring/issues/28) - clickable links are not underlined. Selected links are bolded, underlined, and not clickable.

### Fixed
- [#22](https://github.com/netwayfind/cp-scoring/issues/22) - sort Linux users in order of /etc/passwd
- [#24](https://github.com/netwayfind/cp-scoring/issues/24) - skip collecting software with empty name on Windows
- [#27](https://github.com/netwayfind/cp-scoring/issues/27) - API and UI error handling improvements

## [0.3.0] - 2019-03-30
### Changed
- updated UI dependencies
- [#16](https://github.com/netwayfind/cp-scoring/issues/16) - use postgres for persistent backend, remove sqlite
- [#17](https://github.com/netwayfind/cp-scoring/issues/17) - use configuration file instead of command line arguments
- [#18](https://github.com/netwayfind/cp-scoring/issues/18) - include additional information with state, save state as JSONB
- [#19](https://github.com/netwayfind/cp-scoring/issues/19) - use uint64 for id variables
- [#20](https://github.com/netwayfind/cp-scoring/issues/20) - remove hostname from team_host_tokens

### Fixed
- [#15](https://github.com/netwayfind/cp-scoring/issues/15) - fix issue with UI not loading, use explicit react-plotly.js version

## [0.2.0] - 2019-02-24
### Added
- [#12](https://github.com/netwayfind/cp-scoring/issues/12) [agents] support collecting state of IPv6 connections
- add design document
- [#10](https://github.com/netwayfind/cp-scoring/issues/10) add time input to user PasswordLastSet

### Changed
- rename ScenarioLatestScore struct to TeamScore, rename ScenarioScore struct to ScenarioHostScore
- update UI dependencies
- [#1](https://github.com/netwayfind/cp-scoring/issues/1) improve admin UI, use tabbed/page interface

### Fixed
- use React production.min.js
- [#9](https://github.com/netwayfind/cp-scoring/issues/9) add README instructions to set report.html ownership on Windows
- [#14](https://github.com/netwayfind/cp-scoring/issues/14) only show enabled scenarios in report UI
- [#13](https://github.com/netwayfind/cp-scoring/issues/13) fix item selection bug for scenarios on admin UI
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
