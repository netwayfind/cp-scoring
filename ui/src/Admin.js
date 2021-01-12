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

    return (
      <Fragment>
        <div className="heading">
          <h1>cp-scoring admin</h1>
        </div>
        <div className="navbar">
          <Link className="nav-button" to={`${url}/teams`}>
            Teams
          </Link>
          <Link className="nav-button" to={`${url}/scenarios`}>
            Scenarios
          </Link>
          <Link className="nav-button" to={`${url}/users`}>
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
