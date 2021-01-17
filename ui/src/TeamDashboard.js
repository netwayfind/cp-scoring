import "./App.css";
import { apiGet, apiPost } from "./common/utils";
import HostReport from "./team-dashboard/HostReport";
import ScenarioDesc from "./team-dashboard/ScenarioDesc";

import { Component, Fragment } from "react";
import { Link, Route, Switch, withRouter } from "react-router-dom";

class TeamDashboard extends Component {
  constructor(props) {
    super(props);
    this.state = {
      authenticated: false,
      error: null,
      scenarios: [],
      scenarioHosts: {},
      teamKey: "",
    };

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDidUpdate(prevProps) {
    let scenarioID = this.getScenarioID(this.props);
    let prevScenarioID = this.getScenarioID(prevProps);
    if (scenarioID !== prevScenarioID) {
      this.setState({
        scenarioID: scenarioID
      });
    }
  }

  getScenarioID(props) {
    return Number(props.location.pathname.replace(props.match.url + "/scenario/", "").split("/")[0]);
  }

  getData() {
    apiGet("/api/scoreboard/scenarios").then(
      function (s) {
        let scenarios = s.data;
        this.setState({
          error: s.error,
          scenarios: scenarios,
        });
        if (!s.error) {
          scenarios.forEach(scenario => {
            apiGet(
              "/api/report/" + scenario.ID + "/hostnames?team_key=" + this.state.teamKey
            ).then(
              async function (s) {
                this.setState({
                  error: s.error,
                  scenarioHosts: {
                    ...this.state.scenarioHosts,
                    [scenario.ID]: s.data
                  }
                });
              }.bind(this)
            );
          });
        }
      }.bind(this)
    );
  }

  handleChange(event) {
    let value = event.target.value;
    this.setState({
      [event.target.name]: value,
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
        this.getData();
      }.bind(this)
    );
  }

  render() {
    if (!this.state.authenticated) {
      return (
        <Fragment>
          <form onChange={this.handleChange} onSubmit={this.handleSubmit}>
            <label htmlFor="teamKey">Team key</label>
            <input name="teamKey" required="required" />
            <button type="submit">Submit</button>
          </form>
          <h1>{this.state.error}</h1>
        </Fragment>
      );
    }

    let scenarioID = this.getScenarioID(this.props);
    let scenarios = [];
    if (this.state.scenarios) {
      this.state.scenarios.forEach((scenario) => {
        let scenarioHosts = null;
        if (scenario.ID === scenarioID) {
          let hostnames = [];
          let hosts = this.state.scenarioHosts[scenario.ID];
          if (hosts) {
            hosts.forEach((hostname) => {
              hostnames.push(
                <li key={hostname}>
                  <Link
                    to={`${this.props.match.path}/scenario/${scenarioID}/${hostname}`}
                  >
                    {hostname}
                  </Link>
                </li>
              );
            });
            scenarioHosts = <ul>{hostnames}</ul>;
          }
        }
        let entry = (
          <li key={scenario.ID}>
            <Link to={`${this.props.match.path}/scenario/${scenario.ID}`}>
              {scenario.Name}
            </Link>
            {scenarioHosts}
          </li>
        );
        scenarios.push(entry);
      });
    }

    return (
      <Fragment>
        <div className="navbar">
          <Link className="nav-button" target="_blank" to="/scoreboard">
            Scoreboard
          </Link>
        </div>
        <div className="heading">
          <h1>Team Dashboard</h1>
        </div>
        <div className="toc">
          <h4>Scenarios</h4>
          <ul>{scenarios}</ul>
        </div>
        <div className="content">
          <Switch>
            <Route exact path={`${this.props.match.url}/scenario/:scenarioID`}>
              <ScenarioDesc teamKey={this.state.teamKey} />
            </Route>
            <Route
              exact
              path={`${this.props.match.url}/scenario/:scenarioID/:hostname`}
            >
              <HostReport teamKey={this.state.teamKey} />
            </Route>
          </Switch>
        </div>
      </Fragment>
    );
  }
}

export default withRouter(TeamDashboard);