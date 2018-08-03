'use strict';

const Plot = createPlotlyComponent(Plotly);

class App extends React.Component {
  render() {
    return (
      <div className="App">
        <Scoreboard scenarioID="1"/>
      </div>
    );
  }
}

class Scoreboard extends React.Component {
  constructor() {
    super();
    this.state = {
      scores: []
    }
  }

  populateScores() {
    let id = this.props.scenarioID;
    let url = '/scenarios/' + id + '/scores';
  
    fetch(url)
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json();
    })
    .then(function(data) {
      this.setState({scores: data})
    }.bind(this));
  }

  componentDidMount() {
    this.populateScores();
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

    return (
      <div className="Scoreboard">
        <strong>Scoreboard</strong>
        <p />
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
      </div>
    );
  }
}

ReactDOM.render(<App />, document.getElementById('app'));