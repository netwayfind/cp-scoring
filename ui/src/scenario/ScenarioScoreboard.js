import "../App.css";
import { apiGet } from "../common/utils";

import { Component } from "react";
import { withRouter } from "react-router-dom";

class ScenarioScoreboard extends Component {
  constructor(props) {
    super(props);
    this.state = {
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
          scoreboard: s.data,
        });
      }.bind(this)
    );
  }

  render() {
    let entries = [];
    this.state.scoreboard.forEach((entry) => {
      let timestampStr = new Date(entry.Timestamp * 1000).toLocaleString();
      entries.push(
        <p>
          {entry.TeamName} - {entry.Hostname} - {entry.Score} - {timestampStr}
        </p>
      );
    });
    return (
      <div class="ScenarioScoreboard">
        <ul>{entries}</ul>
      </div>
    );
  }
}

export default withRouter(ScenarioScoreboard);
