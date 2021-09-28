import "../App.css";

import { Component, Fragment } from "react";
import { withRouter } from "react-router-dom";

const ACTION_PRESET = Object.freeze({
  EXEC: "EXEC",
  SH: "sh",
  POWERSHELL: "powershell",
});

const ACTION_PRESET_CHECK = Object.freeze({
  FILE_PERMISSIONS_LINUX: "file permissions (linux)",
  FIREWALL_INBOUND_ALLOW_IPTABLES: "firewall inbound allow (iptables)",
  FIREWALL_INBOUND_ALLOW_UFW: "firewall inbound allow (ufw)",
  FIREWALL_INBOUND_DEFAULT_DROP_LINUX: "firewall inbound default drop (linux)",
  FIREWALL_FORWARD_DEFAULT_DROP_LINUX: "firewall forward default drop (linux)",
  FIREWALL_WINDOWS_DOMAIN: "Firewall (Windows - Domain) set",
  FIREWALL_WINDOWS_PRIVATE: "Firewall (Windows - Private) set",
  FIREWALL_WINDOWS_PUBLIC: "Firewall (Windows - Public) set",
  NETWORK_SERVICE_NOT_AVAILABLE_LINUX: "network service not available (linux)",
  SOFTWARE_INSTALLED_LINUX: "software installed (linux)",
  SOFTWARE_PACKAGES_UPDATED_LINUX: "software packages updated (linux)",
  SOFTWARE_REMOVED_LINUX: "software removed (linux)",
  TMP_APT_PACKAGE_LIST: "TEMP: apt package list",
  TMP_APT_PACKAGE_LIST_REMOVE: "TEMP: apt package list remove",
  USER_ADDED_LINUX: "user added (linux)",
  USER_ADDED_WINDOWS: "user added (windows)",
  USER_ADDED_TO_GROUP_LINUX: "user added to group (linux)",
  USER_ADDED_TO_GROUP_WINDOWS: "user added to group (windows)",
  USER_DISABLED_WINDOWS: "user disabled (windows)",
  USER_DISABLED_ADMINISTRATOR_WINDOWS: "user Administrator disabled (windows)",
  USER_DISABLED_GUEST_WINDOWS: "user Guest disabled (windows)",
  USER_RENAMED_ADMINISTRATOR_WINDOWS: "user Administrator renamed (windows)",
  USER_RENAMED_GUEST_WINDOWS: "user Guest renamed (windows)",
  USER_PASSWORD_CHANGED_LINUX: "user password changed (linux)",
  USER_PASSWORD_CHANGED_WINDOWS: "user password changed (windows)",
  USER_REMOVED_LINUX: "user removed (linux)",
  USER_REMOVED_WINDOWS: "user removed (windows)",
  USER_REMOVED_FROM_GROUP_LINUX: "user removed from group (linux)",
  USER_REMOVED_FROM_GROUP_WINDOWS: "user removed from group (windows)",
});

const ACTION_PRESET_CONFIG = Object.freeze({
  FIREWALL_OFF_WINDOWS: "firewall off (windows)",
  INSTALL_CHOCO: "install chocolatey",
  INSTALL_CHOCO_POST: "install chocolatey (post)",
  INSTALL_PACKAGES_LINUX: "install packages (linux)",
  INSTALL_SOFTWARE_CHOCO: "install software (choco)",
  NET_SHARE_ADD: "net share add",
  NEW_DIR_LINUX: "new directory (linux)",
  NEW_DIR_WINDOWS: "new directory (windows)",
  USER_ADD_LINUX: "user add (linux)",
  USER_ADD_LINUX_SYSTEM: "user add (linux - system)",
  USER_ADD_WINDOWS: "user add (windows)",
});

const CHECK_TYPE = Object.freeze({
  EXEC: "EXEC",
  FILE_CONTAINS: "FILE_CONTAINS",
  FILE_EXIST: "FILE_EXIST",
  FILE_REGEX: "FILE_REGEX",
  FILE_VALUE: "FILE_VALUE",
});

const COMMAND = Object.freeze({
  CHOCO: "C:\\ProgramData\\chocolatey\\bin\\choco.exe",
  CMD: "C:\\Windows\\System32\\cmd.exe",
  NET: "C:\\Windows\\System32\\net.exe",
  POWERSHELL: "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
  SH: "/bin/sh",
});

