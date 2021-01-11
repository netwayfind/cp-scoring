import "./App.css";
import { apiGet } from "./common/utils";
import LinkList from "./common/LinkList";
import ScenarioScoreboard from "./scenario/ScenarioScoreboard";

import { Component } from "react";
import { Route, Switch, withRouter } from "react-router-dom";

class Scoreboard extends Component {
  constructor(props) {
    super(props);
    this.state = {
      scenarios: [],
    };
  }

  componentDidMount() {
    this.getData();
  }

  getData() {
    apiGet("/api/scoreboard/scenarios").then(
      function (s) {
        this.setState({
          error: s.error,
          scenarios: s.data,
        });
      }.bind(this)
    );
  }

  render() {
    return (
      <div className="Scoreboard">
        <LinkList items={this.state.scenarios} path={this.props.match.path} />
        <Switch>
          <Route path={`${this.props.match.url}/:id`}>
            <ScenarioScoreboard parentPath={this.props.match.path} />
          </Route>
        </Switch>
      </div>
    );
  }
}

export default withRouter(Scoreboard);
