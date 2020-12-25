import React from 'react';
import { withRouter } from 'react-router-dom/cjs/react-router-dom.min';

class Team extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            team: {}
        }
    }

    componentDidMount() {
        let id = this.props.match.params.id;
        if (id) {
            this.getTeam(id);
        }
    }

    componentDidUpdate(prevProps) {
        let id = this.props.match.params.id;
        if (id !== prevProps.match.params.id) {
            this.getTeam(id);
        }
    }

    getTeam(id) {
        fetch('/api/teams/' + id)
        .then(async function(response) {
            let error = null;
            let team = {};
            if (response.status === 200) {
                team = await response.json();
            } else {
                error = await response.text();
            }
            this.setState({
                error: error,
                team: team
            });
        }.bind(this));
    }

    render() {
        return (
            <div>
                <h1>{this.state.error}</h1>
                <input value={this.state.team.Name} />
                <input value={this.state.team.POC} />
                <input value={this.state.team.Email} />
                <input type="checkbox" value={this.state.team.Enabled} />
                <input value={this.state.team.Key} />
            </div>
        );
    }
}

export default withRouter(Team);