const OPERATOR = Object.freeze({
  CONTAINS: "CONTAINS",
  EQUAL: "EQUAL",
  NOT_CONTAINS: "NOT_CONTAINS",
  NOT_EQUAL: "NOT_EQUAL",
});

class ScenarioHost extends Component {
  constructor(props) {
    super(props);
    this.state = {
      answers: props.answers,
      checks: props.checks,
      config: props.config,
      hostname: props.hostname,
      presetAddCheck: ACTION_PRESET.EXEC,
      presetAddConfig: ACTION_PRESET.EXEC,
    };

    this.handleAnswerUpdate = this.handleAnswerUpdate.bind(this);
    this.handleCheckAdd = this.handleCheckAdd.bind(this);
    this.handleCheckDelete = this.handleCheckDelete.bind(this);
    this.handleCheckReorder = this.handleCheckReorder.bind(this);
    this.handleCheckUpdate = this.handleCheckUpdate.bind(this);
    this.handleCheckArgAdd = this.handleCheckArgAdd.bind(this);
    this.handleCheckArgDelete = this.handleCheckArgDelete.bind(this);
    this.handleCheckArgUpdate = this.handleCheckArgUpdate.bind(this);
    this.handleConfigAdd = this.handleConfigAdd.bind(this);
    this.handleConfigDelete = this.handleConfigDelete.bind(this);
    this.handleConfigUpdate = this.handleConfigUpdate.bind(this);
    this.handleConfigReorder = this.handleConfigReorder.bind(this);
    this.handleConfigArgAdd = this.handleConfigArgAdd.bind(this);
    this.handleConfigArgDelete = this.handleConfigArgDelete.bind(this);
    this.handleConfigArgUpdate = this.handleConfigArgUpdate.bind(this);
    this.handleSave = this.handleSave.bind(this);
    this.handleUpdatePresetAddCheck = this.handleUpdatePresetAddCheck.bind(
      this
    );
    this.handleUpdatePresetAddConfig = this.handleUpdatePresetAddConfig.bind(
      this
    );
  }

  componentDidUpdate(prevProps) {
    if (this.props.hostname !== prevProps.hostname) {
      this.setState({
        answers: this.props.answers,
        checks: this.props.checks,
        config: this.props.config,
        hostname: this.props.hostname,
      });
    }
  }

  handleAnswerUpdate(i, event) {
    let name = event.target.name;
    let value = event.target.value;
    let answers = [...this.state.answers];
    if (event.target.type === "number") {
      value = Number(value);
    }
    answers[i][name] = value;
    this.setState({
      answers: answers,
    });
  }

  handleCheckAdd() {
    let answers = [...this.state.answers];
    let checks = [...this.state.checks];
    let preset = this.preset(this.state.presetAddCheck);
    answers.push({
      Operator: preset.Operator,
      Value: preset.Value,
      Points: preset.Points,
    });
    checks.push({
      Description: preset.Description,
      Type: preset.Type,
      Command: preset.Command,
      Args: preset.Args,
    });
    this.setState({
      answers: answers,
      checks: checks,
      presetAddCheck: ACTION_PRESET.EXEC,
    });
  }

  handleCheckDelete(i) {
    let answers = [...this.state.answers];
    let checks = [...this.state.checks];
    answers.splice(i, 1);
    checks.splice(i, 1);
    this.setState({
      answers: answers,
      checks: checks,
    });
  }

  handleCheckReorder(event, currentIndex) {
    let newIndex = event.target.value;
    let answers = [...this.state.answers];
    let checks = [...this.state.checks];
    answers.splice(newIndex, 0, answers.splice(currentIndex, 1)[0]);
    checks.splice(newIndex, 0, checks.splice(currentIndex, 1)[0]);
    this.setState({
      answers: answers,
      checks: checks,
    });
  }

  handleCheckUpdate(i, event) {
    let name = event.target.name;
    let value = event.target.value;
    let checks = [...this.state.checks];
    checks[i][name] = value;
    this.setState({
      checks: checks,
    });
  }

  handleCheckArgAdd(i) {
    let checks = [...this.state.checks];
    checks[i]["Args"].push("");
    this.setState({
      checks: checks,
    });
  }

