'use strict';

const Plot = createPlotlyComponent(Plotly);

class App extends React.Component {
  render() {
    return React.createElement("div", {
      className: "App"
    }, React.createElement(Scoreboard, null));
  }

}

class Scoreboard extends React.Component {
  constructor() {
    super();
    this.state = {
      scenarios: [],
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

    fetch(url).then(function (response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }

      return response.json();
    }).then(function (data) {
      this.setState({
        selectedScenarioName: name,
        scores: data
      });
    }.bind(this));
  }

  getScenarios() {
    let url = '/scores/scenarios';
    fetch(url).then(function (response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }

      return response.json();
    }).then(function (data) {
      this.setState({
        scenarios: data
      });
    }.bind(this));
  }

  componentDidMount() {
    this.getScenarios();
  }

  render() {
    let body = [];

    for (let i in this.state.scores) {
      let entry = this.state.scores[i];
      body.push(React.createElement("tr", {
        key: i
      }, React.createElement("td", null, entry.TeamName), React.createElement("td", null, entry.Score)));
    }

    let scenarios = [];

    for (let i in this.state.scenarios) {
      let entry = this.state.scenarios[i];
      scenarios.push(React.createElement("li", {
        id: i
      }, React.createElement("a", {
        href: "#",
        onClick: () => {
          this.populateScores(entry.ID);
        }
      }, entry.Name)));
    }

    let content = null;

    if (this.state.selectedScenarioName != null) {
      content = React.createElement(React.Fragment, null, React.createElement("b", null, "Scenario: "), this.state.selectedScenarioName, React.createElement("table", null, React.createElement("thead", null, React.createElement("tr", null, React.createElement("th", null, "Team"), React.createElement("th", null, "Score"))), React.createElement("tbody", null, body)));
    }

    return React.createElement(React.Fragment, null, React.createElement("div", {
      classname: "heading"
    }, React.createElement("h1", null, "Scoreboard")), React.createElement("hr", null), React.createElement("div", {
      className: "toc",
      id: "toc"
    }, React.createElement("b", null, "Scenarios"), React.createElement("ul", null, scenarios)), React.createElement("div", {
      className: "content",
      id: "content"
    }, content));
  }

}

ReactDOM.render(React.createElement(App, null), document.getElementById('app'));