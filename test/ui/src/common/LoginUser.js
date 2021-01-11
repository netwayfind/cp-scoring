import '../App.css';
import { apiLogin } from '../common/utils';

import { Component } from 'react';
import { withRouter } from 'react-router-dom';

class LoginUser extends Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            username: "",
            password: ""
        }

        this.handleUpdate = this.handleUpdate.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleSubmit(event) {
        event.preventDefault();

        if (this.state.username.length === 0 || this.state.password.length === 0) {
            return;
        }

        apiLogin(this.state.username, this.state.password)
        .then(async function(s) {
            this.setState({
                error: s.error
            });
            if (!s.error) {
                this.props.callback();
                this.props.history.push(this.props.location);
            }
        }.bind(this));
    }

    handleUpdate(event) {
        let value = event.target.value;
        this.setState({
            [event.target.name]: value
        });
    }

    render() {
        return (
            <div className="login">
                <form onChange={this.handleUpdate} onSubmit={this.handleSubmit}>
                    <label htmlFor="username">Username</label>
                    <input name="username" required="required"></input>
                    <br />
                    <label htmlFor="password">Password</label>
                    <input name="password" type="password" required="required"></input>
                    <br />
                    <button type="submit">Submit</button>
                </form>
                <h1>{this.state.error}</h1>
            </div>
        );
    }
}

export default withRouter(LoginUser);