  handleCheckArgDelete(i, j) {
    let checks = [...this.state.checks];
    checks[i]["Args"].splice(j, 1);
    this.setState({
      checks: checks,
    });
  }

  handleCheckArgUpdate(i, j, event) {
    let value = event.target.value;
    let checks = [...this.state.checks];
    checks[i]["Args"][j] = value;
    this.setState({
      checks: checks,
    });
  }

  handleConfigAdd() {
    let config = [...this.state.config];
    let preset = this.preset(this.state.presetAddConfig);
    config.push({
      Description: preset.Description,
      Type: preset.Type,
      Command: preset.Command,
      Args: preset.Args,
    });
    this.setState({
      config: config,
      presetAddConfig: ACTION_PRESET.EXEC,
    });
  }

  handleConfigDelete(i) {
    let config = [...this.state.config];
    config.splice(i, 1);
    this.setState({
      config: config,
    });
  }

  handleConfigReorder(event, currentIndex) {
    let newIndex = event.target.value;
    let config = [...this.state.config];
    config.splice(newIndex, 0, config.splice(currentIndex, 1)[0]);
    this.setState({
      config: config,
    });
  }

  handleConfigUpdate(i, event) {
    let name = event.target.name;
    let value = event.target.value;
    let config = [...this.state.config];
    config[i][name] = value;
    this.setState({
      config: config,
    });
  }

  handleConfigArgAdd(i) {
    let config = [...this.state.config];
    config[i]["Args"].push("");
    this.setState({
      config: config,
    });
  }

  handleConfigArgDelete(i, j) {
    let config = [...this.state.config];
    config[i]["Args"].splice(j, 1);
    this.setState({
      config: config,
    });
  }

  handleConfigArgUpdate(i, j, event) {
    let value = event.target.value;
    let config = [...this.state.config];
    config[i]["Args"][j] = value;
    this.setState({
      config: config,
    });
  }

  handleSave(event) {
    if (event !== null) {
      event.preventDefault();
    }
    this.props.parentCallback(
      this.state.checks,
      this.state.answers,
      this.state.config
    );
  }

  handleUpdatePresetAddCheck(event) {
    let value = event.target.value;
    this.setState({
      presetAddCheck: value,
    });
  }

  handleUpdatePresetAddConfig(event) {
    let value = event.target.value;
    this.setState({
      presetAddConfig: value,
    });
  }

