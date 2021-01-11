import "./App.css";
import { apiGet } from "./common/utils";

import { Component } from "react";
import { withRouter } from "react-router-dom/cjs/react-router-dom.min";
import ReactMarkdown from "react-markdown";

class ScenarioDesc extends Component {
  constructor(props) {
    super(props);
    this.state = {
      scenario: {},
    };
  }

  componentDidMount() {
    let id = 1;
    this.getData(id);
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
      <div className="ScenarioDesc">
        <h1>{this.state.scenario.Name}</h1>
        <ReactMarkdown children={this.state.scenario.Description} />
      </div>
    );
  }
}

export default withRouter(ScenarioDesc);
