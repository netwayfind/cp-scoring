'use strict';

function _objectSpread(target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i] != null ? arguments[i] : {}; var ownKeys = Object.keys(source); if (typeof Object.getOwnPropertySymbols === 'function') { ownKeys = ownKeys.concat(Object.getOwnPropertySymbols(source).filter(function (sym) { return Object.getOwnPropertyDescriptor(source, sym).enumerable; })); } ownKeys.forEach(function (key) { _defineProperty(target, key, source[key]); }); } return target; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

const Plot = createPlotlyComponent(Plotly);

class App extends React.Component {
  constructor() {
    super();
    this.state = {
      authenticated: false,
      args: {}
    };
    this.authCallback = this.authCallback.bind(this);
    this.logout = this.logout.bind(this);
  }

  authCallback(statusCode) {
    if (statusCode == 200) {
      this.setState({
        authenticated: true
      });
    } else {
      this.setState({
        authenticated: false
      });
    }
  }

  logout() {
    let url = "/logout";
    fetch(url, {
      credentials: 'same-origin',
      method: "DELETE"
    }).then(function (_) {
      this.setState({
        authenticated: false
      });
    }.bind(this));
  }

  componentDidMount() {
    // check if logged in by visiting the following URL
    let url = "/templates";
    fetch(url, {
      credentials: 'same-origin'
    }).then(function (response) {
      this.authCallback(response.status);
    }.bind(this));
  }

  analysisRequestCallback(documentType, args) {
    this.setState({
      documentType: documentType,
      args: args
    });
  }

  render() {
    if (!this.state.authenticated) {
      return React.createElement("div", {
        className: "App"
      }, React.createElement(Login, {
        callback: this.authCallback
      }));
    }

    return React.createElement("div", {
      className: "App"
    }, React.createElement("div", {
      className: "heading"
    }, React.createElement("h1", null, "cp-scoring Insight")), React.createElement("div", {
      className: "navbar"
    }, React.createElement("button", {
      className: "right",
      onClick: this.logout
    }, "Logout")), React.createElement("div", {
      className: "toc"
    }, React.createElement(Analysis, {
      requestCallback: this.analysisRequestCallback.bind(this)
    })), React.createElement("div", {
      className: "content"
    }, React.createElement(AnalysisItem, {
      documentType: this.state.documentType,
      args: this.state.args
    })));
  }

}

