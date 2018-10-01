'use strict';

const Plot = createPlotlyComponent(Plotly);

class App extends React.Component {
  render() {
    let scenario = "0";
    let teamKey = "key";
    let hostname = "hostname";
    let query = window.location.search.substring(1);
    let params = query.split("&");

    for (let i = 0; i < params.length; i++) {
      let param = params[i].split("=");

      if (param.length != 2) {
        continue;
      }

      if (param[0] === "scenario") {
        scenario = param[1];
      } else if (param[0] === "team_key") {
        teamKey = param[1];
      } else if (param[0] === "hostname") {
        hostname = param[1];
      }
    }

    return React.createElement("div", {
      className: "App"
    }, React.createElement(ScoreTimeline, {
      scenarioID: scenario,
      teamKey: teamKey,
      hostname: hostname
    }));
  }

}

class ScoreTimeline extends React.Component {
  constructor() {
    super();
    this.state = {
      timestamps: [],
      scores: [],
      report: {}
    };
  }

  populateScores() {
    let scenarioID = this.props.scenarioID;
    let teamKey = this.props.teamKey;
    let hostname = this.props.hostname;
    let url = "/reports/scenario/" + scenarioID + "/timeline?team_key=" + teamKey + "&hostname=" + hostname;
    fetch(url).then(function (response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }

      return response.json();
    }).then(function (data) {
      if (data) {
        // should only be one match
        this.setState({
          scores: data.Scores,
          // timestamps is seconds, need milliseconds
          timestamps: data.Timestamps.map(function (timestamp) {
            return timestamp * 1000;
          })
        });
      }
    }.bind(this));
  }

  populateReport() {
    let scenarioID = this.props.scenarioID;
    let teamKey = this.props.teamKey;
    let hostname = this.props.hostname;
    let url = '/reports/scenario/' + scenarioID + '?team_key=' + teamKey + '&hostname=' + hostname;
    fetch(url).then(function (response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }

      return response.json();
    }).then(function (data) {
      this.setState({
        report: data
      });
    }.bind(this));
  }

  componentDidMount() {
    this.populateScores();
    this.populateReport();
  }

  render() {
    let data = [{
      x: this.state.timestamps,
      y: this.state.scores,
      type: 'scatter',
      mode: 'lines+markers'
    }];
    let layout = {
      xaxis: {
        type: 'date'
      },
      yaxis: {
        fixedrange: true
      }
    };
    let config = {
      displayModeBar: false
    };
    let rows = [];

    if (this.state.report) {
      for (let i in this.state.report.Findings) {
        let finding = this.state.report.Findings[i];

        if (!finding.Hidden) {
          rows.push(React.createElement("li", {
            key: i
          }, finding.Value, " - ", finding.Message));
        } else {
          rows.push(React.createElement("li", {
            key: i
          }, "?"));
        }
      }
    }

    return React.createElement("div", {
      className: "ScoreTimeline"
    }, React.createElement("strong", null, "Score Timeline"), React.createElement("p", null), React.createElement(Plot, {
      data: data,
      layout: layout,
      config: config
    }), React.createElement("ul", null, rows));
  }

}

ReactDOM.render(React.createElement(App, null), document.getElementById('app'));