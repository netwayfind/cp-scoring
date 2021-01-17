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

  getScenarioID(props) {
    return Number(
      props.location.pathname.replace(props.match.url + "/", "").split("/")[0]
    );
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
    let scenarioID = this.getScenarioID(this.props);
    return (
      <Fragment>
        <div className="heading">
          <h1>Scoreboard</h1>
        </div>
        <div className="toc">
          <h4>Scenarios</h4>
          <LinkList
            currentID={scenarioID}
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
