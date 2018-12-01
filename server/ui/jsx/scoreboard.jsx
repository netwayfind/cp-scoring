'use strict';

const Plot = createPlotlyComponent(Plotly);

class App extends React.Component {
  render() {
    return (
      <div className="App">
        <Scoreboard />
      </div>
    );
  }
}

class Scoreboard extends React.Component {
  constructor() {
    super();
    this.state = {
      scenarios: [],
      selectedScenarioName: null,
      scores: []
    }
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
  
    fetch(url)
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json();
    })
    .then(function(data) {
      this.setState({
        selectedScenarioName: name,
        scores: data
      });
    }.bind(this));
  }

  getScenarios() {
    let url = '/scores/scenarios';
    
    fetch(url)
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json();
    })
    .then(function(data) {
      this.setState({scenarios: data})
    }.bind(this));
  }

  componentDidMount() {
    this.getScenarios();
  }

  render() {
    let body = [];
    for (let i in this.state.scores) {
      let entry = this.state.scores[i];
      body.push(
        <tr key={i}>
          <td>{entry.TeamName}</td>
          <td>{entry.Score}</td>
        </tr>
      )
    }

    let scenarios = [];
    for (let i in this.state.scenarios) {
      let entry = this.state.scenarios[i];
      scenarios.push(
        <li id={i}>
          <a href="#" onClick={() => {this.populateScores(entry.ID)}}>{entry.Name}</a>
        </li>
      )
    }

    let content = null;

    if (this.state.selectedScenarioName != null) {
      content = (
        <React.Fragment>
        <b>Scenario: </b>{this.state.selectedScenarioName}
        <table>
          <thead>
            <tr>
              <th>Team</th>
              <th>Score</th>
            </tr>
          </thead>
          <tbody>
            {body}
          </tbody>
        </table>
      </React.Fragment>
      );
    }

    return (
      <React.Fragment>
        <div className="heading">
          <h1>Scoreboard</h1>
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