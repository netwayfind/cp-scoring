import '../App.css';
import { apiGet } from '../utils';

import { Component } from 'react';
import { withRouter } from 'react-router-dom/cjs/react-router-dom.min';

class Scenario extends Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            scenario: {}
        }

        this.getData = this.getData.bind(this);
    }

    componentDidMount() {
        let id = this.props.match.params.id;
        if (id) {
            this.getData(id);
        }
    }

    componentDidUpdate(prevProps) {
        let id = this.props.match.params.id;
        if (id !== prevProps.match.params.id) {
            this.getData(id);
        }
    }

    getData(id) {
        if (!id) {
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

    render() {
        return (
            <div>
                <h1>{this.state.error}</h1>
                <input value={this.state.scenario.Name} />
                <input value={this.state.scenario.Description} />
                <input type="checkbox" value={this.state.scenario.Enabled} />
            </div>
        );
    }
}

export default withRouter(Scenario);
