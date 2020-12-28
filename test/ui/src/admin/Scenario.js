import '../App.css';
import { apiGet, apiPost, apiPut } from '../common/utils';

import { Component } from 'react';
import { withRouter } from 'react-router-dom/cjs/react-router-dom.min';

class Scenario extends Component {
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
            scenario: {}
        }
    }

    getData(id) {
        if (id === undefined) {
            this.setState(this.defaultState);
            return;
        }
        apiGet('/api/scenarios/' + id)
        .then(async function(s) {
            this.setState({
                error: s.error,
                scenario: s.data
            });
        }.bind(this));
    }

    handleSave(event) {
        if (event !== null) {
          event.preventDefault();
        }
        let id = this.state.scenario.ID;
        if (id) {
            // update
            apiPut("/api/scenarios/" + id, this.state.scenario)
            .then(async function(s) {
                this.props.callback();
                let url = this.props.match.url;
                this.props.history.push(url);
            }.bind(this));
        } else {
            // create
            apiPost("/api/scenarios/", this.state.scenario)
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
            scenario: {
                ...this.state.scenario,
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
                    <input onChange={this.handleUpdate} name="ID" disabled value={this.state.scenario.ID || ""} />
                    <label htmlFor="ID">Name</label>
                    <input onChange={this.handleUpdate} name="Name" value={this.state.scenario.Name || ""} />
                    <label htmlFor="ID">Description</label>
                    <input onChange={this.handleUpdate} name="Description" value={this.state.scenario.Description || ""} />
                    <label htmlFor="ID">Enabled</label>
                    <input onChange={this.handleUpdate} name="Enabled" type="checkbox" value={this.state.scenario.Enabled || false} />
                    <button type="submit">Save</button>
                </form>
            </div>
        );
    }
}

export default withRouter(Scenario);
