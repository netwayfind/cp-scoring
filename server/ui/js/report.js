'use strict';

const Plot = createPlotlyComponent(Plotly);

class App extends React.Component {
  render() {
    // check for these in query params
    let teamKey = "";
    let hostToken = "";
    let query = window.location.search.substring(1);
    let params = query.split("&");

    for (let i = 0; i < params.length; i++) {
      let param = params[i].split("=");

      if (param.length != 2) {
        continue;
      }

      if (param[0] === "team_key") {
        teamKey = param[1].trim();
      } else if (param[0] == "host_token") {
        hostToken = param[1].trim();
      }
    }

    if (teamKey.length == 0) {
      return React.createElement(AskTeamKey, {
        hostToken: hostToken
      });
    }

    return React.createElement("div", {
      className: "App"
    }, React.createElement(ScoreTimeline, {
      teamKey: teamKey
    }));
  }

}

class AskTeamKey extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      team_key: null
    };
    this.handleChange = this.handleChange.bind(this);
    this.registerHostToken = this.registerHostToken.bind(this);
    this.submit = this.submit.bind(this);
  }

  handleChange(event) {
    let value = event.target.value;
    this.setState({
      [event.target.name]: value
    });
  }

  registerHostToken(hostToken, teamKey) {
    return new Promise(function (resolve, reject) {
      fetch("/token/team", {
        credentials: 'same-origin',
        method: "POST",
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: "team_key=" + teamKey + "&host_token=" + hostToken
      }).then(async function (response) {
        let r = null;

        if (response.status === 200) {
          window.location = "/ui/report?team_key=" + teamKey;
        } else if (response.status === 301) {
          window.location = window.url;
        } else if (response.status === 400) {
          r = {
            error: "Team key required"
          };
        } else if (response.status === 401) {
          r = {
            error: "Invalid team key"
          };
        } else {
          r = {
            error: await response.text()
          };
        }

        resolve(r);
      });
    });
  }

  submit(event) {
    event.preventDefault();
    let hostToken = this.props.hostToken;
    let teamKey = this.state.team_key; // check team key valid

    fetch("/team_key", {
      credentials: "same-origin",
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      },
      body: "team_key=" + teamKey
    }).then(async function (response) {
      if (response.status === 200) {
        // register host token with team
        if (hostToken != undefined && hostToken != null && hostToken.length > 0) {
          return await this.registerHostToken(hostToken, teamKey);
        } else {
          window.location = "/ui/report?team_key=" + teamKey;
          return {};
        }
      } else if (response.status === 400) {
        return {
          error: "Team key required"
        };
      } else if (response.status === 401) {
        return {
          error: "Invalid team key"
        };
      }

      let text = await response.text();
      return {
        error: text
      };
    }.bind(this)).then(function (s) {
      this.setState(s);
    }.bind(this));
  }

  render() {
    return React.createElement(React.Fragment, null, React.createElement("form", {
      onChange: this.handleChange,
      onSubmit: event => this.submit(event)
    }, React.createElement("input", {
      name: "host_token",
      hidden: true,
      value: this.props.hostToken
    }), React.createElement("label", {
      id: "team_key"
    }, "Enter team key:"), React.createElement("input", {
      name: "team_key"
    }), React.createElement("button", {
      type: "submit"
    }, "Submit")), React.createElement(Error, {
      message: this.state.error
    }));
  }

}

