import "../App.css";
import { apiDelete, apiGet, apiPost, apiPut } from "../common/utils";

import { Component, Fragment } from "react";
import { withRouter } from "react-router-dom/cjs/react-router-dom.min";

class User extends Component {
  constructor(props) {
    super(props);
    this.state = this.defaultState();

    this.getData = this.getData.bind(this);
    this.handleDelete = this.handleDelete.bind(this);
    this.handleUpdate = this.handleUpdate.bind(this);
    this.handleSave = this.handleSave.bind(this);
    this.handleRoleAdd = this.handleRoleAdd.bind(this);
    this.handleRoleAddSelection = this.handleRoleAddSelection.bind(this);
    this.handleRoleDelete = this.handleRoleDelete.bind(this);
  }

  componentDidMount() {
    let id = this.props.match.params.id;
    this.getData(id);
  }

  componentDidUpdate(prevProps) {
    let id = this.props.match.params.id;
    let prevId = prevProps.match.params.id;
    if (id !== prevId) {
      this.getData(id);
    }
  }

  defaultState() {
    return {
      error: null,
      selectedRole: "",
      roles: [],
      user: {},
    };
  }

  getData(id) {
    if (id === undefined) {
      this.setState(this.defaultState);
      return;
    }
    Promise.all([
      apiGet("/api/users/" + id),
      apiGet("/api/users/" + id + "/roles"),
    ]).then(
      async function (responses) {
        let s1 = responses[0];
        let s2 = responses[1];
        this.setState({
          error: s1.error || s2.error,
          selectedRole: "",
          roles: s2.data,
          user: s1.data,
        });
      }.bind(this)
    );
  }

  handleDelete() {
    apiDelete("/api/users/" + this.state.user.ID).then(
      async function (s) {
        if (s.error) {
          this.setState({
            error: s.error,
          });
        } else {
          this.props.parentCallback();
          this.props.history.push(this.props.parentPath);
        }
      }.bind(this)
    );
  }

  handleSave(event) {
    if (event !== null) {
      event.preventDefault();
    }
    let id = this.state.user.ID;
    if (id) {
      // update
      Promise.all([
        apiPut("/api/users/" + id, this.state.user),
        apiPut("/api/users/" + id + "/roles", this.state.roles),
      ]).then(
        async function (responses) {
          let s1 = responses[0];
          let s2 = responses[1];
          if (s1.error || s2.error) {
            this.setState({
              error: s1.error || s2.error,
            });
          } else {
            this.props.parentCallback();
            this.props.history.push(this.props.match.url);
          }
        }.bind(this)
      );
    } else {
      // create
      apiPost("/api/users/", this.state.user).then(
        async function (s) {
          if (s.error) {
            this.setState({
              error: s.error,
            });
          } else {
            this.props.parentCallback();
            this.props.history.push(this.props.match.url + "/" + s.data.ID);
          }
        }.bind(this)
      );
    }
  }

  handleUpdate(event) {
    let value = event.target.value;
    if (event.target.type === "checkbox") {
      value = event.target.checked;
    }
    this.setState({
      user: {
        ...this.state.user,
        [event.target.name]: value,
      },
    });
  }

  handleRoleAdd() {
    if (!this.state.selectedRole) {
      return;
    }
    let roles = [...this.state.roles, this.state.selectedRole];
    this.setState({
      selectedRole: "",
      roles: roles,
    });
  }

  handleRoleAddSelection(event) {
    let value = event.target.value;
    if (!value) {
      return;
    }
    this.setState({
      selectedRole: value,
    });
  }

  handleRoleDelete(i) {
    let roles = [...this.state.roles];
    roles.splice(i, 1);
    this.setState({
      roles: roles,
    });
  }

  render() {
    let roles = ["ADMIN"];
    let roleOptions = [<option key="1" value="" />];
    let roleList = [];
    let rolesField = null;
    if (this.state.user.ID) {
      this.state.roles.forEach((role, i) => {
        roleList.push(
          <li key={i}>
            {role}
            <button type="button" onClick={() => this.handleRoleDelete(i)}>
              -
            </button>
          </li>
        );
        let j = roles.indexOf(role);
        if (j > -1) {
          roles.splice(j, 1);
        }
      });
      roles.forEach((role) => {
        roleOptions.push(<option key={roleOptions.length + 1}>{role}</option>);
      });
      roleList.push(
        <li key="new_role">
          <select
            onChange={this.handleRoleAddSelection}
            value={this.state.selectedRole}
          >
            {roleOptions}
          </select>
          <button onClick={this.handleRoleAdd} type="button">
            +
          </button>
        </li>
      );
      rolesField = (
        <Fragment>
          <p>Roles</p>
          <ul>{roleList}</ul>
        </Fragment>
      );
    }
    return (
      <div>
        <h1>{this.state.error}</h1>
        <form onSubmit={this.handleSave}>
          <label htmlFor="ID">ID</label>
          <input
            className="input-5"
            onChange={this.handleUpdate}
            name="ID"
            disabled
            value={this.state.user.ID || ""}
          />
          <br />
          <label htmlFor="ID">Username</label>
          <input
            onChange={this.handleUpdate}
            name="Username"
            value={this.state.user.Username || ""}
          />
          <br />
          <label htmlFor="ID">Password</label>
          <input
            type="password"
            onChange={this.handleUpdate}
            name="Password"
            value={this.state.user.Password || ""}
          />
          <br />
          <label htmlFor="ID">Email</label>
          <input
            onChange={this.handleUpdate}
            name="Email"
            value={this.state.user.Email || ""}
          />
          <br />
          <label htmlFor="ID">Enabled</label>
          <input
            onChange={this.handleUpdate}
            name="Enabled"
            type="checkbox"
            checked={this.state.user.Enabled || false}
          />
          <br />
          {rolesField}
          <button type="submit">Save</button>
          <button
            type="button"
            disabled={!this.state.user.ID}
            onClick={this.handleDelete}
          >
            Delete
          </button>
        </form>
      </div>
    );
  }
}

export default withRouter(User);
