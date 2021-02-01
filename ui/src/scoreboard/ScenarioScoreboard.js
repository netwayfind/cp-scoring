import "../App.css";
import { apiGet } from "../common/utils";

import { Component, Fragment } from "react";
import { withRouter } from "react-router-dom";

class ScenarioScoreboard extends Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      lastUpdated: new Date(),
      scoreboard: [],
    };
  }

  componentDidMount() {
    let id = this.props.match.params.id;
    this.getData(id);
  }

  componentDidUpdate(prevProps) {
    let id = this.props.match.params.id;
    let prevId = prevProps.match.params.id;
    if (id !== prevId) {
      this.getData(id);
    }
  }

  getData(id) {
    apiGet("/api/scoreboard/scenarios/" + id).then(
      function (s) {
        this.setState({
          error: s.error,
          lastUpdated: new Date(),
          scoreboard: s.data,
        });
      }.bind(this)
    );
  }

  render() {
    let lastUpdatedStr = this.state.lastUpdated.toLocaleString();
    let id = this.props.match.params.id;

    let teamToDetails = {};
    let teamScores = [];
    let teamScoresIndex = 0;
    let currentTeamID = -1;
    let currentTeamScore = {};
    this.state.scoreboard.forEach((entry) => {
      if (entry.TeamID !== currentTeamID) {
        currentTeamScore = {
          score: 0,
          lastUpdated: 0,
        };
        currentTeamID = entry.TeamID;
        teamScoresIndex = teamScores.length;
      }

      let teamHosts = teamToDetails[entry.TeamID];
      if (!teamHosts) {
        teamHosts = [];
      }
      teamHosts.push(
        <tr key={entry.Hostname}>
          <td>{entry.Hostname}</td>
          <td>{entry.Score}</td>
          <td>{new Date(entry.Timestamp * 1000).toLocaleString()}</td>
        </tr>
      );
      teamToDetails[entry.TeamID] = teamHosts;

      let newScore = currentTeamScore.score + entry.Score;
      let newLastUpdated = currentTeamScore.lastUpdated;
      if (newLastUpdated < entry.Timestamp) {
        newLastUpdated = entry.Timestamp;
      }
      currentTeamScore = {
        teamID: entry.TeamID,
        teamName: entry.TeamName,
        score: newScore,
        lastUpdated: newLastUpdated,
      };
      teamScores[teamScoresIndex] = currentTeamScore;
    });

    let tableBody = [];
    teamScores.forEach((entry, i) => {
      tableBody.push(
        <tr key={i}>
          <td className="table-cell">{entry.teamName}</td>
          <td className="table-cell">{entry.score}</td>
          <td className="table-cell">
            {new Date(entry.lastUpdated * 1000).toLocaleString()}
          </td>
          <td>
            <details>
              <table>
                <tr>
                  <th>Host</th>
                  <th>Score</th>
                  <th>Last Updated</th>
                </tr>
                {teamToDetails[entry.teamID]}
              </table>
            </details>
          </td>
        </tr>
      );
    });

    let table = (
      <table>
        <thead>
          <tr>
            <th className="table-cell">Team Name</th>
            <th className="table-cell">Score</th>
            <th className="table-cell">Last Updated</th>
          </tr>
        </thead>
        <tbody>{tableBody}</tbody>
      </table>
    );

    return (
      <Fragment>
        Last updated: {lastUpdatedStr}
        <br />
        <button type="button" onClick={() => this.getData(id)}>
          Refresh
        </button>
        <p />
        {table}
      </Fragment>
    );
  }
}

export default withRouter(ScenarioScoreboard);
