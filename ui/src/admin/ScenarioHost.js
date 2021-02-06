import "../App.css";

import { Component, Fragment } from "react";
import { withRouter } from "react-router-dom";

class ScenarioHost extends Component {
  constructor(props) {
    super(props);
    this.state = {
      answers: props.answers,
      checks: props.checks,
      config: props.config,
      hostname: props.hostname,
    };

    this.handleAnswerUpdate = this.handleAnswerUpdate.bind(this);
    this.handleCheckAdd = this.handleCheckAdd.bind(this);
    this.handleCheckDelete = this.handleCheckDelete.bind(this);
    this.handleCheckReorder = this.handleCheckReorder.bind(this);
    this.handleCheckUpdate = this.handleCheckUpdate.bind(this);
    this.handleCheckArgAdd = this.handleCheckArgAdd.bind(this);
    this.handleCheckArgDelete = this.handleCheckArgDelete.bind(this);
    this.handleCheckArgUpdate = this.handleCheckArgUpdate.bind(this);
    this.handleConfigAdd = this.handleConfigAdd.bind(this);
    this.handleConfigDelete = this.handleConfigDelete.bind(this);
    this.handleConfigUpdate = this.handleConfigUpdate.bind(this);
    this.handleConfigReorder = this.handleConfigReorder.bind(this);
    this.handleConfigArgAdd = this.handleConfigArgAdd.bind(this);
    this.handleConfigArgDelete = this.handleConfigArgDelete.bind(this);
    this.handleConfigArgUpdate = this.handleConfigArgUpdate.bind(this);
    this.handleSave = this.handleSave.bind(this);
  }

  componentDidUpdate(prevProps) {
    if (this.props.hostname !== prevProps.hostname) {
      this.setState({
        answers: this.props.answers,
        checks: this.props.checks,
        config: this.props.config,
        hostname: this.props.hostname,
      });
    }
  }

  handleAnswerUpdate(i, event) {
    let name = event.target.name;
    let value = event.target.value;
    let answers = [...this.state.answers];
    if (event.target.type === "number") {
      value = Number(value);
    }
    answers[i][name] = value;
    this.setState({
      answers: answers,
    });
  }

  handleCheckAdd() {
    let answers = [...this.state.answers];
    let checks = [...this.state.checks];
    answers.push({
      Type: "",
      Value: "",
    });
    checks.push({
      Type: "EXEC",
      Command: "",
      Args: [],
    });
    this.setState({
      answers: answers,
      checks: checks,
    });
  }

  handleCheckDelete(i) {
    let answers = [...this.state.answers];
    let checks = [...this.state.checks];
    answers.splice(i, 1);
    checks.splice(i, 1);
    this.setState({
      answers: answers,
      checks: checks,
    });
  }

  handleCheckReorder(event, currentIndex) {
    let newIndex = event.target.value;
    let answers = [...this.state.answers];
    let checks = [...this.state.checks];
    answers.splice(newIndex, 0, answers.splice(currentIndex, 1)[0]);
    checks.splice(newIndex, 0, checks.splice(currentIndex, 1)[0]);
    this.setState({
      answers: answers,
      checks: checks,
    });
  }

  handleCheckUpdate(i, event) {
    let name = event.target.name;
    let value = event.target.value;
    let checks = [...this.state.checks];
    checks[i][name] = value;
    this.setState({
      checks: checks,
    });
  }

  handleCheckArgAdd(i) {
    let checks = [...this.state.checks];
    checks[i]["Args"].push("");
    this.setState({
      checks: checks,
    });
  }

  handleCheckArgDelete(i, j) {
    let checks = [...this.state.checks];
    checks[i]["Args"].splice(j, 1);
    this.setState({
      checks: checks,
    });
  }

  handleCheckArgUpdate(i, j, event) {
    let value = event.target.value;
    let checks = [...this.state.checks];
    checks[i]["Args"][j] = value;
    this.setState({
      checks: checks,
    });
  }

  handleConfigAdd() {
    let config = [...this.state.config];
    config.push({
      Type: "EXEC",
      Command: "",
      Args: [],
    });
    this.setState({
      config: config,
    });
  }

  handleConfigDelete(i) {
    let config = [...this.state.config];
    config.splice(i, 1);
    this.setState({
      config: config,
    });
  }

  handleConfigReorder(event, currentIndex) {
    let newIndex = event.target.value;
    let config = [...this.state.config];
    config.splice(newIndex, 0, config.splice(currentIndex, 1)[0]);
    this.setState({
      config: config,
    });
  }

  handleConfigUpdate(i, event) {
    let name = event.target.name;
    let value = event.target.value;
    let config = [...this.state.config];
    config[i][name] = value;
    this.setState({
      config: config,
    });
  }

  handleConfigArgAdd(i) {
    let config = [...this.state.config];
    config[i]["Args"].push("");
    this.setState({
      config: config,
    });
  }

  handleConfigArgDelete(i, j) {
    let config = [...this.state.config];
    config[i]["Args"].splice(j, 1);
    this.setState({
      config: config,
    });
  }

  handleConfigArgUpdate(i, j, event) {
    let value = event.target.value;
    let config = [...this.state.config];
    config[i]["Args"][j] = value;
    this.setState({
      config: config,
    });
  }

  handleSave(event) {
    if (event !== null) {
      event.preventDefault();
    }
    this.props.parentCallback(
      this.state.checks,
      this.state.answers,
      this.state.config
    );
  }

