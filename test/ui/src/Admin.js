import './App.css';
import NotFound from './NotFound';
import Scenarios from './admin/Scenarios';
import Teams from './admin/Teams';

import {Link, Route, Switch, useRouteMatch} from 'react-router-dom';

export default function Admin() {
  let { path, url } = useRouteMatch();
  
  return (
    <div className="Admin">
      admin
      <ul>
        <li><Link to={`${url}/teams`}>Teams</Link></li>
        <li><Link to={`${url}/scenarios`}>Scenarios</Link></li>
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
        <Route>
          <NotFound />
        </Route>
      </Switch>
    </div>
  );
}
