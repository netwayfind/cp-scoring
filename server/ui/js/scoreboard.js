'use strict';

class App extends React.Component {
  render() {
    return /*#__PURE__*/React.createElement("div", {
      className: "App"
    }, /*#__PURE__*/React.createElement(Scoreboard, null));
  }

}

class Scoreboard extends React.Component {
  constructor() {
    super();
    this.state = {
      error: null,
      scenarios: [],
      lastCheck: null,
      selectedScenarioID: null,
      selectedScenarioName: null,
      scores: []
    };
  }

  populateScores(id) {
    if (id === undefined || id === null || !id) {
      return;
    }

    let url = '/scores/scenario/' + id;
    id = Number(id);
    let name = null;

    for (let i in this.state.scenarios) {
      let entry = this.state.scenarios[i];

      if (entry.ID === id) {
        name = entry.Name;
        break;
      }
    }

    fetch(url).then(async function (response) {
      let lastCheck = new Date().toLocaleString();

      if (response.status === 200) {
        let data = await response.json();
        return {
          lastCheck: lastCheck,
          error: null,
          selectedScenarioID: id,
          selectedScenarioName: name,
          scores: data
        };
      }

      let text = await response.text();
      return {
        lastCheck: lastCheck,
        error: text
      };
    }).then(function (s) {
      this.setState(s);
    }.bind(this));
  }

  getScenarios() {
    let url = '/scores/scenarios';
    fetch(url).then(async function (response) {
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
  }

  componentDidMount() {
    this.getScenarios();
  }

  render() {
    let body = [];

    for (let i in this.state.scores) {
      let entry = this.state.scores[i];
      let lastTimestamp = 0;
      let totalScore = 0;
      let hostScoreDetails = [];

      for (let j in entry.HostScores) {
        let hostScore = entry.HostScores[j];
        totalScore += hostScore.Score;

        if (hostScore.Timestamp > lastTimestamp) {
          lastTimestamp = hostScore.Timestamp;
        }

        let lastUpdated = new Date(hostScore.Timestamp * 1000).toLocaleString();
        hostScoreDetails.push( /*#__PURE__*/React.createElement("p", null, "Time: ", lastUpdated, "; Hostname: ", hostScore.Hostname, "; Score: ", hostScore.Score, " "));
      }

      let lastTimestampStr = new Date(lastTimestamp * 1000).toLocaleString();
      body.push( /*#__PURE__*/React.createElement("tr", {
        key: i
      }, /*#__PURE__*/React.createElement("td", {
        class: "table-cell"
      }, entry.TeamName), /*#__PURE__*/React.createElement("td", {
        class: "table-cell"
      }, totalScore), /*#__PURE__*/React.createElement("td", {
        class: "table-cell"
      }, lastTimestampStr), /*#__PURE__*/React.createElement("td", {
        class: "table-cell"
      }, /*#__PURE__*/React.createElement("details", null, hostScoreDetails))));
    }

    let scenarios = [];

    for (let i in this.state.scenarios) {
      let entry = this.state.scenarios[i];
      let classes = ["nav-button"];

      if (this.state.selectedScenarioID === entry.ID) {
        classes.push("nav-button-selected");
      }

      scenarios.push( /*#__PURE__*/React.createElement("li", {
        id: i
      }, /*#__PURE__*/React.createElement("a", {
        className: classes.join(" "),
        href: "#",
        onClick: () => {
          this.populateScores(entry.ID);
        }
      }, entry.Name)));
    }

    let content = null;

    if (this.state.selectedScenarioName != null) {
      content = /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("h2", null, this.state.selectedScenarioName), /*#__PURE__*/React.createElement("p", null), "Last updated: ", this.state.lastCheck, /*#__PURE__*/React.createElement("br", null), /*#__PURE__*/React.createElement("button", {
        onClick: () => {
          this.populateScores(this.state.selectedScenarioID);
        }
      }, "Refresh"), /*#__PURE__*/React.createElement("p", null), /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", {
        class: "table-cell"
      }, "Team Name"), /*#__PURE__*/React.createElement("th", {
        class: "table-cell"
      }, "Score"), /*#__PURE__*/React.createElement("th", {
        class: "table-cell"
      }, "Last Updated"))), /*#__PURE__*/React.createElement("tbody", null, body)));
    }

    return /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("div", {
      className: "heading"
    }, /*#__PURE__*/React.createElement("h1", null, "Scoreboard")), /*#__PURE__*/React.createElement("div", {
      className: "toc",
      id: "toc"
    }, /*#__PURE__*/React.createElement("h4", null, "Scenarios"), /*#__PURE__*/React.createElement("ul", null, scenarios)), /*#__PURE__*/React.createElement("div", {
      className: "content",
      id: "content"
    }, /*#__PURE__*/React.createElement(Error, {
      message: this.state.error
    }), content));
  }

}

ReactDOM.render( /*#__PURE__*/React.createElement(App, null), document.getElementById('app'));