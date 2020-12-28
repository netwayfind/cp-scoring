import '../App.css';
import { apiDelete, apiGet, apiPost, apiPut } from '../common/utils';
import ScenarioChecks from './ScenarioChecks';

import { Component } from 'react';
import { withRouter } from 'react-router-dom/cjs/react-router-dom.min';

class Scenario extends Component {
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
            scenario: {},
            checkMap: {}
        }
    }

    getData(id) {
        if (id === undefined) {
            this.setState(this.defaultState);
            return;
        }
        Promise.all([
            apiGet('/api/scenarios/' + id),
            apiGet('/api/scenarios/' + id + '/checks'),
        ])
        .then(async function(responses) {
            let s1 = responses[0];
            let s2 = responses[1];
            this.setState({
                error: s1.error || s2.error,
                scenario: s1.data,
                checkMap: s2.data
            });
        }.bind(this));
    }

    handleDelete() {
        apiDelete("/api/scenarios/" + this.state.scenario.ID)
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
        let id = this.state.scenario.ID;
        if (id) {
            // update
            apiPut("/api/scenarios/" + id, this.state.scenario)
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
            apiPost("/api/scenarios/", this.state.scenario)
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
                    <button type="button" disabled={!this.state.scenario.ID} onClick={this.handleDelete}>Delete</button>
                    <hr />
                    <p>Checks</p>
                    <ScenarioChecks scenarioID={this.state.scenario.ID} checkMap={this.state.checkMap} />
                </form>
            </div>
        );
    }
}

export default withRouter(Scenario);
