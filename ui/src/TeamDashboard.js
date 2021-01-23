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
      teamName: "",
    };

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDidMount() {
    apiGet("/api/login-team/").then(
      async function (s) {
        if (!s.error) {
          this.getData();
          this.setState({
            authenticated: true,
            teamName: s.data,
          });
        }
      }.bind(this)
    );
  }

  componentDidUpdate(prevProps) {
    let scenarioID = this.getScenarioID(this.props);
    let prevScenarioID = this.getScenarioID(prevProps);
    if (scenarioID !== prevScenarioID) {
      this.setState({
        scenarioID: scenarioID,
      });
    }
  }

  getHostname(props) {
    return props.location.pathname
      .replace(props.match.url + "/scenario/", "")
      .split("/")[1];
  }

  getScenarioID(props) {
    return Number(
      props.location.pathname
        .replace(props.match.url + "/scenario/", "")
        .split("/")[0]
    );
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
          this.getScenarioHosts(scenarios);
        }
      }.bind(this)
    );
  }

  getScenarioHosts(scenarios) {
    scenarios.forEach((scenario) => {
      apiGet("/api/report/" + scenario.ID + "/hostnames").then(
        async function (s) {
          this.setState({
            error: s.error,
            scenarioHosts: {
              ...this.state.scenarioHosts,
              [scenario.ID]: s.data,
            },
          });
        }.bind(this)
      );
    });
  }

  handleChange(event) {
    let value = event.target.value;
    this.setState({
      [event.target.name]: value,
    });
  }

  handleSubmit(event) {
    event.preventDefault();

    apiPost("/api/login-team/", {
      TeamKey: this.state.teamKey,
    }).then(
      async function (s) {
        let authenticated = false;
        if (!s.error) {
          this.getData();
          authenticated = true;
        }
        this.setState({
          authenticated: authenticated,
          error: s.error,
          teamName: s.data,
        });
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

    let currentHostname = this.getHostname(this.props);
    let currentScenarioID = this.getScenarioID(this.props);
    let currentScenarioName = null;
    let scenarios = [];
    if (this.state.scenarios) {
      this.state.scenarios.forEach((scenario) => {
        let scenarioHosts = null;
        if (scenario.ID === currentScenarioID) {
          currentScenarioName = scenario.Name;
          let hostnames = [];
          let hosts = this.state.scenarioHosts[scenario.ID];
          if (hosts) {
            hosts.forEach((hostname) => {
              let linkClassesHost = ["nav-button"];
              if (hostname === currentHostname) {
                linkClassesHost.push("nav-button-selected");
              }
              hostnames.push(
                <li key={hostname}>
                  <Link
                    className={linkClassesHost.join(" ")}
                    to={`${this.props.match.path}/scenario/${currentScenarioID}/${hostname}`}
                  >
                    {hostname}
                  </Link>
                </li>
              );
            });
            scenarioHosts = <ul>{hostnames}</ul>;
          }
        }

        let linkClassesScenario = ["nav-button"];
        if (
          this.props.location.pathname ===
          this.props.match.path + "/scenario/" + scenario.ID
        ) {
          linkClassesScenario.push("nav-button-selected");
        }

        let entry = (
          <li key={scenario.ID}>
            <Link
              className={linkClassesScenario.join(" ")}
              to={`${this.props.match.path}/scenario/${scenario.ID}`}
            >
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
          <h1>Team Dashboard - {this.state.teamName}</h1>
        </div>
        <div className="toc">
          <h4>Scenarios</h4>
          <ul>{scenarios}</ul>
        </div>
        <div className="content">
          <Switch>
            <Route exact path={`${this.props.match.url}/scenario/:scenarioID`}>
              <ScenarioDesc />
            </Route>
            <Route
              exact
              path={`${this.props.match.url}/scenario/:scenarioID/:hostname`}
            >
              <HostReport scenarioName={currentScenarioName} />
            </Route>
          </Switch>
        </div>
      </Fragment>
    );
  }
}

export default withRouter(TeamDashboard);
