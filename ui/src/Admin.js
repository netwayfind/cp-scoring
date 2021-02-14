import "./App.css";
import NotFound from "./common/NotFound";
import Scenarios from "./admin/Scenarios";
import Teams from "./admin/Teams";
import Users from "./admin/Users";

import { Component, Fragment } from "react";
import { Link, Route, Switch } from "react-router-dom";
import { withRouter } from "react-router-dom/cjs/react-router-dom.min";

class Admin extends Component {
  render() {
    let url = this.props.match.url;
    let path = this.props.match.path;

    let linkClassesTeams = ["nav-button"];
    if (this.props.location.pathname.startsWith(`${url}/teams`)) {
      linkClassesTeams.push(["nav-button-selected"]);
    }
    let linkClassesScenarios = ["nav-button"];
    if (this.props.location.pathname.startsWith(`${url}/scenarios`)) {
      linkClassesScenarios.push(["nav-button-selected"]);
    }
    let linkClassesUsers = ["nav-button"];
    if (this.props.location.pathname.startsWith(`${url}/users`)) {
      linkClassesUsers.push(["nav-button-selected"]);
    }

    return (
      <Fragment>
        <div className="navbar">
          <Link className="nav-button" target="_blank" to="/insight">
            Insight
          </Link>
        </div>
        <div className="heading">
          <h1>cp-scoring admin</h1>
        </div>
        <div className="navbar">
          <Link className={linkClassesTeams.join(" ")} to={`${url}/teams`}>
            Teams
          </Link>
          <Link
            className={linkClassesScenarios.join(" ")}
            to={`${url}/scenarios`}
          >
            Scenarios
          </Link>
          <Link className={linkClassesUsers.join(" ")} to={`${url}/users`}>
            Users
          </Link>
        </div>
        <div>
          <Switch>
            <Route path={`${path}/teams`}>
              <Teams />
            </Route>
            <Route path={`${path}/scenarios`}>
              <Scenarios />
            </Route>
            <Route path={`${path}/users`}>
              <Users />
            </Route>
            <Route>
              <NotFound />
            </Route>
          </Switch>
        </div>
      </Fragment>
    );
  }
}

export default withRouter(Admin);
