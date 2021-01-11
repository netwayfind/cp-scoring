import "../App.css";
import { apiGet } from "../common/utils";

import Plot from "react-plotly.js";
import { Component } from "react";
import { withRouter } from "react-router-dom";

class HostReport extends Component {
  constructor(props) {
    super(props);
    this.state = {
      lastRefresh: new Date(),
      report: {
        AnswerResults: [],
      },
      timeline: [],
    };
  }

  componentDidMount() {
    let scenarioID = this.props.match.params.scenarioID;
    let hostname = this.props.match.params.hostname;
    this.getData(scenarioID, hostname);
  }

  componentDidUpdate(prevProps) {
    let scenarioID = this.props.match.params.scenarioID;
    let prevScenarioID = prevProps.match.params.scenarioID;
    let hostname = this.props.match.params.hostname;
    let prevHostname = prevProps.match.params.hostname;
    if (scenarioID !== prevScenarioID || hostname !== prevHostname) {
      this.getData(scenarioID, hostname);
    }
  }

  getData(scenarioID, hostname) {
    let teamKey = this.props.teamKey;
    Promise.all([
      apiGet(
        "/api/scenarios/" +
          scenarioID +
          "/report?team_key=" +
          teamKey +
          "&hostname=" +
          hostname
      ),
      apiGet(
        "/api/scenarios/" +
          scenarioID +
          "/report/timeline?team_key=" +
          teamKey +
          "&hostname=" +
          hostname
      ),
    ]).then(
      async function (responses) {
        let s1 = responses[0];
        let s2 = responses[1];
        this.setState({
          error: s1.error || s2.error,
          lastRefresh: new Date(),
          report: s1.data,
          timeline: s2.data,
        });
      }.bind(this)
    );
  }

  render() {
    let timestampStr = new Date(
      this.state.report.Timestamp * 1000
    ).toLocaleString();
    let score = 0;
    let results = [];
    this.state.report.AnswerResults.forEach((result, i) => {
      results.push(
        <li key={i}>
          <strong>{result.Points}</strong> - {result.Description}
        </li>
      );
      score += result.Points;
    });
    let plotlyData = [];
    this.state.timeline.forEach((timeline) => {
      let timestamps = [];
      timeline.Timestamps.forEach((timestamp) => {
        timestamps.push(new Date(timestamp * 1000));
      });
      plotlyData.push({
        x: timestamps,
        y: timeline.Scores,
        type: "scatter",
        mode: "markers",
        fill: "tozeroy",
      });
    });
    let layout = {
      showlegend: false,
      height: 200,
      margin: {
        t: 25,
        b: 50,
        l: 25,
        r: 25,
      },
      xaxis: {
        fixedrange: true,
      },
      yaxis: {
        fixedrange: true,
      },
    };
    let config = {
      staticPlot: true,
    };
    let scenarioID = this.props.match.params.scenarioID;
    let hostname = this.props.match.params.hostname;

    return (
      <div className="HostReport">
        <Plot data={plotlyData} layout={layout} config={config} />
        <p />
        <button
          type="button"
          onClick={(event) => this.getData(scenarioID, hostname, event)}
        >
          Refresh
        </button>
        &nbsp; Last refresh: {this.state.lastRefresh.toLocaleString()}
        <p />
        Hostname: {hostname}
        <br />
        Report time: {timestampStr}
        <br />
        Score: {score}
        <p />
        <ul>{results}</ul>
      </div>
    );
  }
}

export default withRouter(HostReport);
