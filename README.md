cp-scoring
-----

Scoring a bunch of computers.

__Table of Contents__
1. [Intro](#Intro)
1. [Building](#Building)
1. [Deploying and Running](#Deploying\ and\ Running)

# Intro

This is an automated scoring system. The goal is to verify that a set of computers meet standards. The primary use case is to train people on setting up a host computer to a desired configuration, provide them with a computer report + score, and track their progress.

## Principles

* [The Server] has [admins] that set up [templates] and [scenarios]
* [templates] are the desired configuration for [hosts]
* [hosts] are computers or virtual machines
* [scenarios] assign [templates] to [hosts]
* [teams] are given a secret API key and are given one or more [hosts] that must be set up with the secret API key
* [hosts] have [agents] that collect data (e.g. users, processes, etc.)
* [agents] send the collected data to the [The Server]
* [The Server] audits the collected data and provides a report and score for [teams] and their [hosts]

## Design Requirements
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

## Comparison to other tools

There are commercial tools that generally do the similar work to collect data from agents on hosts and generate reports (e.g. Nessus, Rapid7). One motivation was to find software with no purchase, license, or maintenance costs.

There are open-source software, but none of the ones researched met the requirements. Chef and Puppet support large enterprise environments and seemed too heavyweight for the primary use case. Ansible is built for actively configuring other hosts from a master host, so it did not fit the primary use case. Ansible did not meet operational requirements, such as requiring (many) python dependencies and requiring open network ports (e.g. SSH, WinRM) on the hosts.

## Software Components

[The Server] is a binary executable file written in go. [The Server] has a RESTful API that its UI and the [agents] interact with. The UI is a React (JSX) web application.

[agents] are binary executable files written in go. [agents] are built for a particular supported 64-bit operating system:
* Windows
* Linux

# Building

## Language Dependencies

* golang stable (1.10.3+)
* Node.js LTS (8.11.4+)

## UI Dependencies

* react-plotly.js (2.2.0+)
* plotly.js (1.40.1+)

`npm install react-plotly plotly.js`

## Procedure

To build [The Server] and [agents]:

1. `cd cp-scoring`
1. `sh build.sh`

To create a bundle for deployment:
1. `sh bundle.sh`

# Deploying and Running

## [The Server]

Copy the generated cp-scoring.tar.gz (from cp-scoring/target) to the intended location. The archive should have the following:
* public
  * cp-scoring-agent-linux
  * cp-scoring-agent-windows
* ui
* cp-scoring-server-linux

Currently, only Linux is the supported OS to running [The Server].

To extract files:

`tar xzvf cp-scoring.tar.gz`

To run server:

`./cp-scoring-server-linux`

Options:

- -cert: path to X.509 certificate (default: public/server.crt)
- -key: path to RSA private key (default: private/server.key)
- -port: TCP port to listen on (default: 8443)

When [The Server] starts up for the first time, it will set up its database, and generate a private key and public key for encrypting data.

[The Server] requires a HTTPS certificate and key to run. By default, it expects the certificate at public/server.crt and the key at private/server.key. Either create an RSA private key and X.509 certificate at these paths or specify file paths to these files with the -cert and -key arguments.

When the server is ready, [The Server] should be available over https://\<server\>:\<port\> . By default, the port is 8443. The port number can be changed with the -port argument.

## [agents]

Download from https://\<server\>:\<port\>/public . Choose the correct [agent] for the [host].

Windows instructions:
1. Download cp-scoring-agent-windows.exe
1. Open Administrator Command Prompt
1. `cd C:\Users\<user>\Downloads`
1. `cp-scoring-agent-windows.exe -install`
   - This will install files to C:\cp-scoring
1. `cd C:\cp-scoring`
1. `cp-scoring-agent-windows.exe -server https://<server>:<port>`
   - This will download files to allow sending data to [The Server]
1. Restart computer. [agent] will automatically start.
1. If everything is OK, delete cp-scoring-agent-windows.exe in the Downloads folder
1. Create shortcut for C:\cp-scoring\scoreboard.html to Desktop
1. Create shortcut for C:\cp-scoring\report.html to Desktop

Linux instructions:

1. Download cp-scoring-agent-linux
1. Open Terminal
1. `cd ~/Downloads`
1. `chmod +x cp-scoring-agent-linux`
1. `sudo ./cp-scoring-agent-linux -install`
   - This will install files to /opt/cp-scoring
1. `cd /opt/cp-scoring`
1. `sudo cp-scoring-agent-linux -server https://<server>:<port>`
   - This will download files to allow sending data to [The Server]
1. Restart computer. [agent] will automatically start.
1. If everything is OK, delete cp-scoring-agent-linux in the Downloads folder
1. `ln -s /opt/cp-scoring/scoreboard.html ~/Desktop`
1. `ln -s /opt/cp-scoring/report.html ~/Desktop`
