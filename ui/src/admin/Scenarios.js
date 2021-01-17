import "../App.css";
import { apiGet } from "../common/utils";
import LinkList from "../common/LinkList";
import Scenario from "./Scenario";

import { Component, Fragment } from "react";
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

  getScenarioID(props) {
    return Number(
      props.location.pathname.replace(props.match.url + "/", "").split("/")[0]
    );
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
    let scenarioID = this.getScenarioID(this.props);
    let linkClassesAddScenario = ["nav-button"];
    if (this.props.location.pathname === this.props.match.path) {
      linkClassesAddScenario.push("nav-button-selected");
    }
    return (
      <Fragment>
        <div className="toc">
          <Link
            className={linkClassesAddScenario.join(" ")}
            to={this.props.match.path}
          >
            Add Scenario
          </Link>
          <p />
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
      </Fragment>
    );
  }
}

export default withRouter(Scenarios);