class Analysis extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      documentType: "",
      scenarios: [],
      scenarioID: 0,
      teams: [],
      teamID: 0,
      hosts: [],
      hostname: "",
      timeStart: Date.now() * 1000,
      timeEnd: Date.now() * 1000
    };
    this.selectDocumentTypeCallback = this.selectDocumentTypeCallback.bind(this);
    this.selectScenarioCallback = this.selectScenarioCallback.bind(this);
    this.selectTeamCallback = this.selectTeamCallback.bind(this);
    this.selectHostCallback = this.selectHostCallback.bind(this);
    this.selectTimeStartCallback = this.selectTimeStartCallback.bind(this);
    this.selectTimeEndCallback = this.selectTimeEndCallback.bind(this);
    this.updateTime = this.updateTime.bind(this);
    this.submit = this.submit.bind(this);
  }

  componentDidMount() {
    this.populateSelectors();
  }

  componentWillReceiveProps(_) {
    this.populateSelectors();
  }

  selectDocumentTypeCallback(event) {
    event.preventDefault();
    this.setState({
      documentType: event.target.value
    });
  }

  selectScenarioCallback(event) {
    event.preventDefault();
    this.setState({
      scenarioID: event.target.value
    });
  }

  selectTeamCallback(event) {
    event.preventDefault();
    this.setState({
      teamID: event.target.value
    });
  }

  selectHostCallback(event) {
    event.preventDefault();
    this.setState({
      hostname: event.target.value
    });
  }

  selectTimeStartCallback(event) {
    event.preventDefault();
    let updated = this.updateTime(event, this.state.timeStart);

    if (updated === null) {
      return;
    }

    this.setState({
      timeStart: updated
    });
  }

  selectTimeEndCallback(event) {
    event.preventDefault();
    let updated = this.updateTime(event, this.state.timeEnd);

    if (updated === null) {
      return;
    }

    this.setState({
      timeEnd: updated
    });
  }

  updateTime(event, original) {
    let current = new Date(Math.trunc(original * 1000));

    if (event.target.type === "date") {
      let parts = event.target.value.split("-");

      if (parts.length != 3) {
        return null;
      }

      current.setFullYear(parts[0]); // months start counting at 0

      current.setMonth(parts[1] - 1);
      current.setDate(parts[2]);
    } else if (event.target.type === "time") {
      let parts = event.target.value.split(":");

      if (parts.length != 3) {
        return null;
      }

      current.setHours(parts[0]);
      current.setMinutes(parts[1]);
      current.setSeconds(parts[2]);
    }

    let value = Math.trunc(current.getTime() / 1000);

    if (Number.isNaN(value)) {
      return null;
    }

    return value;
  }

  populateSelectors() {
    fetch('/scenarios', {
      credentials: 'same-origin'
    }).then(function (response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }

      return response.json();
    }).then(function (data) {
      this.setState({
        scenarios: data
      });
    }.bind(this));
    fetch('/teams', {
      credentials: 'same-origin'
    }).then(function (response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }

      return response.json();
    }).then(function (data) {
      this.setState({
        teams: data
      });
    }.bind(this));
    fetch('/hosts', {
      credentials: 'same-origin'
    }).then(function (response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }

      return response.json();
    }).then(function (data) {
      this.setState({
        hosts: data
      });
    }.bind(this));
  }

  submit() {
    let args = {
      'scenario_id': this.state.scenarioID,
      'team_id': this.state.teamID,
      'hostname': this.state.hostname,
      'time_start': this.state.timeStart,
      'time_end': this.state.timeEnd
    };
    this.props.requestCallback(this.state.documentType, args);
  }

  render() {
    // form document type options
    let documentTypeOptions = [];
    documentTypeOptions.push(React.createElement("option", {
      key: "-1",
      value: ""
    }));
    documentTypeOptions.push(React.createElement("option", {
      key: "0",
      value: "reports"
    }, "Reports"));
    documentTypeOptions.push(React.createElement("option", {
      key: "1",
      value: "states"
    }, "States")); // form scenario options

    let scenarioOptions = [];
    scenarioOptions.push(React.createElement("option", {
      key: "-1",
      value: ""
    }));

    for (let i in this.state.scenarios) {
      let scenario = this.state.scenarios[i];
      scenarioOptions.push(React.createElement("option", {
        key: i,
        value: scenario.ID
      }, scenario.Name));
    } // form team options


    let teamOptions = [];
    teamOptions.push(React.createElement("option", {
      key: "-1",
      value: ""
    }));

    for (let i in this.state.teams) {
      let team = this.state.teams[i];
      teamOptions.push(React.createElement("option", {
        key: i,
        value: team.ID
      }, team.Name));
    } // form host options


    let hostOptions = [];
    hostOptions.push(React.createElement("option", {
      key: "-1",
      value: ""
    }));

    for (let i in this.state.hosts) {
      let host = this.state.hosts[i];
      hostOptions.push(React.createElement("option", {
        key: i,
        value: host.Hostname
      }, host.Hostname));
    } // form time start


    let d = new Date(this.state.timeStart * 1000);
    let startDate = ("000" + d.getFullYear()).slice(-4);
    startDate += "-";
    startDate += ("0" + (d.getMonth() + 1)).slice(-2);
    startDate += "-";
    startDate += ("0" + d.getDate()).slice(-2);
    let startTime = ("000" + d.getHours()).slice(-2);
    startTime += ":";
    startTime += ("000" + d.getMinutes()).slice(-2);
    startTime += ":";
    startTime += ("000" + d.getSeconds()).slice(-2); // form time end

    d = new Date(this.state.timeEnd * 1000);
    let endDate = ("000" + d.getFullYear()).slice(-4);
    endDate += "-";
    endDate += ("0" + (d.getMonth() + 1)).slice(-2);
    endDate += "-";
    endDate += ("0" + d.getDate()).slice(-2);
    let endTime = ("000" + d.getHours()).slice(-2);
    endTime += ":";
    endTime += ("000" + d.getMinutes()).slice(-2);
    endTime += ":";
    endTime += ("000" + d.getSeconds()).slice(-2);
    return React.createElement(React.Fragment, null, React.createElement("label", {
      name: "type"
    }, "Document Type"), React.createElement("select", {
      value: this.state.documentType,
      onChange: this.selectDocumentTypeCallback
    }, documentTypeOptions), React.createElement("br", null), React.createElement("label", {
      name: "scenarios"
    }, "Scenarios"), React.createElement("select", {
      value: this.state.scenarioID,
      onChange: this.selectScenarioCallback
    }, scenarioOptions), React.createElement("br", null), React.createElement("label", {
      name: "teams"
    }, "Teams"), React.createElement("select", {
      value: this.state.teamID,
      onChange: this.selectTeamCallback
    }, teamOptions), React.createElement("br", null), React.createElement("label", {
      name: "hosts"
    }, "Hosts"), React.createElement("select", {
      value: this.state.hostname,
      onChange: this.selectHostCallback
    }, hostOptions), React.createElement("br", null), React.createElement("label", {
      name: "timeStart"
    }, "Time start"), React.createElement("input", {
      type: "date",
      value: startDate,
      onChange: this.selectTimeStartCallback
    }), React.createElement("input", {
      type: "time",
      value: startTime,
      onChange: this.selectTimeStartCallback
    }), React.createElement("br", null), React.createElement("label", {
      name: "timeEnd"
    }, "Time end"), React.createElement("input", {
      type: "date",
      value: endDate,
      onChange: this.selectTimeEndCallback
    }), React.createElement("input", {
      type: "time",
      value: endTime,
      onChange: this.selectTimeEndCallback
    }), React.createElement("p", null), React.createElement("button", {
      onClick: this.submit
    }, "Submit"));
  }

}

