cp-scoring
-----

Scoring a bunch of computers.

__Table of Contents__
1. [Intro](#intro)
1. [Building](#building)
1. [Deploying and Running](#deploying-and-running)

# Intro

This is an automated scoring system. The goal is to verify that a set of computers meet standards. The primary use case is to train people on setting up a host computer to a desired configuration, provide them with a computer report + score, and track their progress.

[Design](DESIGN.md)


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

Copy the generated cp-scoring-\<version\>.tar.gz (from cp-scoring/target) to the intended location. The archive should have the following:
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
1. `icacls report.html /setowner LOCALSERVICE`
   - This is necessary to finish setting up [agent]
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
