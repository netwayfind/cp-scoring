import "./App.css";
import NotFound from "./common/NotFound";
import Scenarios from "./admin/Scenarios";
import Teams from "./admin/Teams";
import Users from "./admin/Users";

import { Link, Route, Switch, useRouteMatch } from "react-router-dom";

export default function Admin() {
  let { path, url } = useRouteMatch();

  return (
    <div className="Admin">
      admin
      <ul>
        <li>
          <Link to={`${url}/teams`}>Teams</Link>
        </li>
        <li>
          <Link to={`${url}/scenarios`}>Scenarios</Link>
        </li>
        <li>
          <Link to={`${url}/users`}>Users</Link>
        </li>
      </ul>
      <Switch>
        <Route exact path={path}>
          ????
        </Route>
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
  );
}
