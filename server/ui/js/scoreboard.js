'use strict';

const Plot = createPlotlyComponent(Plotly);

class App extends React.Component {
  render() {
    return React.createElement("div", {
      className: "App"
    }, React.createElement(Scoreboard, {
      scenarioID: "1"
    }));
  }

}

class Scoreboard extends React.Component {
  constructor() {
    super();
    this.state = {
      scores: []
    };
  }

  populateScores() {
    let id = this.props.scenarioID;
    let url = '/scores/scenario/' + id;
    fetch(url).then(function (response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }

      return response.json();
    }).then(function (data) {
      this.setState({
        scores: data
      });
    }.bind(this));
  }

  componentDidMount() {
    this.populateScores();
  }

  render() {
    let body = [];

    for (let i in this.state.scores) {
      let entry = this.state.scores[i];
      body.push(React.createElement("tr", {
        key: i
      }, React.createElement("td", null, entry.TeamName), React.createElement("td", null, entry.Score)));
    }

    return React.createElement("div", {
      className: "Scoreboard"
    }, React.createElement("strong", null, "Scoreboard"), React.createElement("p", null), React.createElement("table", null, React.createElement("thead", null, React.createElement("tr", null, React.createElement("th", null, "Team"), React.createElement("th", null, "Score"))), React.createElement("tbody", null, body)));
  }

}

ReactDOM.render(React.createElement(App, null), document.getElementById('app'));