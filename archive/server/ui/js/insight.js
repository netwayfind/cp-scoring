'use strict';

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
    this.analysisRequestCallback = this.analysisRequestCallback.bind(this);
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

  analysisRequestCallback(args) {
    this.setState({
      args: args
    });
  }

  render() {
    if (!this.state.authenticated) {
      return /*#__PURE__*/React.createElement("div", {
        className: "App"
      }, /*#__PURE__*/React.createElement(Login, {
        callback: this.authCallback
      }));
    }

    return /*#__PURE__*/React.createElement("div", {
      className: "App"
    }, /*#__PURE__*/React.createElement("div", {
      className: "heading"
    }, /*#__PURE__*/React.createElement("h1", null, "cp-scoring Insight")), /*#__PURE__*/React.createElement("div", {
      className: "navbar"
    }, /*#__PURE__*/React.createElement("button", {
      className: "right",
      onClick: this.logout
    }, "Logout")), /*#__PURE__*/React.createElement("div", {
      className: "toc"
    }, /*#__PURE__*/React.createElement(AnalysisConfig, {
      requestCallback: this.analysisRequestCallback
    })), /*#__PURE__*/React.createElement("div", {
      className: "content"
    }, /*#__PURE__*/React.createElement(AnalysisResults, {
      args: this.state.args,
      selectedCallback: this.analysisSelectedCallback
    })));
  }

}

