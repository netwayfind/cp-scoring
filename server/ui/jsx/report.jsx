'use strict';

const Plot = createPlotlyComponent(Plotly);

class App extends React.Component {
  render() {
    let teamKey = "";
    let query = window.location.search.substring(1);
    let params = query.split("&");
    for (let i = 0; i < params.length; i++) {
      let param = params[i].split("=");
      if (param.length != 2) {
        continue;
      }
      if (param[0] === "team_key") {
        teamKey = param[1];
      }
    }
    return (
      <div className="App">
        <ScoreTimeline teamKey={teamKey}/>
      </div>
    );
  }
}

class ScoreTimeline extends React.Component {
  constructor() {
    super();
    this.state = {
      scenarioName: "",
      hostname: "",
      timestamps: [],
      scores: [],
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
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json();
    })
    .then(function(data) {
      if (data) {
        this.setState({
          scenarioHosts: data
        })
      }
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
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json();
    })
    .then(function(data) {
      if (data) {
        // should only be one match
        this.setState({
          scores: data.Scores,
          // timestamps is seconds, need milliseconds
          timestamps: data.Timestamps.map(function(timestamp) {
            return timestamp * 1000;
          })
        })
      }
    }.bind(this));

    url = '/reports/scenario/' + scenarioID + '?team_key=' + teamKey + '&hostname=' + hostname;  
    fetch(url)
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json();
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

    let data = [
      {
        x: this.state.timestamps,
        y: this.state.scores,
        type: 'scatter',
        mode: 'markers',
        fill: 'tozeroy'
      }
    ];

    let layout = {
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
          <p />
          Last Updated: {lastUpdated}
          <p />
          Score: {score}
          <ul>
            <li>Points earned: {pointsEarned}</li>
            <li>Points lost: {pointsLost}</li>
          </ul>
          <p />
          Findings:
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
        <hr />
        <div className="toc" id="toc">
          <b>Scenarios</b>
          <ul>
            {scenarios}
          </ul>
        </div>
        <div className="content" id="content">
          {content}
        </div>
      </React.Fragment>
    );
  }
}

ReactDOM.render(<App />, document.getElementById('app'));