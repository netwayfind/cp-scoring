import './App.css';
import Admin from './Admin';
import LoginUser from './common/LoginUser';
import NotFound from './common/NotFound';
import Report from './Report';
import ScenarioDesc from './ScenarioDesc';
import Scoreboard from './Scoreboard';

import { Component } from 'react';
import { BrowserRouter, Redirect, Route, Switch } from 'react-router-dom';
import { apiGet, apiLogout } from './common/utils';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      authenticated: false
    }

    this.handleLoginSuccess = this.handleLoginSuccess.bind(this);
    this.handleLogout = this.handleLogout.bind(this);
  }

  componentDidMount() {
    apiGet('/api/login/')
    .then(async function(s) {
      if (!s.error) {
        this.setState({
          authenticated: true
        });
      }
    }.bind(this));
  }

  handleLoginSuccess() {
    this.setState({
      authenticated: true
    });
  }

  handleLogout() {
    apiLogout()
    .then(async function(s) {
      if (!s.error) {
        this.setState({
          authenticated: false
        });
      }
    }.bind(this));
  }

  render() {
    let destLogin;
    let destAdmin;
    let logout;
    if (this.state.authenticated) {
      // authenticated
      destLogin = (<Redirect to="/admin" />);
      destAdmin = (<Admin />);
      logout = (<button type="button" onClick={this.handleLogout}>Log out</button>)
    } else {
      // not authenticated
      destLogin = (<LoginUser callback={this.handleLoginSuccess} location="/admin" />);
      destAdmin = (<Redirect to="/login" />);
      logout = null;
    }

    return (
      <div className="App">
        <BrowserRouter basename="/ui">
          <Switch>
            <Route exact path="/">
              <Redirect to="/scoreboard" />
            </Route>
            <Route exact path="/login">
              {destLogin}
            </Route>
            <Route path="/admin">
              {logout}
              {destAdmin}
            </Route>
            <Route path="/report">
              <Report />
            </Route>
            <Route path="/scoreboard">
              <Scoreboard />
            </Route>
            <Route path="/scenarioDesc">
              <ScenarioDesc />
            </Route>
            <Route>
              <NotFound />
            </Route>
          </Switch>
        </BrowserRouter>
      </div>
    );
  }
}

export default App;
