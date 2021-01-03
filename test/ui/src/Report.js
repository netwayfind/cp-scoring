import './App.css';
import { apiGet } from './common/utils';
import HostReport from './report/HostReport';

import { Component } from 'react';
import { Link, Route, Switch, withRouter } from 'react-router-dom';

class Report extends Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            hostnames: [],
            scenarioID: 1,
            teamKey: 55555555
        }
    }

    componentDidMount() {
        this.getData(this.state.scenarioID, this.state.teamKey);
    }

    getData(scenarioID, teamKey) {
        apiGet('/api/scenarios/' + scenarioID + '/report/hostnames?team_key=' + teamKey)
        .then(async function(s) {
            this.setState({
                error: s.error,
                hostnames: s.data,
            })
        }.bind(this));
    }

    render() {
        let hostnames = [];
        this.state.hostnames.forEach(hostname => {
            hostnames.push(
                <li key={hostname}><Link to={`${this.props.match.url}/${hostname}`}>{hostname}</Link></li>
            );
        });
        return (
            <div className="Report">
                <ul>{hostnames}</ul>
                <hr />
                <Switch>
                    <Route path={`${this.props.match.url}/:hostname`}>
                        <HostReport scenarioID={this.state.scenarioID} teamKey={this.state.teamKey} />
                    </Route>
                </Switch>
            </div>
        );
    }
}

export default withRouter(Report);
