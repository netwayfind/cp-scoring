'use strict';

class Error extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    if (this.props.message === null) {
      return null;
    }

    return (
      <div class="error">
        {this.props.message}
      </div>
    )
  }
}

class Login extends React.Component {
  constructor() {
    super();
    this.state = {
      username: "",
      password: "",
      error: null
    }

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    let value = event.target.value;
    this.setState({
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
    .then(async function(response) {
      this.props.callback(response.status);
      if (response.status === 200) {
        return {
          error: null
        }
      }
      let text = await response.text();
      return {
        error: text
      }
    }.bind(this))
    .then(function(s) {
      this.setState(s);
    }.bind(this));
  }

  render() {
    return (
      <div className="login">
        <form onChange={this.handleChange} onSubmit={this.handleSubmit}>
          <label htmlFor="username">Username</label>
          <input name="username" required="required"></input>
          <br />
          <label htmlFor="password">Password</label>
          <input name="password" type="password" required="required"></input>
          <br />
          <button type="submit">Submit</button>
        </form>
        <Error message={this.state.error} />
      </div>
    );
  }
}