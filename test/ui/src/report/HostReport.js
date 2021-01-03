import '../App.css';
import { apiGet } from '../common/utils';

import { Component } from 'react';
import { withRouter } from 'react-router-dom';

class HostReport extends Component {
    constructor(props) {
        super(props);
        this.state = {
            report: {
                AnswerResults: []
            },
            timeline: []
        }
    }

    componentDidMount() {
        let hostname = this.props.match.params.hostname;
        let scenarioID = this.props.scenarioID;
        let teamKey = this.props.teamKey;
        this.getData(scenarioID, teamKey, hostname);
    }

    componentDidUpdate(prevProps) {
        let hostname = this.props.match.params.hostname;
        let prevHostname = prevProps.match.params.hostname;
        let scenarioID = this.props.scenarioID;
        let teamKey = this.props.teamKey;
        if (hostname !== prevHostname) {
            this.getData(scenarioID, teamKey, hostname);
        }
    }

    getData(scenarioID, teamKey, hostname) {
        Promise.all([
            apiGet('/api/scenarios/' + scenarioID + '/report?team_key=' + teamKey + '&hostname=' + hostname),
            apiGet('/api/scenarios/' + scenarioID + '/report/timeline?team_key=' + teamKey + '&hostname=' + hostname)
        ])
        .then(async function(responses) {
            let s1 = responses[0];
            let s2 = responses[1];
            this.setState({
                error: s1.error || s2.error,
                report: s1.data,
                timeline: s2.data
            })
        }.bind(this));
    }

    render() {
        let timestampStr = new Date(this.state.report.Timestamp * 1000).toLocaleString();
        let score = 0;
        let results = [];
        this.state.report.AnswerResults.forEach((result, i) => {
            results.push(
                <li key={i}><strong>{result.Points}</strong> - {result.Description}</li>
            );
            score += result.Points;
        });
        return (
            <div className="HostReport">
                [Timeline]
                <p />
                Last Updated: {timestampStr}
                <p />
                Score: {score}
                <p />
                Results:
                <p />
                <ul>{results}</ul>
            </div>
        )
    }
}

export default withRouter(HostReport);
