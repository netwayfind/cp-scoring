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
        this.handleSaveChecks = this.handleSaveChecks.bind(this);
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
            checkMap: {},
            answerMap: {}
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
            apiGet('/api/scenarios/' + id + '/answers'),
        ])
        .then(async function(responses) {
            let s1 = responses[0];
            let s2 = responses[1];
            let s3 = responses[2];
            this.setState({
                error: s1.error || s2.error || s3.error,
                scenario: s1.data,
                checkMap: s2.data,
                answerMap: s3.data
            });
        }.bind(this));
    }

    handleDelete() {
        apiDelete('/api/scenarios/' + this.state.scenario.ID)
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
            apiPut('/api/scenarios/' + id, this.state.scenario)
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
            apiPost('/api/scenarios/', this.state.scenario)
            .then(async function(s) {
                if (s.error) {
                    this.setState({
                        error: s.error
                    });
                } else {
                    this.props.parentCallback();
                    this.props.history.push(this.props.match.url + '/' + s.data.ID);
                }
            }.bind(this));
        }
    }

    handleSaveChecks(checkMap, answerMap) {
        let id = this.state.scenario.ID;
        Promise.all([
            apiPut('/api/scenarios/' + id + '/checks', checkMap),
            apiPut('/api/scenarios/' + id + '/answers', answerMap)
        ])
        .then(async function(responses) {
            let s1 = responses[0];
            let s2 = responses[1];
            let error = s1.error || s2.error;
            if (error) {
                this.setState({
                    error: error
                });
            } else {
                this.props.parentCallback();
                this.props.history.push(this.props.match.url);
            }
            this.setState({
                error: s1.error || s2.error
            });
        }.bind(this));
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
                    <label htmlFor="Name">Name</label>
                    <input onChange={this.handleUpdate} name="Name" value={this.state.scenario.Name || ""} />
                    <label htmlFor="Description">Description</label>
                    <input onChange={this.handleUpdate} name="Description" value={this.state.scenario.Description || ""} />
                    <label htmlFor="Enabled">Enabled</label>
                    <input onChange={this.handleUpdate} name="Enabled" type="checkbox" value={this.state.scenario.Enabled || false} />
                    <button type="submit">Save</button>
                    <button type="button" disabled={!this.state.scenario.ID} onClick={this.handleDelete}>Delete</button>
                </form>
                <hr />
                <p>Checks</p>
                <ScenarioChecks scenarioID={this.state.scenario.ID} checkMap={this.state.checkMap} answerMap={this.state.answerMap} parentCallback={this.handleSaveChecks} />
            </div>
        );
    }
}

export default withRouter(Scenario);