  preset(p) {
    let description = p;
    let type = CHECK_TYPE.EXEC;
    let command = "";
    let args = [];
    let operator = "";
    let value = "";
    let points = 0;
    if (p === ACTION_PRESET.EXEC) {
      // default
    } else if (p === ACTION_PRESET.SH) {
      command = COMMAND.SH;
      args = ["-c", ""];
    } else if (p === ACTION_PRESET.POWERSHELL) {
      command = COMMAND.POWERSHELL;
      args = ["-command", ""];
    } else if (p === ACTION_PRESET_CHECK.FILE_PERMISSIONS_LINUX) {
      command = COMMAND.SH;
      args = ["-c", "stat -c %a file"];
      operator = OPERATOR.EQUAL;
      value = "###";
    } else if (p === ACTION_PRESET_CHECK.FIREWALL_INBOUND_DEFAULT_DROP_LINUX) {
      command = COMMAND.SH;
      args = ["-c", "iptables -nvL INPUT | grep -q 'policy DROP' ; echo $?"];
      operator = OPERATOR.EQUAL;
      value = "0";
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.FIREWALL_FORWARD_DEFAULT_DROP_LINUX) {
      command = COMMAND.SH;
      args = ["-c", "iptables -nvL FORWARD | grep -q 'policy DROP' ; echo $?"];
      operator = OPERATOR.EQUAL;
      value = "0";
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.FIREWALL_INBOUND_ALLOW_IPTABLES) {
      command = COMMAND.SH;
      args = [
        "-c",
        "iptables -nvL INPUT | grep 'protocol dpt:port' | grep -q ACCEPT ; echo $?",
      ];
      operator = OPERATOR.EQUAL;
      value = "0";
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.FIREWALL_INBOUND_ALLOW_UFW) {
      command = COMMAND.SH;
      args = [
        "-c",
        "iptables -nvL ufw-user-input | grep 'protocol dpt:port' | grep -q ACCEPT ; echo $?",
      ];
      operator = OPERATOR.EQUAL;
      value = "0";
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.FIREWALL_WINDOWS_DOMAIN) {
      command = COMMAND.POWERSHELL;
      args = [
        "-command",
        "Get-NetFirewallProfile -Name Domain | Select-Object Enabled | ConvertTo-Json -Compress",
      ];
      operator = OPERATOR.EQUAL;
      value = "{\"Enabled\":1}";
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.FIREWALL_WINDOWS_PUBLIC) {
      command = COMMAND.POWERSHELL;
      args = [
        "-command",
        "Get-NetFirewallProfile -Name Public | Select-Object Enabled | ConvertTo-Json -Compress",
      ];
      operator = OPERATOR.EQUAL;
      value = "{\"Enabled\":1}";
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.FIREWALL_WINDOWS_PRIVATE) {
      command = COMMAND.POWERSHELL;
      args = [
        "-command",
        "Get-NetFirewallProfile -Name Private | Select-Object Enabled | ConvertTo-Json -Compress",
      ];
      operator = OPERATOR.EQUAL;
      value = "{\"Enabled\":1}";
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.NETWORK_SERVICE_NOT_AVAILABLE_LINUX) {
      command = COMMAND.SH;
      args = ["-c", "ss -ntlp | grep ':port ' | grep -q service ; echo $?"];
      operator = OPERATOR.NOT_EQUAL;
      value = "0";
      points = -1;
    } else if (p === ACTION_PRESET_CHECK.SOFTWARE_INSTALLED_LINUX) {
      command = COMMAND.SH;
      args = ["-c", "grep -q '^software/' apt; echo $?"];
      operator = OPERATOR.EQUAL;
      value = "0";
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.SOFTWARE_PACKAGES_UPDATED_LINUX) {
      command = COMMAND.SH;
      args = ["-c", "apt list --upgradable | wc -l"];
      operator = OPERATOR.EQUAL;
      value = "1";
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.SOFTWARE_REMOVED_LINUX) {
      command = COMMAND.SH;
      args = ["-c", "grep -q '^software/' apt; echo $?"];
      operator = OPERATOR.NOT_EQUAL;
      value = "0";
    } else if (p === ACTION_PRESET_CHECK.TMP_APT_PACKAGE_LIST) {
      command = COMMAND.SH;
      args = ["-c", "apt list --installed > apt"];
    } else if (p === ACTION_PRESET_CHECK.TMP_APT_PACKAGE_LIST_REMOVE) {
      command = COMMAND.SH;
      args = ["-c", "rm apt"];
    } else if (p === ACTION_PRESET_CHECK.USER_ADDED_LINUX) {
      command = COMMAND.SH;
      args = ["-c", "grep -q '^user:' /etc/passwd; echo $?"];
      operator = OPERATOR.EQUAL;
      value = "0";
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.USER_ADDED_WINDOWS) {
      command = COMMAND.POWERSHELL;
      args = [
        "-command",
        'Get-LocalUser | Where-Object {$_.Name -eq "user"} | Measure-Object | Select-Object Count | ConvertTo-Json -Compress',
      ];
      operator = OPERATOR.EQUAL;
      value = '{"Count":1}';
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.USER_ADDED_TO_GROUP_LINUX) {
      command = COMMAND.SH;
      args = ["-c", "grep '^group:' /etc/group | grep -q user ; echo $?"];
      operator = OPERATOR.EQUAL;
      value = "0";
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.USER_ADDED_TO_GROUP_WINDOWS) {
      command = COMMAND.POWERSHELL;
      args = [
        "-command",
        'Get-LocalGroupMember -Name group | Where-Object {$_.Name -eq "hostname\\user"} | Measure-Object | Select-Object Count | ConvertTo-Json -Compress',
      ];
      operator = OPERATOR.EQUAL;
      value = '{"Count":1}';
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.USER_DISABLED_WINDOWS) {
      command = COMMAND.POWERSHELL;
      args = [
        "-command",
        'Get-LocalUser | Where-Object {($_.Name -eq "user") -and ($_.Enabled -eq $false)} | Measure-Object | Select-Object Count | ConvertTo-Json -Compress',
      ];
      operator = OPERATOR.EQUAL;
      value = '{"Count":1}';
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.USER_DISABLED_ADMINISTRATOR_WINDOWS) {
      command = COMMAND.POWERSHELL;
      args = [
        "-command",
        'Get-LocalUser | Where-Object {($_.SID -like "*-500") -and ($_.Enabled -eq $false)} | Measure-Object | Select-Object Count | ConvertTo-Json -Compress',
      ];
      operator = OPERATOR.EQUAL;
      value = '{"Count":1}';
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.USER_DISABLED_GUEST_WINDOWS) {
      command = COMMAND.POWERSHELL;
      args = [
        "-command",
        'Get-LocalUser | Where-Object {($_.SID -like "*-501") -and ($_.Enabled -eq $false)} | Measure-Object | Select-Object Count | ConvertTo-Json -Compress',
      ];
      operator = OPERATOR.EQUAL;
      value = '{"Count":1}';
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.USER_RENAMED_ADMINISTRATOR_WINDOWS) {
      command = COMMAND.POWERSHELL;
      args = [
        "-command",
        'Get-LocalUser | Where-Object {($_.SID -like "*-500") -and ($_.Name -eq "Administrator")} | Measure-Object | Select-Object Count | ConvertTo-Json -Compress',
      ];
      operator = OPERATOR.EQUAL;
      value = '{"Count":1}';
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.USER_RENAMED_GUEST_WINDOWS) {
      command = COMMAND.POWERSHELL;
      args = [
        "-command",
        'Get-LocalUser | Where-Object {($_.SID -like "*-501") -and ($_.Name -eq "Guest")} | Measure-Object | Select-Object Count | ConvertTo-Json -Compress',
      ];
      operator = OPERATOR.EQUAL;
      value = '{"Count":0}';
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.USER_PASSWORD_CHANGED_LINUX) {
      command = COMMAND.SH;
      args = [
        "-c",
        "grep '^user' /etc/shadow | grep -q 'passwordHash' ; echo $?",
      ];
      operator = OPERATOR.EQUAL;
      value = "1";
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.USER_PASSWORD_CHANGED_WINDOWS) {
      command = COMMAND.CMD;
      args = [
        "/C",
        "net user user",
      ];
      operator = OPERATOR.NOT_CONTAINS;
      value = "Password last set            1/11/2021 12:00:00 PM";
      points = 1;
    } else if (p === ACTION_PRESET_CHECK.USER_REMOVED_LINUX) {
      command = COMMAND.SH;
      args = ["-c", "grep -q '^user:' /etc/passwd; echo $?"];
      operator = OPERATOR.NOT_EQUAL;
      value = "0";
    } else if (p === ACTION_PRESET_CHECK.USER_REMOVED_WINDOWS) {
      command = COMMAND.POWERSHELL;
      args = [
        "-command",
        'Get-LocalUser | Where-Object {$_.Name -eq "user"} | Measure-Object | Select-Object Count | ConvertTo-Json -Compress',
      ];
      operator = OPERATOR.EQUAL;
      value = '{"Count":0}';
    } else if (p === ACTION_PRESET_CHECK.USER_REMOVED_FROM_GROUP_LINUX) {
      command = COMMAND.SH;
      args = ["-c", "grep '^group:' /etc/group | grep -q user ; echo $?"];
      operator = OPERATOR.NOT_EQUAL;
      value = "0";
    } else if (p === ACTION_PRESET_CHECK.USER_REMOVED_FROM_GROUP_WINDOWS) {
      command = COMMAND.POWERSHELL;
      args = [
        "-command",
        'Get-LocalGroupMember -Name group | Where-Object {$_.Name -eq "hostname\\user"} | Measure-Object | Select-Object Count | ConvertTo-Json -Compress',
      ];
      operator = OPERATOR.EQUAL;
      value = '{"Count":0}';
    } else if (p === ACTION_PRESET_CONFIG.FIREWALL_OFF_WINDOWS) {
      command = COMMAND.CMD;
      args = ["/C", "netsh advfirewall set allprofiles state off"];
    } else if (p === ACTION_PRESET_CONFIG.INSTALL_CHOCO) {
      command = COMMAND.POWERSHELL;
      args = [
        "-command",
        "Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))",
      ];
    } else if (p === ACTION_PRESET_CONFIG.INSTALL_CHOCO_POST) {
      command = COMMAND.POWERSHELL;
      args = [
        "feature",
        "enable",
        "--name",
        "allowGlobalConfirmation",
      ];
    }else if (p === ACTION_PRESET_CONFIG.INSTALL_PACKAGES_LINUX) {
      command = COMMAND.SH;
      args = ["-c", "apt-get update && apt-get -q install -y packages"];
    } else if (p === ACTION_PRESET_CONFIG.INSTALL_SOFTWARE_CHOCO) {
      command = COMMAND.CHOCO;
      args = ["install", "software"];
    } else if (p === ACTION_PRESET_CONFIG.NET_SHARE_ADD) {
      command = COMMAND.NET;
      args = ["share", "share=path"];
    } else if (p === ACTION_PRESET_CONFIG.NEW_DIR_LINUX) {
      command = COMMAND.SH;
      args = ["-c", "mkdir directory"];
    } else if (p === ACTION_PRESET_CONFIG.NEW_DIR_WINDOWS) {
      command = COMMAND.CMD;
      args = ["/C", "mkdir directory"];
    } else if (p === ACTION_PRESET_CONFIG.USER_ADD_LINUX) {
      command = COMMAND.SH;
      args = [
        "-c",
        "useradd -m -s /bin/bash username && echo username:password | chpasswd",
      ];
    } else if (p === ACTION_PRESET_CONFIG.USER_ADD_LINUX_SYSTEM) {
      command = COMMAND.SH;
      args = [
        "-c",
        "useradd -r -s /bin/bash username && echo username:password | chpasswd",
      ];
    } else if (p === ACTION_PRESET_CONFIG.USER_ADD_WINDOWS) {
      command = COMMAND.NET;
      args = ["user", "username", "password", "/add"];
    } else {
      description = "unsupported preset";
    }
    return {
      Description: description,
      Type: type,
      Command: command,
      Args: args,
      Operator: operator,
      Value: value,
      Points: points,
    };
  }

