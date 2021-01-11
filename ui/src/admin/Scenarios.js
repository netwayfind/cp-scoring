import "../App.css";
import { apiGet } from "../common/utils";
import LinkList from "../common/LinkList";
import Scenario from "./Scenario";

import { Component } from "react";
import { Link, Route, Switch } from "react-router-dom";
import { withRouter } from "react-router-dom/cjs/react-router-dom.min";

class Scenarios extends Component {
  constructor(props) {
    super(props);
    this.state = this.defaultState();

    this.getData = this.getData.bind(this);
  }

  componentDidMount() {
    this.getData();
  }

  defaultState() {
    return {
      error: null,
      scenarios: [],
    };
  }

  getData() {
    apiGet("/api/scenarios/").then(
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
      <div className="Scenarios">
        <Link to={this.props.match.path}>Add Scenario</Link>
        <p />
        <LinkList
          items={this.state.scenarios}
          path={this.props.match.path}
          label="Name"
        />
        <Switch>
          <Route path={`${this.props.match.url}/:id`}>
            <Scenario
              parentCallback={this.getData}
              parentPath={this.props.match.path}
            />
          </Route>
          <Route>
            <Scenario
              parentCallback={this.getData}
              parentPath={this.props.match.path}
            />
          </Route>
        </Switch>
      </div>
    );
  }
}

export default withRouter(Scenarios);
