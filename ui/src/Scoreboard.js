import "./App.css";
import { apiGet } from "./common/utils";
import LinkList from "./common/LinkList";
import ScenarioScoreboard from "./scoreboard/ScenarioScoreboard";

import { Component, Fragment } from "react";
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
      <Fragment>
        <div className="heading">
          <h1>Scoreboard</h1>
        </div>
        <div className="toc">
          <h4>Scenarios</h4>
          <LinkList
            items={this.state.scenarios}
            path={this.props.match.path}
            label="Name"
          />
        </div>
        <div className="content">
          <Switch>
            <Route path={`${this.props.match.url}/:id`}>
              <ScenarioScoreboard parentPath={this.props.match.path} />
            </Route>
          </Switch>
        </div>
      </Fragment>
    );
  }
}

export default withRouter(Scoreboard);
