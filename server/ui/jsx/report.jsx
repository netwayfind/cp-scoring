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

  populateHostReport(scenarioID, teamKey, hostname) {
    if (scenarioID === "" || teamKey === "" || hostname === "") {
      return;
    }

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
        mode: 'lines+markers'
      }
    ];

    let layout = {
      xaxis: {
        type: 'date'
      },
      yaxis: {
        fixedrange: true
      }
    }

    let config = {
      displayModeBar: false
    }

    let lastUpdated = null;

    let rows = [];
    if (this.state.report) {
      if (this.state.report.Timestamp) {
        lastUpdated = new Date(this.state.report.Timestamp * 1000).toLocaleString();
      }

      for (let i in this.state.report.Findings) {
        let finding = this.state.report.Findings[i];
        if (!finding.Hidden) {
          rows.push(
            <li key={i}>
              {finding.Value} - {finding.Message}
            </li>
          );
        }
        else {
          rows.push(
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
        let scenarioID = scenarioHosts.ScenarioID;
        let hosts = [];
        for (let i in scenarioHosts.Hosts) {
          let host = scenarioHosts.Hosts[i];
          let hostname = host.Hostname;
          hosts.push(
            <li key={i}><a href="#" onClick={() => this.populateHostReport(scenarioID, teamKey, hostname)}>{hostname}</a></li>
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

    return (
      <React.Fragment>
        <div classname="heading">
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
          <strong>Score Timeline</strong>
          <p />
          <Plot data={data} layout={layout} config={config}/>
          <p />
          Hostname: {this.props.hostname}
          <p />
          Last Updated: {lastUpdated}
          <p />
          Findings:
          <br />
          <ul>{rows}</ul>
        </div>
      </React.Fragment>
    );
  }
}

ReactDOM.render(<App />, document.getElementById('app'));