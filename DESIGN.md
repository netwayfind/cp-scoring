__Table of Contents__
1. [Principles](#principles)
1. [Design Requirements](#design-requirements)

# Principles

* [The Server] has [admins] that set up [templates] and [scenarios]
* [templates] are the desired configuration for [hosts]
* [hosts] are computers or virtual machines
* [scenarios] assign [templates] to [hosts]
* [teams] are given a secret API key and are given one or more [hosts] that must be set up with the secret API key
* [hosts] have [agents] that collect data (e.g. users, processes, etc.)
* [agents] send the collected data to the [The Server]
* [The Server] audits the collected data and provides a report and score for [teams] and their [hosts]

# Design Requirements
General
- software must run on 64-bit Windows and 64-bit Linux
- software must run on bare-metal computer or on virtual machine
- software must be self-contained, must not require installing dependencies on host
- software distributables must be minimal, an executable file and a few files for its configuration
- software may have different executable files for each operating system platform (Windows, Linux)
- communication must be initiated from [agents] to [The Server]; no expectation of an open network port to [hosts]

[The Server]
- [The Server] must be able to be simply copied to the intended [host]
- [The Server] must accept configuration options at start up, and must use default settings or generate files (e.g. self-signed certificate) when not provided options at start up
- [The Server] must have the following available to download:
  - [agents]
  - openpgp public key to encrypt files
  - X.509 certificate to verify identity
- [The Server] must have a admin web page
- [The Server] must have a public scoreboard
- [The Server] must have a [teams]/[hosts] report page
- [The Server] must have a HTTPS RESTful API that has an endpoint to accept data from [agents] and endpoints to support the admin web page, public scoreboard, and [teams]/[hosts] report page
- [The Server] must authenticate and authorize all API endpoints, except for the public scoreboard
- [The Server] must generate a report and a score from the [agents] data
- [The Server] must persist [agents] data, web admin interface settings, reports, and scores to a database

[agents]
- [agents] must be able to be simply copied to the [host], and have minimal configuration such as the [The Server] URL
- [agents] must periodically collect data about the host it is on
- [agents] must encrypt the collected data using the openpgp public key of the [The Server]. This encrypted data may be saved to disk until it can be sent.
- [agents] must send the encrypted data to the [The Server] and must use the given X.509 certificate to verify the connection. If the [The Server] is not available, the [agents] must try again later.