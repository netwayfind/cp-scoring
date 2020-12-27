import '../App.css';
import { apiGet } from '../utils';
import Scenario from './Scenario';

import { Component } from 'react';
import { Link, Route, Switch } from 'react-router-dom';
import { withRouter } from 'react-router-dom/cjs/react-router-dom.min';

class Scenarios extends Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            scenarios: []
        }

        this.getData = this.getData.bind(this);
    }

    componentDidMount() {
        this.getData();
    }

    getData() {
        apiGet("/api/scenarios/")
        .then(function(s) {
            this.setState({
                error: s.error,
                scenarios: s.data
            });
        }.bind(this));
    }

    render() {
        let scenarios = [];
        this.state.scenarios.forEach((scenario, i) => {
            scenarios.push(
                <li key={i}><Link to={`${this.props.match.path}/${scenario.ID}`}>{scenario.Name}</Link></li>
            );
        });
        return (
            <div className="Scenarios">
                <ul>{scenarios}</ul>
                <Switch>
                    <Route path={`${this.props.match.url}/:id`}>
                        <Scenario />
                    </Route>
                    <Route>
                        <Scenario />
                    </Route>
                </Switch>
            </div>
        );
    }
}

export default withRouter(Scenarios);
