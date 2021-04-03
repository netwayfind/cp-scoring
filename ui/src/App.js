import "./App.css";
import { apiGet, apiLogout } from "./common/utils";
import Admin from "./Admin";
import LoginUser from "./common/LoginUser";
import NotFound from "./common/NotFound";
import Insight from "./Insight";
import Scoreboard from "./Scoreboard";
import TeamDashboard from "./TeamDashboard";

import { Component, Fragment } from "react";
import { BrowserRouter, Redirect, Route, Switch } from "react-router-dom";

const basePath = "/ui";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      authenticated: false,
      username: "",
    };

    this.handleLoginSuccess = this.handleLoginSuccess.bind(this);
    this.handleLogout = this.handleLogout.bind(this);
  }

  componentDidMount() {
    let path = this.getPath();
    if (path !== "/login") {
      return;
    }
    apiGet("/api/login/").then(
      async function (s) {
        if (!s.error) {
          this.setState({
            authenticated: true,
            username: s.data,
          });
        }
      }.bind(this)
    );
  }

  getPath() {
    return window.location.pathname.replace(basePath, "");
  }

  getRedirectTo() {
    let params = new URLSearchParams(window.location.search);
    return params.get("redirectTo");
  }

  handleLoginSuccess(username) {
    this.setState({
      authenticated: true,
      username: username,
    });
  }

  handleLogout() {
    apiLogout().then(
      async function (s) {
        if (!s.error) {
          this.setState({
            authenticated: false,
          });
        }
      }.bind(this)
    );
  }

  render() {
    let destAdmin;
    let destInsight;
    let destLogin;
    let logout;
    if (this.state.authenticated) {
      // authenticated
      logout = (
        <Fragment>
          {this.state.username} -
          <button type="button" onClick={this.handleLogout}>
            Log out
          </button>
        </Fragment>
      );
      destAdmin = (
        <Fragment>
          {logout}
          <Admin />
        </Fragment>
      );
      destInsight = (
        <Fragment>
          {logout}
          <Insight />
        </Fragment>
      );
      destLogin = <Redirect to={this.getRedirectTo() || "/admin"} />;
    } else {
      // not authenticated
      destAdmin = <Redirect to="/login?redirectTo=/admin" />;
      destInsight = <Redirect to="/login?redirectTo=/insight" />;
      destLogin = <LoginUser callback={this.handleLoginSuccess} />;
    }

    return (
      <div className="App">
        <BrowserRouter basename={basePath}>
          <Switch>
            <Route exact path="/">
              <Redirect to="/scoreboard" />
            </Route>
            <Route exact path="/login">
              {destLogin}
            </Route>
            <Route path="/admin">{destAdmin}</Route>
            <Route path="/team-dashboard">
              <TeamDashboard />
            </Route>
            <Route path="/insight">{destInsight}</Route>
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
}

export default App;
