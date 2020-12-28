import '../App.css';
import { apiGet, apiPost, apiPut } from '../common/utils';

import { Component } from 'react';
import { withRouter } from 'react-router-dom/cjs/react-router-dom.min';

class Team extends Component {
    constructor(props) {
        super(props);
        this.state = this.defaultState();

        this.getData = this.getData.bind(this);
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
            team: {}
        }
    }

    getData(id) {
        if (id === undefined) {
            this.setState(this.defaultState);
            return;
        }
        apiGet('/api/teams/' + id)
        .then(async function(s) {
            this.setState({
                error: s.error,
                team: s.data
            });
        }.bind(this));
    }

    handleSave(event) {
        if (event !== null) {
          event.preventDefault();
        }
        let id = this.state.team.ID;
        if (id) {
            // update
            apiPut("/api/teams/" + id, this.state.team)
            .then(async function(s) {
                this.props.callback();
                let url = this.props.match.url;
                this.props.history.push(url);
            }.bind(this));
        } else {
            // create
            apiPost("/api/teams/", this.state.team)
            .then(async function(s) {
                this.props.callback();
                let url = this.props.match.url + "/" + s.data.ID;
                this.props.history.push(url);
            }.bind(this));
        }
    }

    handleUpdate(event) {
        let value = event.target.value;
        if (event.target.type === 'checkbox') {
          value = event.target.checked;
        }
        this.setState({
            team: {
                ...this.state.team,
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
                    <input onChange={this.handleUpdate} name="ID" disabled value={this.state.team.ID || ""} />
                    <label htmlFor="ID">Name</label>
                    <input onChange={this.handleUpdate} name="Name" value={this.state.team.Name || ""} />
                    <label htmlFor="ID">POC</label>
                    <input onChange={this.handleUpdate} name="POC" value={this.state.team.POC || ""} />
                    <label htmlFor="ID">Email</label>
                    <input onChange={this.handleUpdate} name="Email" value={this.state.team.Email || ""} />
                    <label htmlFor="ID">Enabled</label>
                    <input onChange={this.handleUpdate} name="Enabled" type="checkbox" value={this.state.team.Enabled || false} />
                    <label htmlFor="ID">Key</label>
                    <input onChange={this.handleUpdate} name="Key" value={this.state.team.Key || ""} />
                    <button type="submit">Save</button>
                </form>
            </div>
        );
    }
}

export default withRouter(Team);
