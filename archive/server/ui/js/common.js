'use strict';

class Error extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    if (this.props.message === null) {
      return null;
    }

    return /*#__PURE__*/React.createElement("div", {
      class: "error"
    }, this.props.message);
  }

}

class Login extends React.Component {
  constructor() {
    super();
    this.state = {
      username: "",
      password: "",
      error: null
    };
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
    }).then(async function (response) {
      this.props.callback(response.status);

      if (response.status === 200) {
        return {
          error: null
        };
      }

      let text = await response.text();
      return {
        error: text
      };
    }.bind(this)).then(function (s) {
      this.setState(s);
    }.bind(this));
  }

  render() {
    return /*#__PURE__*/React.createElement("div", {
      className: "login"
    }, /*#__PURE__*/React.createElement("form", {
      onChange: this.handleChange,
      onSubmit: this.handleSubmit
    }, /*#__PURE__*/React.createElement("label", {
      htmlFor: "username"
    }, "Username"), /*#__PURE__*/React.createElement("input", {
      name: "username",
      required: "required"
    }), /*#__PURE__*/React.createElement("br", null), /*#__PURE__*/React.createElement("label", {
      htmlFor: "password"
    }, "Password"), /*#__PURE__*/React.createElement("input", {
      name: "password",
      type: "password",
      required: "required"
    }), /*#__PURE__*/React.createElement("br", null), /*#__PURE__*/React.createElement("button", {
      type: "submit"
    }, "Submit")), /*#__PURE__*/React.createElement(Error, {
      message: this.state.error
    }));
  }

}