import './App.css';
import Admin from './Admin';
import NotFound from './common/NotFound';
import Scoreboard from './Scoreboard';

import { BrowserRouter, Redirect, Route, Switch } from 'react-router-dom';

export default function App() {
  return (
    <div className="App">
      <BrowserRouter basename="/ui">
        <Switch>
          <Route exact path="/">
            <Redirect to="/scoreboard" />
          </Route>
          <Route path="/admin">
            <Admin />
          </Route>
          <Route path="/scoreboard">
            <Scoreboard />
          </Route>
          <Route>
            <NotFound />
          </Route>
        </Switch>
      </BrowserRouter>
    </div>
  );
}
