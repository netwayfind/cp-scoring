import '../App.css';
import LinkList from '../common/LinkList';
import { apiGet } from '../common/utils';
import User from './User';

import { Component } from 'react';
import { Link, Route, Switch } from 'react-router-dom';
import { withRouter } from 'react-router-dom/cjs/react-router-dom.min';

class Users extends Component {
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
            users: []
        }
    }

    getData() {
        apiGet("/api/users/")
        .then(function(s) {
            this.setState({
                error: s.error,
                users: s.data
            });
        }.bind(this));
    }

    render() {
        return (
            <div className="Users">
                <Link to={this.props.match.path}>Add User</Link>
                <p />
                <LinkList items={this.state.users} path={this.props.match.path} label="Username" />
                <Switch>
                    <Route path={`${this.props.match.url}/:id`}>
                        <User parentCallback={this.getData} parentPath={this.props.match.path} />
                    </Route>
                    <Route>
                        <User parentCallback={this.getData} parentPath={this.props.match.path}/>
                    </Route>
                </Switch>
            </div>
        );
    }
}

export default withRouter(Users);
