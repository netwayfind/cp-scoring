import '../App.css';
import { apiGet } from '../utils';
import Team from './Team';

import { Component } from 'react';
import { Link, Route, Switch } from 'react-router-dom';
import { withRouter } from 'react-router-dom/cjs/react-router-dom.min';

class Teams extends Component {
    constructor(props) {
        super(props);
        this.state = this.defaultState();

        this.getData = this.getData.bind(this);
    }

    componentDidMount() {
        this.getData();
    }

    defaultState() {
        return {
            error: null,
            teams: []
        }
    }

    getData() {
        apiGet("/api/teams/")
        .then(function(s) {
            this.setState({
                error: s.error,
                teams: s.data
            });
        }.bind(this));
    }

    render() {
        let teams = [];
        this.state.teams.forEach((team, i) => {
            teams.push(
                <li key={i}><Link to={`${this.props.match.path}/${team.ID}`}>{team.Name}</Link></li>
            );
        });
        return (
            <div className="Teams">
                <Link to={this.props.match.path}>Add Team</Link>
                <ul>{teams}</ul>
                <Switch>
                    <Route path={`${this.props.match.url}/:id`}>
                        <Team callback={this.getData} />
                    </Route>
                    <Route>
                        <Team callback={this.getData} />
                    </Route>
                </Switch>
            </div>
        );
    }
}

export default withRouter(Teams);
