import '../App.css';
import Team from './Team';

import {Link, Route, Switch, useRouteMatch} from 'react-router-dom';

export default function Teams() {
    let { path, url } = useRouteMatch();

    return (
        <div className="Teams">
            <ul>
                <li><Link to={`${url}/1`}>1</Link></li>
                <li><Link to={`${url}/2`}>2</Link></li>
            </ul>
            <Switch>
                <Route path={`${path}/:id`}>
                    <Team />
                </Route>
            </Switch>
        </div>
    );
}
