import { apiGet } from "./common/utils";

import { Component, Fragment } from "react";
import { Link, withRouter } from "react-router-dom/cjs/react-router-dom.min";

class Insight extends Component {
  constructor(props) {
    super(props);
    this.state = {
      hostname: "",
      hostnames: [],
      report: {},
      scenario: {},
      scenarioID: "",
      scenarios: [],
      team: {},
      teamID: "",
      teams: [],
    };

    this.handleUpdateHostname = this.handleUpdateHostname.bind(this);
    this.handleUpdateScenario = this.handleUpdateScenario.bind(this);
    this.handleUpdateTeam = this.handleUpdateTeam.bind(this);
  }

  componentDidMount() {
    this.getData();
  }

  getData() {
    apiGet("/api/scenarios").then(
      function (s) {
        this.setState({
          error: s.error,
          scenarios: s.data,
        });
      }.bind(this)
    );
    apiGet("/api/teams").then(
      function (s) {
        this.setState({
          error: s.error,
          teams: s.data,
        });
      }.bind(this)
    );
  }

  getScenarioReport(hostname) {
    apiGet(
      "/api/insight/" +
        this.state.scenarioID +
        "?hostname=" +
        hostname +
        "&team_id=" +
        this.state.teamID
    ).then(
      function (s) {
        this.setState({
          error: s.error,
          report: s.data,
        });
      }.bind(this)
    );
  }

  getScenarioHostnames() {
    if (!this.state.scenarioID || !this.state.teamID) {
      this.setState({
        hostnames: [],
      });
      return;
    }
    apiGet(
      "/api/insight/" +
        this.state.scenarioID +
        "/hostnames?team_id=" +
        this.state.teamID
    ).then(
      function (s) {
        this.setState({
          error: s.error,
          hostnames: s.data,
        });
      }.bind(this)
    );
  }

  handleUpdateHostname(event) {
    let hostname = event.target.value;
    this.setState({
      hostname: hostname,
    });
    this.getScenarioReport(hostname);
  }

  handleUpdateScenario(event) {
    let id = event.target.value;
    let scenario = {};
    if (id) {
      id = Number(id);
      scenario = this.state.scenarios.find((s) => s.ID === id);
    }
    this.setState({
      hostname: "",
      scenarioID: id,
      scenario: scenario,
    });
    this.getScenarioHostnames();
  }

  handleUpdateTeam(event) {
    let id = event.target.value;
    let team = {};
    if (id) {
      id = Number(id);
      team = this.state.teams.find((s) => s.ID === id);
    }
    this.setState({
      hostname: "",
      teamID: id,
      team: team,
    });
    this.getScenarioHostnames();
  }

  render() {
    let hostnameOptions = [];
    hostnameOptions.push(<option key=""></option>);
    for (let i in this.state.hostnames) {
      let hostname = this.state.hostnames[i];
      hostnameOptions.push(
        <option key={hostname} value={hostname}>
          {hostname}
        </option>
      );
    }
    let scenarioOptions = [];
    scenarioOptions.push(<option key=""></option>);
    for (let i in this.state.scenarios) {
      let scenario = this.state.scenarios[i];
      scenarioOptions.push(
        <option key={scenario.ID} value={scenario.ID}>
          {scenario.Name}
        </option>
      );
    }
    let teamOptions = [];
    teamOptions.push(<option key=""></option>);
    for (let i in this.state.teams) {
      let team = this.state.teams[i];
      teamOptions.push(
        <option key={team.ID} value={team.ID}>
          {team.Name}
        </option>
      );
    }

    let scenarioName = "";
    if (this.state.scenario) {
      scenarioName = this.state.scenario.Name;
    }
    let teamName = "";
    if (this.state.team) {
      teamName = this.state.team.Name;
    }

    let report = "No report found";
    if (this.state.report) {
      let reportTime = null;
      if (this.state.report.Timestamp) {
        reportTime = new Date(
          this.state.report.Timestamp * 1000
        ).toLocaleString();
      }
      let results = [];
      if (this.state.report.AnswerResults) {
        this.state.report.AnswerResults.forEach((result, i) => {
          let entry = (
            <li key={i}>
              <strong>{result.Points}</strong> - {result.Description}
            </li>
          );
          results.push(<li key={i}>{entry}</li>);
        });
      }
      report = (
        <Fragment>
          Report time: {reportTime}
          <p />
          Results: <ul>{results}</ul>
        </Fragment>
      );
    }

    return (
      <Fragment>
        <div className="navbar">
          <Link className="nav-button" target="_blank" to="/admin">
            Admin
          </Link>
        </div>
        <div className="heading">
          <h1>cp-scoring insight</h1>
        </div>
        <div className="toc">
          <h4>Scenarios</h4>
          <select
            onChange={this.handleUpdateScenario}
            value={this.state.scenarioID}
          >
            {scenarioOptions}
          </select>
          <h4>Teams</h4>
          <select onChange={this.handleUpdateTeam} value={this.state.teamID}>
            {teamOptions}
          </select>
        </div>
        <div className="content">
          Scenario: {scenarioName}
          <br />
          Team: {teamName}
          <p />
          Host:{" "}
          <select
            onChange={this.handleUpdateHostname}
            value={this.state.hostname}
          >
            {hostnameOptions}
          </select>
          <p />
          {report}
        </div>
      </Fragment>
    );
  }
}

export default withRouter(Insight);
