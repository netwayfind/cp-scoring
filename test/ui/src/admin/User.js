import '../App.css';
import { apiDelete, apiGet, apiPost, apiPut } from '../common/utils';

import { Component } from 'react';
import { withRouter } from 'react-router-dom/cjs/react-router-dom.min';

class User extends Component {
    constructor(props) {
        super(props);
        this.state = this.defaultState();

        this.getData = this.getData.bind(this);
        this.handleDelete = this.handleDelete.bind(this);
        this.handleUpdate = this.handleUpdate.bind(this);
        this.handleSave = this.handleSave.bind(this);
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
            user: {}
        }
    }

    getData(id) {
        if (id === undefined) {
            this.setState(this.defaultState);
            return;
        }
        apiGet('/api/users/' + id)
        .then(async function(s) {
            this.setState({
                error: s.error,
                user: s.data
            });
        }.bind(this));
    }

    handleDelete() {
        apiDelete("/api/users/" + this.state.user.ID)
        .then(async function(s) {
            if (s.error) {
                this.setState({
                    error: s.error
                });
            } else {
                this.props.parentCallback();
                this.props.history.push(this.props.parentPath);
            }
        }.bind(this));
    }

    handleSave(event) {
        if (event !== null) {
          event.preventDefault();
        }
        let id = this.state.user.ID;
        if (id) {
            // update
            apiPut("/api/users/" + id, this.state.user)
            .then(async function(s) {
                if (s.error) {
                    this.setState({
                        error: s.error
                    });
                } else {
                    this.props.parentCallback();
                    this.props.history.push(this.props.match.url);
                }
            }.bind(this));
        } else {
            // create
            apiPost("/api/users/", this.state.user)
            .then(async function(s) {
                if (s.error) {
                    this.setState({
                        error: s.error
                    });
                } else {
                    this.props.parentCallback();
                    this.props.history.push(this.props.match.url + "/" + s.data.ID);
                }
            }.bind(this));
        }
    }

    handleUpdate(event) {
        let value = event.target.value;
        if (event.target.type === 'checkbox') {
          value = event.target.checked;
        }
        this.setState({
            user: {
                ...this.state.user,
                [event.target.name]: value
            }
        });
    }

    render() {
        return (
            <div>
                <h1>{this.state.error}</h1>
                <form onSubmit={this.handleSave}>
                    <label htmlFor="ID">ID</label>
                    <input onChange={this.handleUpdate} name="ID" disabled value={this.state.user.ID || ""} />
                    <label htmlFor="ID">Username</label>
                    <input onChange={this.handleUpdate} name="Username" value={this.state.user.Username || ""} />
                    <label htmlFor="ID">Password</label>
                    <input type="password" onChange={this.handleUpdate} name="Password" value={this.state.user.Password || ""} />
                    <label htmlFor="ID">Email</label>
                    <input onChange={this.handleUpdate} name="Email" value={this.state.user.Email || ""} />
                    <label htmlFor="ID">Enabled</label>
                    <input onChange={this.handleUpdate} name="Enabled" type="checkbox" value={this.state.user.Enabled || false} />
                    <button type="submit">Save</button>
                    <button type="button" disabled={!this.state.user.ID} onClick={this.handleDelete}>Delete</button>
                </form>
            </div>
        );
    }
}

export default withRouter(User);
