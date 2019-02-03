__Table of Contents__
1. [Principles](#principles)
1. [Design Requirements](#design-requirements)
1. [Persistent Backing Store](#persistent-backing-store)

# Principles

* [The Server] has [admins] that set up [templates] and [scenarios]
* [templates] are the desired configuration for [hosts]
* [hosts] are computers or virtual machines
* [scenarios] assign [templates] to [hosts]
* [teams] are given a secret API key and are given one or more [hosts] that must be set up with the secret API key
* [hosts] have [agents] that collect data (e.g. users, processes, etc.)
* [hosts] obtain a unique host token from [The Server] to uniquely identify the instance
* [agents] send the collected data and host token to the [The Server]
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
- [The Server] must persist [agents] data, web admin interface settings, reports, and scores to a persistent backing store

[agents]
- [agents] must be able to be simply copied to the [host], and have minimal configuration such as the [The Server] URL
- [agents] must periodically collect data about the host it is on
- [agents] must encrypt the collected data using the openpgp public key of the [The Server]. This encrypted data may be saved to disk until it can be sent.
- [agents] must send the encrypted data to the [The Server] and must use the given X.509 certificate to verify the connection. If the [The Server] is not available, the [agents] must try again later.

# Persistent Backing Store
The persistent backing store saves the data for [The Server] and from the [agents]. These are the available implementations:

- sqlite

sqlite is intended to be for development and small test environments with a few [hosts] and [teams], limited time operation (several hours), and infrequent data access (a few times a minute).

See [persistence.go](server/persistence.go) for interface. The persistent backing store must handle the following items and actions:

- [admins]
  - add
  - update password hash
  - delete
  - select all usernames
  - check if there is an account with username
  - retrieve password hash for username. This is for authentication.
- [hosts]
  - add
  - update
  - delete
  - select all
  - select instance with matching ID
  - retrieve ID for hostname (first match)
- [hosts] state data
  - add
  - select all that match criteria (e.g. [hosts], [teams], IP address of submission, hostname, ID, time), ordered by state timestamp ascending
  - delete all for a single [hosts]; no single state data instance allowed
- [teams]
  - add
  - update
  - delete
  - select all
  - select instance with matching ID
  - retrieve ID for secret API key
  - retrieve ID for host token
- [templates]
  - add
  - update
  - delete
  - select all
  - select all associated with [scenarios] ID and hostname
  - select instance with matching ID
- [scenarios]
  - add
  - update
  - delete
  - select all (optional: filter for enabled)
  - select all associated with hostname
  - select instance with matching ID
- reports
  - add
  - select latest for [scenarios] ID and host token
- scores
  - add
  - select scores over time for [scenarios] ID and hostname
  - select latest for [scenarios] ID
- host tokens
  - add
  - associate with [teams] ID and hostname
  - select instance with matching [teams] ID and hostname