  render() {
    let actionOptions = [
      <option key="1">EXEC</option>,
      <option key="2">FILE_EXIST</option>,
      <option key="2">FILE_REGEX</option>,
      <option key="2">FILE_VALUE</option>,
    ];
    let operatorOptions = [
      <option key="1" value="" />,
      <option key="2">EQUAL</option>,
      <option key="3">NOT_EQUAL</option>,
    ];

    let checkList = [];
    let checks = this.state.checks;
    let checksPositionOptions = [];
    for (let i in checks) {
      checksPositionOptions.push(
        <option key={i} value={i}>
          {Number(i) + 1}
        </option>
      );
    }
    checks.forEach((check, i) => {
      let args = [];
      if (check.Args) {
        check.Args.forEach((arg, j) => {
          args.push(
            <li key={j}>
              <input
                className="input-50"
                onChange={(event) => this.handleCheckArgUpdate(i, j, event)}
                value={arg}
              ></input>
              <button
                type="button"
                onClick={() => this.handleCheckArgDelete(i, j)}
              >
                -
              </button>
            </li>
          );
        });
      }
      args.push(
        <li key="arg_add">
          <button type="button" onClick={() => this.handleCheckArgAdd(i)}>
            Add Arg
          </button>
        </li>
      );
      let answer = this.state.answers[i];
      checkList.push(
        <li key={i}>
          <select
            onChange={(event) => this.handleCheckReorder(event, i)}
            value={i}
          >
            {checksPositionOptions}
          </select>
          <details>
            <summary>{answer.Description}</summary>
            <button type="button" onClick={() => this.handleCheckDelete(i)}>
              Delete Check
            </button>
            <p />
            <label htmlFor="Description">Description</label>
            <input
              className="input-20"
              name="Description"
              onChange={(event) => this.handleAnswerUpdate(i, event)}
              value={answer.Description}
            />
            <br />
            <label htmlFor="Type">Type</label>
            <select
              name="Type"
              onChange={(event) => this.handleCheckUpdate(i, event)}
              value={check.Type}
            >
              {actionOptions}
            </select>
            <br />
            {check.Type === "EXEC" ? (
              <Fragment>
                <label htmlFor="Command">Command</label>
                <input
                  className="input-50"
                  name="Command"
                  onChange={(event) => this.handleCheckUpdate(i, event)}
                  value={check.Command}
                />
                <br />
              </Fragment>
            ) : null}
            <label htmlFor="Args">Args</label>
            <ul>{args}</ul>
            <label htmlFor="Answer">Answer</label>
            <select
              name="Operator"
              onChange={(event) => this.handleAnswerUpdate(i, event)}
              value={answer.Operator}
            >
              {operatorOptions}
            </select>
            <input
              className="input-50"
              name="Value"
              onChange={(event) => this.handleAnswerUpdate(i, event)}
              value={answer.Value}
            />
            <br />
            <label htmlFor="Points">Points</label>
            <input
              className="input-5"
              name="Points"
              onChange={(event) => this.handleAnswerUpdate(i, event)}
              value={answer.Points}
              type="number"
              steps="1"
            />
          </details>
        </li>
      );
    });
    checkList.push(
      <li key="check_add">
        <button type="button" onClick={this.handleCheckAdd}>
          Add Check
        </button>
      </li>
    );

    let configList = [];
    let config = this.state.config;
    let configPositionOptions = [];
    for (let i in config) {
      configPositionOptions.push(
        <option key={i} value={i}>
          {Number(i) + 1}
        </option>
      );
    }
    config.forEach((conf, i) => {
      let args = [];
      if (conf.Args) {
        conf.Args.forEach((arg, j) => {
          args.push(
            <li key={j}>
              <input
                className="input-50"
                onChange={(event) => this.handleConfigArgUpdate(i, j, event)}
                value={arg}
              ></input>
              <button
                type="button"
                onClick={() => this.handleConfigArgDelete(i, j)}
              >
                -
              </button>
            </li>
          );
        });
      }
      args.push(
        <li key="arg_add">
          <button type="button" onClick={() => this.handleConfigArgAdd(i)}>
            Add Arg
          </button>
        </li>
      );
      configList.push(
        <li key={i}>
          <select
            onChange={(event) => this.handleConfigReorder(event, i)}
            value={i}
          >
            {configPositionOptions}
          </select>
          <details>
            <summary>
              Command: {conf.Command}, Args: [{conf.Args.join(" ") || ""}]
            </summary>
            <button type="button" onClick={() => this.handleConfigDelete(i)}>
              Delete Config
            </button>
            <p />
            <label htmlFor="Type">Type</label>
            <select disabled name="Type" value="EXEC">
              {actionOptions}
            </select>
            <br />
            <label htmlFor="Command">Command</label>
            <input
              className="input-50"
              name="Command"
              onChange={(event) => this.handleConfigUpdate(i, event)}
              value={conf.Command}
            />
            <br />
            <label htmlFor="Args">Args</label>
            <ul>{args}</ul>
          </details>
        </li>
      );
    });
    configList.push(
      <li key="config_add">
        <button type="button" onClick={this.handleConfigAdd}>
          Add Config
        </button>
      </li>
    );

    return (
      <form onSubmit={this.handleSave}>
        <p>Checks</p>
        <ol>{checkList}</ol>
        <p>Config</p>
        <ol>{configList}</ol>
        <button type="submit">Save Host</button>
      </form>
    );
  }
}

export default withRouter(ScenarioHost);
