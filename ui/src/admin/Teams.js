import "../App.css";
import LinkList from "../common/LinkList";
import { apiGet } from "../common/utils";
import Team from "./Team";

import { Component, Fragment } from "react";
import { Link, Route, Switch } from "react-router-dom";
import { withRouter } from "react-router-dom/cjs/react-router-dom.min";

class Teams extends Component {
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
      teams: [],
    };
  }

  getTeamID(props) {
    return Number(
      props.location.pathname.replace(props.match.url + "/", "").split("/")[0]
    );
  }

  getData() {
    apiGet("/api/teams/").then(
      function (s) {
        this.setState({
          error: s.error,
          teams: s.data,
        });
      }.bind(this)
    );
  }

  render() {
    let teamID = this.getTeamID(this.props);
    let linkClassesAddTeam = ["nav-button"];
    if (this.props.location.pathname === this.props.match.path) {
      linkClassesAddTeam.push("nav-button-selected");
    }
    return (
      <Fragment>
        <div className="toc">
          <Link
            className={linkClassesAddTeam.join(" ")}
            to={this.props.match.path}
          >
            Add Team
          </Link>
          <p />
          <LinkList
            currentID={teamID}
            items={this.state.teams}
            path={this.props.match.path}
            label="Name"
          />
        </div>
        <div className="content">
          <Switch>
            <Route path={`${this.props.match.url}/:id`}>
              <Team
                parentCallback={this.getData}
                parentPath={this.props.match.path}
              />
            </Route>
            <Route>
              <Team
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

export default withRouter(Teams);
