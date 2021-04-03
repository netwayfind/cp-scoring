import "../App.css";
import { apiGet } from "../common/utils";

import { Component, Fragment } from "react";
import { withRouter } from "react-router-dom";
import ReactMarkdown from "react-markdown";

class ScenarioDesc extends Component {
  constructor(props) {
    super(props);
    this.state = this.defaultState();
  }

  componentDidMount() {
    let id = this.props.match.params.scenarioID;
    this.getData(id);
  }

  componentDidUpdate(prevProps) {
    let id = this.props.match.params.scenarioID;
    let prevId = prevProps.match.params.scenarioID;
    if (id !== prevId) {
      this.getData(id);
    }
  }

  defaultState() {
    return {
      error: null,
      scenario: {},
    };
  }

  getData(id) {
    apiGet("/api/scenario-desc/" + id).then(
      function (s) {
        this.setState({
          error: s.error,
          scenario: s.data,
        });
      }.bind(this)
    );
  }

  render() {
    return (
      <Fragment>
        <div className="heading">
          <h1>{this.state.scenario.Name}</h1>
        </div>
        <div>
          <ReactMarkdown children={this.state.scenario.Description} />
        </div>
      </Fragment>
    );
  }
}

export default withRouter(ScenarioDesc);
