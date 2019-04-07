'use strict';

class Login extends React.Component {
  constructor() {
    super();
    this.state = {
      username: "",
      password: "",
      messages: ""
    }

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    let value = event.target.value;
    this.setState({
      ...this.state.credentials,
      [event.target.name]: value
    });
  }

  handleSubmit(event) {
    event.preventDefault();

    if (this.state.username.length == 0 || this.state.password.length == 0) {
      return;
    }

    var url = "/login";

    fetch(url, {
      credentials: 'same-origin',
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
      },
      body: "username=" + this.state.username + "&password=" + this.state.password
    })
    .then(function(response) {
      this.props.callback(response.status);
      if (response.status >= 400) {
        this.setState({
          messages: "Rejected login. Try again."
        });
      }
    }.bind(this));
  }

  render() {
    return (
      <div className="login">
        {this.state.messages}
        <form onChange={this.handleChange} onSubmit={this.handleSubmit}>
          <label htmlFor="username">Username</label>
          <input name="username" required="required"></input>
          <br />
          <label htmlFor="password">Password</label>
          <input name="password" type="password" required="required"></input>
          <br />
          <button type="submit">Submit</button>
        </form>
      </div>
    );
  }
}