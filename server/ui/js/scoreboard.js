'use strict';

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
      error: null,
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

    fetch(url).then(async function (response) {
      if (response.status === 200) {
        let data = await response.json();
        return {
          error: null,
          selectedScenarioName: name,
          scores: data
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
      let lastUpdated = new Date(entry.Timestamp * 1000).toLocaleString();
      body.push(React.createElement("tr", {
        key: i
      }, React.createElement("td", {
        class: "table-cell"
      }, entry.TeamName), React.createElement("td", {
        class: "table-cell"
      }, entry.Score), React.createElement("td", {
        class: "table-cell"
      }, lastUpdated)));
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
      content = React.createElement(React.Fragment, null, React.createElement("h2", null, this.state.selectedScenarioName), React.createElement("p", null), React.createElement("table", null, React.createElement("thead", null, React.createElement("tr", null, React.createElement("th", {
        class: "table-cell"
      }, "Team Name"), React.createElement("th", {
        class: "table-cell"
      }, "Score"), React.createElement("th", {
        class: "table-cell"
      }, "Last Updated"))), React.createElement("tbody", null, body)));
    }

    return React.createElement(React.Fragment, null, React.createElement("div", {
      className: "heading"
    }, React.createElement("h1", null, "Scoreboard")), React.createElement("div", {
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