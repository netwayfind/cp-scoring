import "../App.css";
import LinkList from "../common/LinkList";
import { apiGet } from "../common/utils";
import User from "./User";

import { Component, Fragment } from "react";
import { Link, Route, Switch } from "react-router-dom";
import { withRouter } from "react-router-dom/cjs/react-router-dom.min";

class Users extends Component {
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
      users: [],
    };
  }

  getUserID(props) {
    return Number(
      props.location.pathname.replace(props.match.url + "/", "").split("/")[0]
    );
  }

  getData() {
    apiGet("/api/users/").then(
      function (s) {
        this.setState({
          error: s.error,
          users: s.data,
        });
      }.bind(this)
    );
  }

  render() {
    let pathNewUser = this.props.match.path + "/new";
    let userID = this.getUserID(this.props);
    let linkClassesAddUser = ["nav-button"];
    if (this.props.location.pathname === pathNewUser) {
      linkClassesAddUser.push("nav-button-selected");
    }
    return (
      <Fragment>
        <div className="toc">
          <Link className={linkClassesAddUser.join(" ")} to={pathNewUser}>
            Add User
          </Link>
          <p />
          <LinkList
            currentID={userID}
            items={this.state.users}
            path={this.props.match.path}
            label="Username"
          />
        </div>
        <div className="content">
          <Switch>
            <Route path={pathNewUser}>
              <User
                parentCallback={this.getData}
                parentPath={this.props.match.path}
              />
            </Route>
            <Route path={`${this.props.match.url}/:id`}>
              <User
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

export default withRouter(Users);