class ScoreTimeline extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      scenarioName: "",
      hostname: "",
      timelines: [],
      report: {},
      scenarioHosts: []
    };
  }

  populateScenarioHosts() {
    let teamKey = this.props.teamKey;

    if (teamKey === "") {
      return;
    }

    let url = "/reports?team_key=" + teamKey;
    fetch(url).then(async function (response) {
      if (response.status === 200) {
        let data = await response.json();
        return {
          error: null,
          scenarioHosts: data
        };
      } else if (response.status === 401) {
        window.location = "/ui/report";
      }

      let text = await response.text();
      return {
        error: text
      };
    }).then(function (s) {
      this.setState(s);
    }.bind(this));
  }

  populateHostReport(scenarioName, scenarioID, teamKey, hostname) {
    if (scenarioID === "" || teamKey === "" || hostname === "") {
      return;
    }

    this.setState({
      scenarioName: scenarioName,
      hostname: hostname
    });
    let url = "/reports/scenario/" + scenarioID + "/timeline?team_key=" + teamKey + "&hostname=" + hostname;
    fetch(url).then(async function (response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }

      return await response.json();
    }).then(function (data) {
      if (data) {
        // may have multiple timelines
        let timelines = [];

        for (let i in data) {
          timelines.push({
            scores: data[i].Scores,
            // timestamps is seconds, need milliseconds
            timestamps: data[i].Timestamps.map(function (timestamp) {
              return timestamp * 1000;
            })
          });
        }

        this.setState({
          timelines: timelines
        });
      }
    }.bind(this));
    url = '/reports/scenario/' + scenarioID + '?team_key=' + teamKey + '&hostname=' + hostname;
    fetch(url).then(async function (response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }

      return await response.json();
    }).then(function (data) {
      this.setState({
        report: data
      });
    }.bind(this));
  }

  componentDidMount() {
    this.populateScenarioHosts();
  }

  render() {
    let teamKey = this.props.teamKey;
    let data = []; // plot each timeline

    for (let i in this.state.timelines) {
      let timeline = this.state.timelines[i];
      data.push({
        x: timeline.timestamps,
        y: timeline.scores,
        type: 'scatter',
        mode: 'markers',
        fill: 'tozeroy'
      });
    }

    let layout = {
      showlegend: false,
      height: 200,
      margin: {
        t: 25,
        b: 50,
        l: 25,
        r: 25
      },
      xaxis: {
        fixedrange: true,
        type: 'date'
      },
      yaxis: {
        fixedrange: true
      }
    };
    let config = {
      staticPlot: true
    };
    let lastUpdated = null;
    let score = 0;
    let pointsEarned = 0;
    let pointsLost = 0;
    let findings = [];

    if (this.state.report) {
      if (this.state.report.Timestamp) {
        lastUpdated = new Date(this.state.report.Timestamp * 1000).toLocaleString();
      }

      let fontWeight = null;

      for (let i in this.state.report.Findings) {
        let finding = this.state.report.Findings[i];
        score += finding.Value;

        if (finding.Value >= 0) {
          fontWeight = "normal";
          pointsEarned += finding.Value;
        } else {
          fontWeight = "bold";
          pointsLost += finding.Value;
        }

        if (finding.Show) {
          findings.push(React.createElement("li", {
            key: i
          }, React.createElement("span", {
            style: {
              fontWeight: fontWeight
            }
          }, finding.Value, " - ", finding.Message)));
        } else {
          findings.push(React.createElement("li", {
            key: i
          }, "?"));
        }
      }
    }

    let scenarios = [];

    if (this.state.scenarioHosts) {
      for (let i in this.state.scenarioHosts) {
        let scenarioHosts = this.state.scenarioHosts[i];
        let scenarioName = scenarioHosts.ScenarioName;
        let scenarioID = scenarioHosts.ScenarioID;
        let hosts = [];

        for (let i in scenarioHosts.Hosts) {
          let host = scenarioHosts.Hosts[i];
          let hostname = host.Hostname;
          hosts.push(React.createElement("li", {
            key: i
          }, React.createElement("a", {
            href: "#",
            onClick: () => this.populateHostReport(scenarioName, scenarioID, teamKey, hostname)
          }, hostname)));
        }

        scenarios.push(React.createElement("li", {
          key: i
        }, scenarioHosts.ScenarioName, React.createElement("ul", null, hosts)));
      }
    }

    let content = null;

    if (this.state.hostname) {
      content = React.createElement(React.Fragment, null, React.createElement("h2", null, this.state.scenarioName), React.createElement("h3", null, this.state.hostname), React.createElement(Plot, {
        data: data,
        layout: layout,
        config: config
      }), React.createElement("br", null), "Host instances found: ", this.state.timelines.length, React.createElement("p", null), "Latest Report: ", lastUpdated, React.createElement("p", null), "Report Score: ", score, React.createElement("ul", null, React.createElement("li", null, "Points earned: ", pointsEarned), React.createElement("li", null, "Points lost: ", pointsLost)), React.createElement("p", null), "Report Findings:", React.createElement("br", null), React.createElement("ul", null, findings));
    }

    return React.createElement(React.Fragment, null, React.createElement("div", {
      className: "heading"
    }, React.createElement("h1", null, "Team Reports")), React.createElement("div", {
      className: "toc",
      id: "toc"
    }, React.createElement("h4", null, "Scenarios"), React.createElement("ul", null, scenarios)), React.createElement("div", {
      className: "content",
      id: "content"
    }, React.createElement(Error, {
      message: this.state.error
    }), content));
  }

}

ReactDOM.render(React.createElement(App, null), document.getElementById('app'));