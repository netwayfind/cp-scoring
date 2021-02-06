import "../App.css";
import { apiDelete, apiGet, apiPost, apiPut } from "../common/utils";
import ScenarioHost from "./ScenarioHost";

import { Component, Fragment } from "react";
import { withRouter } from "react-router-dom/cjs/react-router-dom.min";

class Scenario extends Component {
  constructor(props) {
    super(props);
    this.state = this.defaultState();

    this.getData = this.getData.bind(this);
    this.handleDelete = this.handleDelete.bind(this);
    this.handleUpdate = this.handleUpdate.bind(this);
    this.handleSave = this.handleSave.bind(this);
    this.handleSaveHost = this.handleSaveHost.bind(this);
    this.handleHostnameDelete = this.handleHostnameDelete.bind(this);
    this.handleHostnameRename = this.handleHostnameRename.bind(this);
    this.handleHostnameSelect = this.handleHostnameSelect.bind(this);
    this.handleNewHostnameAdd = this.handleNewHostnameAdd.bind(this);
    this.handleNewHostnameUpdate = this.handleNewHostnameUpdate.bind(this);
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

  defaultState() {
    return {
      error: null,
      currentScenarioHostname: "",
      currentScenarioHost: {},
      newScenarioHostname: "",
      scenario: {},
      scenarioHosts: {},
    };
  }

  getData(id) {
    if (id === undefined) {
      this.setState(this.defaultState);
      return;
    }
    Promise.all([
      apiGet("/api/scenarios/" + id),
      apiGet("/api/scenarios/" + id + "/hosts"),
    ]).then(
      async function (responses) {
        let s1 = responses[0];
        let s2 = responses[1];
        this.setState({
          error: s1.error || s2.error,
          currentScenarioHostname: "",
          currentScenarioHost: {},
          newScenarioHostname: "",
          scenario: s1.data,
          scenarioHosts: s2.data,
        });
      }.bind(this)
    );
  }

  handleDelete() {
    apiDelete("/api/scenarios/" + this.state.scenario.ID).then(
      async function (s) {
        if (s.error) {
          this.setState({
            error: s.error,
          });
        } else {
          this.props.parentCallback();
          this.props.history.push(this.props.parentPath);
        }
      }.bind(this)
    );
  }

  handleSave(event) {
    if (event !== null) {
      event.preventDefault();
    }
    let id = this.state.scenario.ID;
    if (id) {
      // update
      apiPut("/api/scenarios/" + id, this.state.scenario).then(
        async function (s) {
          if (s.error) {
            this.setState({
              error: s.error,
            });
          } else {
            this.props.parentCallback();
            this.props.history.push(this.props.match.url);
          }
        }.bind(this)
      );
    } else {
      // create
      apiPost("/api/scenarios/", this.state.scenario).then(
        async function (s) {
          if (s.error) {
            this.setState({
              error: s.error,
            });
          } else {
            this.props.parentCallback();
            this.props.history.push(this.props.match.url + "/" + s.data.ID);
          }
        }.bind(this)
      );
    }
  }

  handleSaveHost(checks, answers, config) {
    let id = this.state.scenario.ID;
    let scenarioHost = {
      Checks: checks,
      Answers: answers,
      Config: config,
    };
    let scenarioHosts = {
      ...this.state.scenarioHosts,
    };
    scenarioHosts[this.state.currentScenarioHostname] = scenarioHost;
    apiPut("/api/scenarios/" + id + "/hosts", scenarioHosts).then(
      async function (s) {
        if (s.error) {
          this.setState({
            error: s.error,
          });
        } else {
          this.setState({
            scenarioHosts: scenarioHosts,
          });
          this.props.parentCallback();
          this.props.history.push(this.props.match.url);
        }
      }.bind(this)
    );
  }

  handleUpdate(event) {
    let value = event.target.value;
    if (event.target.type === "checkbox") {
      value = event.target.checked;
    }
    this.setState({
      scenario: {
        ...this.state.scenario,
        [event.target.name]: value,
      },
    });
  }

  handleHostnameDelete() {
    let hostname = this.state.currentScenarioHostname;
    if (!hostname) {
      return;
    }
    let scenarioHosts = {
      ...this.state.scenarioHosts,
    };
    delete scenarioHosts[hostname];
    this.setState({
      currentScenarioHostname: "",
      currentScenarioHost: {},
      scenarioHosts: scenarioHosts,
    });
  }

  handleHostnameRename(event) {
    let hostname = event.target.value;
    if (!hostname) {
      return;
    }
    let oldHostname = this.state.currentScenarioHostname;
    let scenarioHost = this.state.currentScenarioHost;
    let scenarioHosts = {
      ...this.state.scenarioHosts,
      [hostname]: scenarioHost,
    };
    delete scenarioHosts[oldHostname];
    this.setState({
      currentScenarioHostname: hostname,
      scenarioHosts: scenarioHosts,
    });
  }

  handleHostnameSelect(event) {
    let hostname = event.target.value;
    this.setState({
      currentScenarioHostname: hostname,
      currentScenarioHost: this.state.scenarioHosts[hostname] || {},
    });
  }

  handleNewHostnameAdd() {
    let hostname = this.state.newScenarioHostname;
    if (!hostname) {
      return;
    }
    let scenarioHost = {
      Checks: [],
      Answers: [],
      Config: [],
    };
    let scenarioHosts = {
      ...this.state.scenarioHosts,
    };
    scenarioHosts[hostname] = scenarioHost;
    this.setState({
      currentScenarioHostname: hostname,
      currentScenarioHost: scenarioHost,
      newScenarioHostname: "",
      scenarioHosts: scenarioHosts,
    });
  }

  handleNewHostnameUpdate(event) {
    let hostname = event.target.value;
    this.setState({
      newScenarioHostname: hostname,
    });
  }

  render() {
    let hostOptions = [];
    hostOptions.push(<option key=""></option>);
    for (let hostname in this.state.scenarioHosts) {
      hostOptions.push(<option key={hostname}>{hostname}</option>);
    }

    let answers = this.state.currentScenarioHost.Answers || [];
    let checks = this.state.currentScenarioHost.Checks || [];
    let config = this.state.currentScenarioHost.Config || [];
    let hostname = this.state.currentScenarioHostname;

    let scenarioHosts = null;
    if (this.state.scenario.ID) {
      let scenarioHost = null;
      if (hostname) {
        scenarioHost = (
          <Fragment>
            <label>Hostname</label>
            <input onChange={this.handleHostnameRename} value={hostname} />
            <br />
            <button type="button" onClick={this.handleHostnameDelete}>
              Delete Host
            </button>
            <ScenarioHost
              answers={answers}
              checks={checks}
              config={config}
              hostname={hostname}
              parentCallback={this.handleSaveHost}
            />
          </Fragment>
        );
      }
      scenarioHosts = (
        <Fragment>
          <p>Hosts</p>
          <select onChange={this.handleHostnameSelect} value={hostname}>
            {hostOptions}
          </select>
          <br />
          <input
            onChange={this.handleNewHostnameUpdate}
            value={this.state.newScenarioHostname}
          />
          <button type="button" onClick={this.handleNewHostnameAdd}>
            +
          </button>
          <p />
          {scenarioHost}
        </Fragment>
      );
    }

    return (
      <div>
        <h1>{this.state.error}</h1>
        <form onSubmit={this.handleSave}>
          <label htmlFor="ID">ID</label>
          <input
            className="input-5"
            name="ID"
            disabled
            onChange={this.handleUpdate}
            value={this.state.scenario.ID || ""}
          />
          <br />
          <label htmlFor="Name">Name</label>
          <input
            className="input-50"
            name="Name"
            onChange={this.handleUpdate}
            value={this.state.scenario.Name || ""}
          />
          <br />
          <label htmlFor="Description">Description</label>
          <textarea
            cols="80"
            name="Description"
            onChange={this.handleUpdate}
            rows="25"
            value={this.state.scenario.Description || ""}
          />
          <br />
          <label htmlFor="Enabled">Enabled</label>
          <input
            name="Enabled"
            type="checkbox"
            onChange={this.handleUpdate}
            checked={this.state.scenario.Enabled || false}
          />
          <br />
          <button type="submit">Save</button>
          <button
            type="button"
            disabled={!this.state.scenario.ID}
            onClick={this.handleDelete}
          >
            Delete
          </button>
        </form>
        <hr />
        {scenarioHosts}
      </div>
    );
  }
}

export default withRouter(Scenario);