class AnalysisItem extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      timeline: {},
      diffs: {},
      selected: {}
    };
    this.plotClick = this.plotClick.bind(this);
  }

  componentDidMount() {
    this.getData(this.props);
  }

  componentWillReceiveProps(newProps) {
    this.getData(newProps);
  }

  getData(props) {
    if (props.args === null || props.args === undefined) {
      return;
    }

    let urlTimeline = null;
    let urlDiffs = null;

    if (props.documentType === 'reports') {
      urlTimeline = '/analysis/reports/timeline';
      urlDiffs = '/analysis/reports/diffs';
    } else if (props.documentType === 'states') {
      urlTimeline = '/analysis/states/timeline';
      urlDiffs = '/analysis/states/diffs';
    } else {
      return;
    }

    let params = Object.entries(props.args).map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`).join('&');
    urlTimeline = urlTimeline + '?' + params;
    urlDiffs = urlDiffs + '?' + params;
    fetch(urlTimeline, {
      credentials: 'same-origin'
    }).then(function (response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }

      return response.json();
    }).then(function (data) {
      this.setState({
        timeline: data
      });
    }.bind(this));
    fetch(urlDiffs, {
      credentials: 'same-origin'
    }).then(function (response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }

      return response.json();
    }).then(function (data) {
      this.setState({
        diffs: data
      });
    }.bind(this));
  }

  plotClick(plotlyEvent) {
    // only accept left click
    if (plotlyEvent.event.buttons != 1) {
      return;
    }

    let i = plotlyEvent.points[0].pointIndex;
    let type = plotlyEvent.points[0].y;

    if (type.endsWith('(diff)')) {
      // diffs should have been previously retrieved
      this.setState({
        // choose first instance
        selected: this.state.diffs[0][i]
      });
    } else if (type.endsWith('(reports)') || type.endsWith('(states)')) {
      // get report/state
      let args = _objectSpread({}, this.props.args);

      delete args['time_start'];
      delete args['time_end'];
      args['timestamp'] = Math.trunc(plotlyEvent.points[0].data.x[i] / 1000);
      let params = Object.entries(args).map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`).join('&');
      let url = '/analysis/' + this.props.documentType + '?' + params;
      fetch(url, {
        credentials: 'same-origin'
      }).then(function (response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }

        return response.json();
      }).then(function (data) {
        // choose first instance
        this.setState({
          selected: data[0]
        });
      }.bind(this));
    }

    return false;
  }

  render() {
    let config = {
      displaylogo: false
    };
    let layout = {
      hovermode: 'closest',
      xaxis: {
        type: 'date'
      },
      yaxis: {
        autorange: 'reversed'
      }
    };
    let traces = [];
    let name = this.props.args.hostname + ' (' + this.props.documentType + ')';

    for (let i in this.state.timeline) {
      let hostInstance = this.state.timeline[i];
      let trace = {
        name: name,
        mode: 'markers',
        x: hostInstance.map(document => document.Document * 1000),
        y: hostInstance.map(_ => name)
      };
      traces.push(trace);
    } // plot diffs


    name = this.props.args.hostname + ' (diff)';

    for (let i in this.state.diffs) {
      let hostInstance = this.state.diffs[i];
      let trace = {
        name: name,
        mode: 'markers',
        x: hostInstance.map(diff => diff.Timestamp * 1000),
        y: hostInstance.map(_ => name)
      };
      traces.push(trace);
    }

    let selected = null; // diff

    if (this.state.selected.Changes != undefined) {
      let time = new Date(this.state.selected.Timestamp * 1000).toLocaleString();
      let changes = [];

      for (let i in this.state.selected.Changes) {
        let change = this.state.selected.Changes[i];
        changes.push(React.createElement("li", {
          key: i
        }, change.Type, " - ", change.Key, " - ", change.Item));
      }

      selected = React.createElement(React.Fragment, null, "Time: ", time, React.createElement("br", null), "Changes:", React.createElement("ul", null, changes));
    } // report
    else if (this.props.documentType === 'reports') {
        let time = new Date(this.state.selected.Timestamp * 1000).toLocaleString();
        let findings = [];

        for (let i in this.state.selected.Findings) {
          let finding = this.state.selected.Findings[i];
          findings.push(React.createElement("li", {
            key: i
          }, finding.Show, " - ", finding.Value, " - ", finding.Message));
        }

        selected = React.createElement(React.Fragment, null, "Time: ", time, React.createElement("br", null), "Findings:", React.createElement("ul", null, findings));
      } // state
      else if (this.props.documentType === 'states') {
          let time = new Date(this.state.selected.Timestamp * 1000).toLocaleString();
          let errors = [];

          for (let i in this.state.selected.Errors) {
            let error = this.state.selected.Errors[i];
            errors.push(React.createElement("li", {
              key: i
            }, error));
          }

          let users = [];

          for (let i in this.state.selected.Users) {
            let user = this.state.selected.Users[i];
            let passwordLastSet = new Date(user.PasswordLastSet * 1000).toLocaleString();
            users.push(React.createElement("li", {
              key: i
            }, user.ID, " - ", user.Name, " - ", user.AccountActive, " - ", user.AccountExpires, " - ", passwordLastSet, " - ", user.PasswordExpires));
          }

          let groups = [];

          for (let group in this.state.selected.Groups) {
            let members = this.state.selected.Groups[group];

            if (members.length === 0) {
              groups.push(React.createElement("li", {
                key: group
              }, group));
            } else {
              let membersStr = members.map(member => member.Name).join(', ');
              groups.push(React.createElement("li", {
                key: group
              }, group, " - [", membersStr, "]"));
            }
          }

          let software = [];

          for (let i in this.state.selected.Software) {
            let sw = this.state.selected.Software[i];
            software.push(React.createElement("li", {
              key: i
            }, sw.Name, " - ", sw.Version));
          }

          let processes = [];

          for (let i in this.state.selected.Processes) {
            let process = this.state.selected.Processes[i];
            processes.push(React.createElement("li", {
              key: i
            }, process.PID, " - ", process.User, " - ", process.CommandLine));
          }

          let conns = [];

          for (let i in this.state.selected.NetworkConnections) {
            let conn = this.state.selected.NetworkConnections[i];
            conns.push(React.createElement("li", {
              key: i
            }, conn.Protocol, " - ", conn.LocalAddress, ":", conn.LocalPort, " - ", conn.RemoteAddress, ":", conn.RemotePort, " - ", conn.State));
          }

          selected = React.createElement(React.Fragment, null, "Time: ", time, React.createElement("br", null), "OS: ", this.state.selected.OS, React.createElement("br", null), "Hostname: ", this.state.selected.Hostname, React.createElement("br", null), "Errors:", React.createElement("ul", null, errors), React.createElement("br", null), "Users:", React.createElement("ul", null, users), React.createElement("br", null), "Groups:", React.createElement("ul", null, groups), React.createElement("br", null), "Software:", React.createElement("ul", null, software), React.createElement("br", null), "Processes:", React.createElement("ul", null, processes), React.createElement("br", null), "Network connections:", React.createElement("ul", null, conns), React.createElement("br", null));
        }

    return React.createElement(React.Fragment, null, "Timeline", React.createElement("ul", null, React.createElement(Plot, {
      data: traces,
      layout: layout,
      config: config,
      onClick: this.plotClick
    })), React.createElement("p", null), "Selected", React.createElement("ul", null, selected));
  }

}

ReactDOM.render(React.createElement(App, null), document.getElementById('app'));