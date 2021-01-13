import "./App.css";
import { apiGet, apiPost } from "./common/utils";
import HostReport from "./report/HostReport";

import { Component, Fragment } from "react";
import { Link, Route, Switch, withRouter } from "react-router-dom";

class Report extends Component {
  constructor(props) {
    super(props);
    this.state = {
      authenticated: false,
      error: null,
      hostnames: [],
      scenarios: [],
      scenarioID: null,
      teamKey: "",
    };

    this.handleChange = this.handleChange.bind(this);
    this.handleScenarioUpdate = this.handleScenarioUpdate.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDidMount() {
    this.getScenarios();
  }

  getScenarios() {
    apiGet("/api/scoreboard/scenarios").then(
      function (s) {
        this.setState({
          error: s.error,
          scenarios: s.data,
        });
      }.bind(this)
    );
  }

  getData(scenarioID, teamKey) {
    apiGet("/api/report/" + scenarioID + "/hostnames?team_key=" + teamKey).then(
      async function (s) {
        this.setState({
          error: s.error,
          hostnames: s.data,
        });
      }.bind(this)
    );
  }

  handleChange(event) {
    let value = event.target.value;
    this.setState({
      [event.target.name]: value,
    });
  }

  handleScenarioUpdate(scenarioID, event) {
    event.preventDefault();

    this.getData(scenarioID, this.state.teamKey);
    this.setState({
      scenarioID: scenarioID,
    });
  }

  handleSubmit(event) {
    event.preventDefault();

    apiPost("/api/login/team", {
      TeamKey: this.state.teamKey,
    }).then(
      async function (s) {
        let authenticated = false;
        if (!s.error) {
          authenticated = true;
        }
        this.setState({
          authenticated: authenticated,
          error: s.error,
        });
      }.bind(this)
    );
  }

  render() {
    if (!this.state.authenticated) {
      return (
        <div className="Report">
          <form onChange={this.handleChange} onSubmit={this.handleSubmit}>
            <label htmlFor="teamKey">Team key</label>
            <input name="teamKey" required="required" />
            <button type="submit">Submit</button>
          </form>
          <h1>{this.state.error}</h1>
        </div>
      );
    }

    let scenarios = [];
    if (this.state.scenarios) {
      this.state.scenarios.forEach((scenario) => {
        scenarios.push(
          <li key={scenario.ID}>
            <button
              type="button"
              disabled={this.state.scenarioID === scenario.ID}
              onClick={(event) => this.handleScenarioUpdate(scenario.ID, event)}
            >
              {scenario.Name}
            </button>
          </li>
        );
      });
    }

    let hostnames;
    if (this.state.hostnames.length > 0) {
      hostnames = [];
      this.state.hostnames.forEach((hostname) => {
        hostnames.push(
          <li key={hostname}>
            <Link
              to={`${this.props.match.url}/${this.state.scenarioID}/${hostname}`}
            >
              {hostname}
            </Link>
          </li>
        );
      });
    }

    return (
      <Fragment>
        <div className="heading">
          <h1>Report</h1>
        </div>
        <div className="toc">
          <h4>Scenarios</h4>
          <ul>{scenarios}</ul>
          <hr />
          <h4>Hosts</h4>
          <ul>{hostnames}</ul>
        </div>
        <div className="content">
          <Switch>
            <Route path={`${this.props.match.url}/:scenarioID/:hostname`}>
              <HostReport teamKey={this.state.teamKey} />
            </Route>
          </Switch>
        </div>
      </Fragment>
    );
  }
}

export default withRouter(Report);