  render() {
    let actionOptions = [];
    for (let type in CHECK_TYPE) {
      let value = CHECK_TYPE[type];
      actionOptions.push(<option key={type}>{value}</option>);
    }
    let operatorOptions = [];
    operatorOptions.push(<option key="" />);
    for (let operator in OPERATOR) {
      let value = OPERATOR[operator];
      operatorOptions.push(<option key={operator}>{value}</option>);
    }

    let presetAddCheckOptions = [];
    for (let preset in ACTION_PRESET) {
      let value = ACTION_PRESET[preset];
      presetAddCheckOptions.push(<option key={preset}>{value}</option>);
    }
    for (let preset in ACTION_PRESET_CHECK) {
      let value = ACTION_PRESET_CHECK[preset];
      presetAddCheckOptions.push(<option key={preset}>{value}</option>);
    }

    let presetAddConfigOptions = [];
    for (let preset in ACTION_PRESET) {
      let value = ACTION_PRESET[preset];
      presetAddConfigOptions.push(<option key={preset}>{value}</option>);
    }
    for (let preset in ACTION_PRESET_CONFIG) {
      let value = ACTION_PRESET_CONFIG[preset];
      presetAddConfigOptions.push(<option key={preset}>{value}</option>);
    }

    let checkList = [];
    let checks = this.state.checks;
    let checksPositionOptions = [];
    for (let i in checks) {
      checksPositionOptions.push(
        <option key={i} value={i}>
          {Number(i) + 1}
        </option>
      );
    }
    checks.forEach((check, i) => {
      let args = [];
      if (check.Args) {
        check.Args.forEach((arg, j) => {
          args.push(
            <li key={j}>
              <input
                className="input-50"
                onChange={(event) => this.handleCheckArgUpdate(i, j, event)}
                value={arg}
              ></input>
              <button
                type="button"
                onClick={() => this.handleCheckArgDelete(i, j)}
              >
                -
              </button>
            </li>
          );
        });
      }
      args.push(
        <li key="arg_add">
          <button type="button" onClick={() => this.handleCheckArgAdd(i)}>
            Add Arg
          </button>
        </li>
      );
      let answer = this.state.answers[i];
      checkList.push(
        <li key={i}>
          <select
            onChange={(event) => this.handleCheckReorder(event, i)}
            value={i}
          >
            {checksPositionOptions}
          </select>
          <details>
            <summary>{check.Description}</summary>
            <button type="button" onClick={() => this.handleCheckDelete(i)}>
              Delete Check
            </button>
            <p />
            <label htmlFor="Description">Description</label>
            <input
              className="input-20"
              name="Description"
              onChange={(event) => this.handleCheckUpdate(i, event)}
              value={check.Description}
            />
            <br />
            <label htmlFor="Type">Type</label>
            <select
              name="Type"
              onChange={(event) => this.handleCheckUpdate(i, event)}
              value={check.Type}
            >
              {actionOptions}
            </select>
            <br />
            {check.Type === CHECK_TYPE.EXEC ? (
              <Fragment>
                <label htmlFor="Command">Command</label>
                <input
                  className="input-50"
                  name="Command"
                  onChange={(event) => this.handleCheckUpdate(i, event)}
                  value={check.Command}
                />
                <br />
              </Fragment>
            ) : null}
            <label htmlFor="Args">Args</label>
            <ul>{args}</ul>
            <label htmlFor="Answer">Answer</label>
            <select
              name="Operator"
              onChange={(event) => this.handleAnswerUpdate(i, event)}
              value={answer.Operator}
            >
              {operatorOptions}
            </select>
            <input
              className="input-50"
              name="Value"
              onChange={(event) => this.handleAnswerUpdate(i, event)}
              value={answer.Value}
            />
            <br />
            <label htmlFor="Points">Points</label>
            <input
              className="input-5"
              name="Points"
              onChange={(event) => this.handleAnswerUpdate(i, event)}
              value={answer.Points}
              type="number"
              steps="1"
            />
          </details>
        </li>
      );
    });
    checkList.push(
      <li key="check_add">
        <select
          onChange={this.handleUpdatePresetAddCheck}
          value={this.state.presetAddCheck}
        >
          {presetAddCheckOptions}
        </select>
        <button type="button" onClick={this.handleCheckAdd}>
          Add Check
        </button>
      </li>
    );

    let configList = [];
    let config = this.state.config;
    let configPositionOptions = [];
    for (let i in config) {
      configPositionOptions.push(
        <option key={i} value={i}>
          {Number(i) + 1}
        </option>
      );
    }
    config.forEach((conf, i) => {
      let args = [];
      if (conf.Args) {
        conf.Args.forEach((arg, j) => {
          args.push(
            <li key={j}>
              <input
                className="input-50"
                onChange={(event) => this.handleConfigArgUpdate(i, j, event)}
                value={arg}
              ></input>
              <button
                type="button"
                onClick={() => this.handleConfigArgDelete(i, j)}
              >
                -
              </button>
            </li>
          );
        });
      }
      args.push(
        <li key="arg_add">
          <button type="button" onClick={() => this.handleConfigArgAdd(i)}>
            Add Arg
          </button>
        </li>
      );
      configList.push(
        <li key={i}>
          <select
            onChange={(event) => this.handleConfigReorder(event, i)}
            value={i}
          >
            {configPositionOptions}
          </select>
          <details>
            <summary>{conf.Description}</summary>
            <button type="button" onClick={() => this.handleConfigDelete(i)}>
              Delete Config
            </button>
            <p />
            <label htmlFor="Description">Description</label>
            <input
              className="input-20"
              name="Description"
              onChange={(event) => this.handleConfigUpdate(i, event)}
              value={conf.Description}
            />
            <br />
            <label htmlFor="Type">Type</label>
            <select disabled name="Type" value={CHECK_TYPE.EXEC}>
              {actionOptions}
            </select>
            <br />
            <label htmlFor="Command">Command</label>
            <input
              className="input-50"
              name="Command"
              onChange={(event) => this.handleConfigUpdate(i, event)}
              value={conf.Command}
            />
            <br />
            <label htmlFor="Args">Args</label>
            <ul>{args}</ul>
          </details>
        </li>
      );
    });
    configList.push(
      <li key="config_add">
        <select
          onChange={this.handleUpdatePresetAddConfig}
          value={this.state.presetAddConfig}
        >
          {presetAddConfigOptions}
        </select>
        <button type="button" onClick={this.handleConfigAdd}>
          Add Config
        </button>
      </li>
    );

    return (
      <form onSubmit={this.handleSave}>
        <p>Checks</p>
        <ol>{checkList}</ol>
        <p>Config</p>
        <ol>{configList}</ol>
        <button type="submit">Save Host</button>
      </form>
    );
  }
}

export default withRouter(ScenarioHost);
