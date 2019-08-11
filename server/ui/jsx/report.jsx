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
      return (
        <AskTeamKey />
      );
    }
    return (
      <div className="App">
        <ScoreTimeline teamKey={teamKey}/>
      </div>
    );
  }
}

class AskTeamKey extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      error: null,
      team_key: null
    }

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

    let teamKey = this.state.team_key;

    // check team key valid
    fetch("/team_key", {
      credentials: "same-origin",
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      },
      body: "team_key=" + teamKey
    })
    .then(async function(response) {
      if (response.status === 200) {
        window.location = "/ui/report?team_key=" + teamKey;
        return {}
      }
      else if (response.status === 400) {
        return {
          error: "Team key required"
        }
      }
      else if (response.status === 401) {
        return {
          error: "Invalid team key"
        }
      }
      let text = await response.text();
      return {
        error: text
      }
    }.bind(this))
    .then(function(s) {
      this.setState(s);
    }.bind(this));
  }

  render() {
    return (
      <React.Fragment>
        <form onChange={this.handleChange} onSubmit={event => this.submit(event)}>
          <label id="team_key">Enter team key:</label>
          <input name="team_key" />
          <button type="submit">Submit</button>
        </form>
        <Error message={this.state.error} />
      </React.Fragment>
    );
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
    }
  }

  populateScenarioHosts() {
    let teamKey = this.props.teamKey;

    if (teamKey === "") {
      return;
    }

    let url = "/reports?team_key=" + teamKey;

    fetch(url)
    .then(async function(response) {
      if (response.status === 200) {
        let data = await response.json();
        return {
          error: null,
          scenarioHosts: data
        }
      }
      else if (response.status === 401) {
        window.location = "/ui/report";
      }
      let text = await response.text();
      return {
        error: text
      }
    })
    .then(function(s) {
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
    fetch(url)
    .then(async function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return await response.json();
    })
    .then(function(data) {
      if (data) {
        // may have multiple timelines
        let timelines = [];
        for (let i in data) {
          timelines.push({
            scores: data[i].Scores,
            // timestamps is seconds, need milliseconds
            timestamps: data[i].Timestamps.map(function(timestamp) {
              return timestamp * 1000;
            })
          });
        }
        this.setState({
          timelines: timelines
        })
      }
    }.bind(this));

    url = '/reports/scenario/' + scenarioID + '?team_key=' + teamKey + '&hostname=' + hostname;  
    fetch(url)
    .then(async function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return await response.json();
    })
    .then(function(data) {
      this.setState({
        report: data
      })
    }.bind(this));
  }

  componentDidMount() {
    this.populateScenarioHosts();
  }

  render() {
    let teamKey = this.props.teamKey;

    let data = [];
    // plot each timeline
    for (let i in this.state.timelines) {
      let timeline = this.state.timelines[i];
      data.push({
        x: timeline.timestamps,
        y: timeline.scores,
        type: 'scatter',
        mode: 'markers',
        fill: 'tozeroy'
      })
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
    }

    let config = {
      staticPlot: true
    }

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
        }
        else {
          fontWeight = "bold";
          pointsLost += finding.Value;
        }
        if (finding.Show) {
          findings.push(
            <li key={i}>
              <span style={{fontWeight: fontWeight}}>{finding.Value} - {finding.Message}</span>
            </li>
          );
        }
        else {
          findings.push(
            <li key={i}>
              ?
            </li>
          )
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
          hosts.push(
            <li key={i}><a href="#" onClick={() => this.populateHostReport(scenarioName, scenarioID, teamKey, hostname)}>{hostname}</a></li>
          );
        }
        scenarios.push(
          <li key={i}>{scenarioHosts.ScenarioName}
            <ul>
              {hosts}
            </ul>
          </li>
        );
      }
    }

    let content = null;
    if (this.state.hostname) {
      content = (
        <React.Fragment>
          <h2>{this.state.scenarioName}</h2>
          <h3>{this.state.hostname}</h3>
          <Plot data={data} layout={layout} config={config}/>
          <br />
          Host instances found: {this.state.timelines.length}
          <p />
          Latest Report: {lastUpdated}
          <p />
          Report Score: {score}
          <ul>
            <li>Points earned: {pointsEarned}</li>
            <li>Points lost: {pointsLost}</li>
          </ul>
          <p />
          Report Findings:
          <br />
          <ul>
            {findings}
          </ul>
        </React.Fragment>
      );
    }

    return (
      <React.Fragment>
        <div className="heading">
          <h1>Team Reports</h1>
        </div>
        <div className="toc" id="toc">
          <h4>Scenarios</h4>
          <ul>
            {scenarios}
          </ul>
        </div>
        <div className="content" id="content">
          <Error message={this.state.error} />
          {content}
        </div>
      </React.Fragment>
    );
  }
}

ReactDOM.render(<App />, document.getElementById('app'));