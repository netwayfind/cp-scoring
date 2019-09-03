'use strict';

const Plot = createPlotlyComponent(Plotly);

class App extends React.Component {
  render() {
    // check for these in query params
    let teamKey = "";
    let query = window.location.search.substring(1);
    let params = query.split("&");

    for (let i = 0; i < params.length; i++) {
      let param = params[i].split("=");

      if (param.length != 2) {
        continue;
      }

      if (param[0] === "team_key") {
        teamKey = param[1].trim();
      }
    }

    if (teamKey.length == 0) {
      return React.createElement(AskTeamKey, null);
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
    this.submit = this.submit.bind(this);
  }

  handleChange(event) {
    let value = event.target.value;
    this.setState({
      [event.target.name]: value
    });
  }

  submit(event) {
    event.preventDefault();
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
        window.location = "/ui/report?team_key=" + teamKey;
        return {};
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
    }, React.createElement("label", {
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
      timelines: [],
      report: {},
      scenarioHosts: [],
      selectedScenarioID: null,
      selectedScenarioName: null,
      selectedScenarioHostname: null,
      lastCheck: null
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
      selectedScenarioID: scenarioID,
      selectedScenarioName: scenarioName,
      selectedScenarioHostname: hostname
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
          timelines: timelines,
          lastCheck: new Date().toLocaleString()
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
        report: data,
        lastCheck: new Date().toLocaleString()
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

    if (this.state.selectedScenarioHostname) {
      content = React.createElement(React.Fragment, null, React.createElement("h2", null, this.state.selectedScenarioName), React.createElement("h3", null, this.state.selectedScenarioHostname), React.createElement("p", null), "Last updated: ", this.state.lastCheck, React.createElement("br", null), React.createElement("button", {
        onClick: () => this.populateHostReport(this.state.selectedScenarioName, this.state.selectedScenarioID, this.props.teamKey, this.state.selectedScenarioHostname)
      }, "Refresh"), React.createElement("p", null), React.createElement(Plot, {
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