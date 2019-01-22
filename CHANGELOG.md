# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