class AnalysisConfig extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      scenarios: [],
      scenarioID: null,
      teams: [],
      teamID: null,
      timeStart: Date.now(),
      timeEnd: Date.now()
    };
    this.selectScenarioCallback = this.selectScenarioCallback.bind(this);
    this.selectTeamCallback = this.selectTeamCallback.bind(this);
    this.selectTimeStartCallback = this.selectTimeStartCallback.bind(this);
    this.selectTimeEndCallback = this.selectTimeEndCallback.bind(this);
    this.updateTime = this.updateTime.bind(this);
    this.submit = this.submit.bind(this);
  }

  componentDidMount() {
    this.populateSelectors();
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
    let current = new Date(Math.trunc(original));

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

    let value = Math.trunc(current.getTime());

    if (Number.isNaN(value)) {
      return null;
    }

    return value;
  }

  populateSelectors() {
    fetch('/scenarios', {
      credentials: 'same-origin'
    }).then(async function (response) {
      if (response.status === 200) {
        let data = await response.json();
        return {
          error: null,
          scenarios: data
        };
      }

      let text = await response.text();
      return {
        error: text
      };
    }).then(function (s) {
      this.setState(s);
    }.bind(this));
    fetch('/teams', {
      credentials: 'same-origin'
    }).then(async function (response) {
      if (response.status === 200) {
        let data = await response.json();
        return {
          error: null,
          teams: data
        };
      }

      let text = await response.text();
      return {
        error: text
      };
    }).then(function (s) {
      this.setState(s);
    }.bind(this));
  }

  submit() {
    let args = {
      'scenario_id': this.state.scenarioID,
      'team_id': this.state.teamID,
      'time_start': Math.trunc(this.state.timeStart / 1000),
      'time_end': Math.trunc(this.state.timeEnd / 1000)
    };
    let params = Object.entries(args).map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`).join('&');
    let requestReports = fetch('/analysis/reports/timeline?' + params, {
      credentials: 'same-origin'
    });
    let requestReportDiffs = fetch('/analysis/reports/diffs?' + params, {
      credentials: 'same-origin'
    });
    let requestStates = fetch('/analysis/states/timeline?' + params, {
      credentials: 'same-origin'
    });
    let requestStateDiffs = fetch('/analysis/states/diffs?' + params, {
      credentials: 'same-origin'
    });
    let requestScores = fetch('/reports/scenario/' + this.state.scenarioID + '/timeline?hostname=*&' + params, {
      credentials: 'same-origin'
    });
    this.setState({
      error: "Running query..."
    });
    this.props.requestCallback({
      reportTimeline: null,
      reportDiffs: null,
      stateTimeline: null,
      stateDiffs: null,
      scores: null
    });
    Promise.all([requestReports, requestReportDiffs, requestStates, requestStateDiffs, requestScores]).then(async function (responses) {
      let j = [];
      let errors = [];

      for (let r in responses) {
        let response = responses[r];

        if (response.status >= 400) {
          errors.push(await response.text());
        }

        j.push(await response.json());
      }

      if (errors.length === 0) {
        this.setState({
          error: null
        });
      } else {
        this.setState({
          error: errors.join(", ")
        });
      }

      return Promise.all(j);
    }.bind(this)).then(function (data) {
      let error = null;

      if (data != undefined && data != null && data.length > 0) {
        let emptyData = true;

        for (let i = 0; i < data.length; i++) {
          let entry = data[i];

          if (entry != undefined && entry != null && Object.keys(entry).length > 0) {
            emptyData = false;
            break;
          }
        }

        if (emptyData) {
          error = "No data found";
        }
      }

      this.setState({
        error: error
      });
      this.props.requestCallback({
        reportTimeline: data[0],
        reportDiffs: data[1],
        stateTimeline: data[2],
        stateDiffs: data[3],
        scores: data[4]
      });
    }.bind(this));
  }

  render() {
    // form scenario options
    let scenarioOptions = [];
    scenarioOptions.push( /*#__PURE__*/React.createElement("option", {
      key: "-1",
      value: ""
    }));

    for (let i in this.state.scenarios) {
      let scenario = this.state.scenarios[i];
      scenarioOptions.push( /*#__PURE__*/React.createElement("option", {
        key: i,
        value: scenario.ID
      }, scenario.Name));
    } // form team options


    let teamOptions = [];
    teamOptions.push( /*#__PURE__*/React.createElement("option", {
      key: "-1",
      value: ""
    }));

    for (let i in this.state.teams) {
      let team = this.state.teams[i];
      teamOptions.push( /*#__PURE__*/React.createElement("option", {
        key: i,
        value: team.ID
      }, team.Name));
    } // form time start


    let d = new Date(this.state.timeStart);
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

    d = new Date(this.state.timeEnd);
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
    return /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("label", {
      name: "scenarios"
    }, "Scenarios"), /*#__PURE__*/React.createElement("select", {
      value: this.state.scenarioID,
      onChange: this.selectScenarioCallback
    }, scenarioOptions), /*#__PURE__*/React.createElement("br", null), /*#__PURE__*/React.createElement("label", {
      name: "teams"
    }, "Teams"), /*#__PURE__*/React.createElement("select", {
      value: this.state.teamID,
      onChange: this.selectTeamCallback
    }, teamOptions), /*#__PURE__*/React.createElement("br", null), /*#__PURE__*/React.createElement("label", {
      name: "timeStart"
    }, "Time start"), /*#__PURE__*/React.createElement("input", {
      type: "date",
      value: startDate,
      onChange: this.selectTimeStartCallback
    }), /*#__PURE__*/React.createElement("input", {
      type: "time",
      value: startTime,
      onChange: this.selectTimeStartCallback
    }), /*#__PURE__*/React.createElement("br", null), /*#__PURE__*/React.createElement("label", {
      name: "timeEnd"
    }, "Time end"), /*#__PURE__*/React.createElement("input", {
      type: "date",
      value: endDate,
      onChange: this.selectTimeEndCallback
    }), /*#__PURE__*/React.createElement("input", {
      type: "time",
      value: endTime,
      onChange: this.selectTimeEndCallback
    }), /*#__PURE__*/React.createElement("p", null), /*#__PURE__*/React.createElement("button", {
      onClick: this.submit
    }, "Submit"), /*#__PURE__*/React.createElement(Error, {
      message: this.state.error
    }));
  }

}

class AnalysisResults extends React.Component {
  constructor(props) {
    super(props);
    let config = {
      displaylogo: false
    };
    let layout = {
      height: 450,
      hovermode: 'closest',
      barmode: 'stack',
      xaxis: {
        type: 'date'
      },
      yaxis: {
        domain: [0.70, 1],
        visible: false
      },
      yaxis2: {
        domain: [0.25, 0.60]
      },
      yaxis3: {
        domain: [0, 0.15]
      }
    };
    this.state = {
      error: null,
      config: config,
      layout: layout,
      traces: [],
      selected: {}
    };
    this.plotClick = this.plotClick.bind(this);
  }

  componentWillReceiveProps(newProps) {
    let traces = []; // states

    for (let i in newProps.args.stateTimeline) {
      let hostInstance = newProps.args.stateTimeline[i];
      let name = i + ' - A.states';
      let trace = {
        name: name,
        mode: 'markers',
        x: hostInstance.map(document => document.Document * 1000),
        y: hostInstance.map(_ => name)
      };
      traces.push(trace);
    } // state diffs


    for (let i in newProps.args.stateDiffs) {
      let hostInstance = newProps.args.stateDiffs[i];
      let name = i + ' - B.states diff';
      let trace = {
        name: name,
        type: 'markers+lines',
        x: hostInstance.map(diff => diff.Timestamp * 1000),
        y: hostInstance.map(diff => diff.Changes.length),
        yaxis: 'y2'
      };
      traces.push(trace);
    } // reports


    for (let i in newProps.args.reportTimeline) {
      let hostInstance = newProps.args.reportTimeline[i];
      let name = i + ' - C.reports';
      let trace = {
        name: name,
        mode: 'markers',
        x: hostInstance.map(document => document.Document * 1000),
        y: hostInstance.map(_ => name)
      };
      traces.push(trace);
    } // reports diffs


    for (let i in newProps.args.reportDiffs) {
      let hostInstance = newProps.args.reportDiffs[i];
      let name = i + ' - D.reports diff';
      let trace = {
        name: name,
        type: 'markers+lines',
        x: hostInstance.map(diff => diff.Timestamp * 1000),
        y: hostInstance.map(diff => diff.Changes.length),
        yaxis: 'y2'
      };
      traces.push(trace);
    } // scores


    for (let i in newProps.args.scores) {
      let hostInstance = newProps.args.scores[i];
      let name = i + ' - E.scores';
      let trace = {
        name: name,
        mode: 'markers',
        fill: 'tozeroy',
        x: hostInstance.Timestamps.map(timestamp => timestamp * 1000),
        y: hostInstance.Scores,
        yaxis: 'y3'
      };
      traces.push(trace);
    } // sort traces by name


    traces.sort(function (a, b) {
      if (a.name < b.name) {
        return -1;
      }

      if (a.name > b.name) {
        return 1;
      }

      return 0;
    }); // reverse traces to go from top to bottom in legend

    traces.reverse();
    this.setState({
      error: null,
      reportTimeline: newProps.args.reportTimeline,
      reportDiffs: newProps.args.reportDiffs,
      stateTimeline: newProps.args.stateTimeline,
      stateDiffs: newProps.args.stateDiffs,
      traces: traces
    });
  }

  plotClick(plotlyEvent) {
    plotlyEvent.event.preventDefault(); // only accept left click

    if (plotlyEvent.event.buttons != 1) {
      return false;
    }

    let index = plotlyEvent.points[0].pointIndex;
    let timestamp = Math.trunc(plotlyEvent.points[0].data.x[index] / 1000);
    let type = plotlyEvent.points[0].data.name;

    if (type.endsWith('reports diff') || type.endsWith('states diff')) {
      let diffs = null;

      if (type.endsWith('reports diff')) {
        diffs = this.state.reportDiffs;
      } else if (type.endsWith('states diff')) {
        diffs = this.state.stateDiffs;
      } else {
        return false;
      } // find diff that matches timestamp


      let selected = null;

      for (let i in diffs) {
        if (diffs[i].length <= index) {
          continue;
        }

        if (diffs[i][index].Timestamp != timestamp) {
          continue;
        }

        selected = diffs[i][index];
      }

      if (selected === null) {
        selected = {};
      }

      this.setState({
        selected: selected
      });
    } else if (type.endsWith('reports') || type.endsWith('states')) {
      let documentType = null;
      let timeline = null;

      if (type.endsWith('reports')) {
        documentType = 'reports';
        timeline = this.state.reportTimeline;
      } else if (type.endsWith('states')) {
        documentType = 'states';
        timeline = this.state.stateTimeline;
      } else {
        return false;
      } // get report/state ID that matches timestamp and position


      let id = null;

      for (let i in timeline) {
        if (timeline[i].length <= index) {
          continue;
        }

        if (timeline[i][index].Document === timestamp) {
          id = timeline[i][index].ID;
        }
      }

      if (id === null) {
        this.setState({
          selected: {}
        });
        return false;
      }

      let url = '/analysis/' + documentType + '?id=' + id;
      fetch(url, {
        credentials: 'same-origin'
      }).then(async function (response) {
        if (response.status === 200) {
          let data = await response.json(); // add state ID for state

          if (documentType === "states") {
            data["StateID"] = id;
          }

          return {
            error: null,
            selected: data
          };
        }

        let text = await response.text();
        return {
          error: text
        };
      }).then(function (s) {
        this.setState(s);
      }.bind(this));
    }

    return false;
  }

  render() {
    return /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement(Error, {
      message: this.state.error
    }), "Timeline", /*#__PURE__*/React.createElement("br", null), /*#__PURE__*/React.createElement(Plot, {
      data: this.state.traces,
      layout: this.state.layout,
      config: this.state.config,
      onClick: this.plotClick
    }), /*#__PURE__*/React.createElement("br", null), "Selected", /*#__PURE__*/React.createElement("div", {
      className: "analysisSelected"
    }, /*#__PURE__*/React.createElement(AnalysisSelected, {
      selected: this.state.selected
    })));
  }

}

class AnalysisSelected extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      showAllFindings: false
    };
    this.updateShowAllFindings = this.updateShowAllFindings.bind(this);
  }

  updateShowAllFindings(event) {
    event.preventDefault();
    this.setState({
      showAllFindings: !this.state.showAllFindings
    });
  }

  render() {
    let selected = null;

    if (!this.props.selected) {
      selected = /*#__PURE__*/React.createElement(React.Fragment, null, "No result");
    } // diff
    else if (this.props.selected.Changes != undefined) {
        let time = new Date(this.props.selected.Timestamp * 1000).toLocaleString();
        let changes = [];

        for (let i in this.props.selected.Changes) {
          let change = this.props.selected.Changes[i];
          changes.push( /*#__PURE__*/React.createElement("li", {
            key: i
          }, change.Type, " - ", change.Key, " - ", JSON.stringify(change.Item)));
        }

        selected = /*#__PURE__*/React.createElement(React.Fragment, null, "Time: ", time, /*#__PURE__*/React.createElement("br", null), "Changes:", /*#__PURE__*/React.createElement("ul", null, changes));
      } // report
      else if (this.props.selected.Findings != undefined) {
          let time = new Date(this.props.selected.Timestamp * 1000).toLocaleString();
          let findings = [];

          for (let i in this.props.selected.Findings) {
            let finding = this.props.selected.Findings[i]; // always show findings where Show is true
            // only show findings where Show is false when show all flag is true

            if (finding.Show || this.state.showAllFindings) {
              findings.push( /*#__PURE__*/React.createElement("li", {
                key: i
              }, finding.Value, " - ", finding.Message));
            }
          }

          selected = /*#__PURE__*/React.createElement(React.Fragment, null, "Time: ", time, /*#__PURE__*/React.createElement("br", null), "Findings:", /*#__PURE__*/React.createElement("br", null), /*#__PURE__*/React.createElement("button", {
            disabled: this.state.reportShowAllFindings,
            onClick: this.updateShowAllFindings
          }, "Show/Hide Findings"), /*#__PURE__*/React.createElement("ul", null, findings));
        } // state
        else if (this.props.selected.Users != undefined) {
            let time = new Date(this.props.selected.Timestamp * 1000).toLocaleString();
            let errors = [];

            for (let i in this.props.selected.Errors) {
              let error = this.props.selected.Errors[i];
              errors.push( /*#__PURE__*/React.createElement("li", {
                key: i
              }, error));
            }

            let users = [];

            for (let i in this.props.selected.Users) {
              let user = this.props.selected.Users[i];
              let passwordLastSet = new Date(user.PasswordLastSet * 1000).toLocaleString();
              users.push( /*#__PURE__*/React.createElement("li", {
                key: i
              }, user.ID, " - ", user.Name, " - ", user.AccountActive, " - ", user.AccountExpires, " - ", passwordLastSet, " - ", user.PasswordExpires));
            }

            let groups = [];

            for (let group in this.props.selected.Groups) {
              let members = this.props.selected.Groups[group];

              if (members.length === 0) {
                groups.push( /*#__PURE__*/React.createElement("li", {
                  key: group
                }, group));
              } else {
                let membersStr = members.map(member => member.Name).join(', ');
                groups.push( /*#__PURE__*/React.createElement("li", {
                  key: group
                }, group, " - [", membersStr, "]"));
              }
            }

            let software = [];

            for (let i in this.props.selected.Software) {
              let sw = this.props.selected.Software[i];
              software.push( /*#__PURE__*/React.createElement("li", {
                key: i
              }, sw.Name, " - ", sw.Version));
            }

            let processes = [];

            for (let i in this.props.selected.Processes) {
              let process = this.props.selected.Processes[i];
              processes.push( /*#__PURE__*/React.createElement("li", {
                key: i
              }, process.PID, " - ", process.User, " - ", process.CommandLine));
            }

            let conns = [];

            for (let i in this.props.selected.NetworkConnections) {
              let conn = this.props.selected.NetworkConnections[i];
              conns.push( /*#__PURE__*/React.createElement("li", {
                key: i
              }, conn.Protocol, " - ", conn.LocalAddress, ":", conn.LocalPort, " - ", conn.RemoteAddress, ":", conn.RemotePort, " - ", conn.State));
            }

            let tasks = [];

            for (let i in this.props.selected.ScheduledTasks) {
              let task = this.props.selected.ScheduledTasks[i];
              let enabledStr = "Enabled";

              if (!task.Enabled) {
                enabledStr = "Disabled";
              }

              tasks.push( /*#__PURE__*/React.createElement("li", {
                key: i
              }, task.Name, " @ ", task.Path, " - ", enabledStr));
            }

            let profiles = [];

            for (let i in this.props.selected.WindowsFirewallProfiles) {
              let profile = this.props.selected.WindowsFirewallProfiles[i];
              let enabledStr = "Enabled";

              if (!profile.Enabled) {
                enabledStr = "Disabled";
              }

              profiles.push( /*#__PURE__*/React.createElement("li", {
                key: i
              }, "Profile: ", profile.Name, " - ", enabledStr, " - Inbound: ", profile.DefaultInboundAction, " - Outbound: ", profile.DefaultOutboundAction));
            }

            let rules = [];

            for (let i in this.props.selected.WindowsFirewallRules) {
              let rule = this.props.selected.WindowsFirewallRules[i];
              let enabledStr = "Enabled";

              if (!rule.Enabled) {
                enabledStr = "Disabled";
              }

              rules.push( /*#__PURE__*/React.createElement("li", {
                key: i
              }, "Display Name: ", rule.DisplayName, ", ", enabledStr, ", Protocol: ", rule.Protocol, ", Local Port: ", rule.LocalPort, ", Remote Address: ", rule.RemoteAddress, ", Remote Port: ", rule.RemotePort, ", Direction: ", rule.Direction, ", Action: ", rule.Action));
            }

            let settings = [];

            for (let i in this.props.selected.WindowsSettings) {
              let setting = this.props.selected.WindowsSettings[i];
              settings.push( /*#__PURE__*/React.createElement("li", {
                key: i
              }, "Key: ", setting.Key, ", Value: ", setting.Value));
            }

            selected = /*#__PURE__*/React.createElement(React.Fragment, null, "State ID: ", this.props.selected.StateID, /*#__PURE__*/React.createElement("br", null), "Time: ", time, /*#__PURE__*/React.createElement("br", null), "OS: ", this.props.selected.OS, /*#__PURE__*/React.createElement("br", null), "Hostname: ", this.props.selected.Hostname, /*#__PURE__*/React.createElement("br", null), "Errors:", /*#__PURE__*/React.createElement("ul", null, errors), /*#__PURE__*/React.createElement("br", null), "Users:", /*#__PURE__*/React.createElement("ul", null, users), /*#__PURE__*/React.createElement("br", null), "Groups:", /*#__PURE__*/React.createElement("ul", null, groups), /*#__PURE__*/React.createElement("br", null), "Software:", /*#__PURE__*/React.createElement("ul", null, software), /*#__PURE__*/React.createElement("br", null), "Processes:", /*#__PURE__*/React.createElement("ul", null, processes), /*#__PURE__*/React.createElement("br", null), "Network connections:", /*#__PURE__*/React.createElement("ul", null, conns), /*#__PURE__*/React.createElement("br", null), "Scheduled tasks:", /*#__PURE__*/React.createElement("ul", null, tasks), /*#__PURE__*/React.createElement("br", null), "Windows Firewall profiles:", /*#__PURE__*/React.createElement("ul", null, profiles), /*#__PURE__*/React.createElement("br", null), "Windows Firewall rules:", /*#__PURE__*/React.createElement("ul", null, rules), /*#__PURE__*/React.createElement("br", null), "Windows settings:", /*#__PURE__*/React.createElement("ul", null, settings), /*#__PURE__*/React.createElement("br", null));
          }

    return selected;
  }

}

ReactDOM.render( /*#__PURE__*/React.createElement(App, null), document.getElementById('app'));