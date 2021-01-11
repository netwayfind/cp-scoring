import "../App.css";
import LinkList from "../common/LinkList";
import { apiGet } from "../common/utils";
import Team from "./Team";

import { Component } from "react";
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
    return (
      <div className="Teams">
        <Link to={this.props.match.path}>Add Team</Link>
        <p />
        <LinkList
          items={this.state.teams}
          path={this.props.match.path}
          label="Name"
        />
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
    );
  }
}

export default withRouter(Teams);
