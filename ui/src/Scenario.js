import "./App.css";
import ScenarioDesc from "./scenario/ScenarioDesc";

import { Component } from "react";
import {
  Route,
  Switch,
  withRouter,
} from "react-router-dom/cjs/react-router-dom.min";

class Scenario extends Component {
  render() {
    return (
      <Switch>
        <Route path={`${this.props.match.url}/:id`}>
          <ScenarioDesc />
        </Route>
        <Route>
          <p>No scenario selected</p>
        </Route>
      </Switch>
    );
  }
}

export default withRouter(Scenario);
