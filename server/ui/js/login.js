'use strict';

function _objectSpread(target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i] != null ? arguments[i] : {}; var ownKeys = Object.keys(source); if (typeof Object.getOwnPropertySymbols === 'function') { ownKeys = ownKeys.concat(Object.getOwnPropertySymbols(source).filter(function (sym) { return Object.getOwnPropertyDescriptor(source, sym).enumerable; })); } ownKeys.forEach(function (key) { _defineProperty(target, key, source[key]); }); } return target; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

class Login extends React.Component {
  constructor() {
    super();
    this.state = {
      username: "",
      password: "",
      messages: ""
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    let value = event.target.value;
    this.setState(_objectSpread({}, this.state.credentials, {
      [event.target.name]: value
    }));
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
    }).then(function (response) {
      this.props.callback(response.status);

      if (response.status >= 400) {
        this.setState({
          messages: "Rejected login. Try again."
        });
      }
    }.bind(this));
  }

  render() {
    return React.createElement("div", {
      className: "login"
    }, this.state.messages, React.createElement("form", {
      onChange: this.handleChange,
      onSubmit: this.handleSubmit
    }, React.createElement("label", {
      htmlFor: "username"
    }, "Username"), React.createElement("input", {
      name: "username",
      required: "required"
    }), React.createElement("br", null), React.createElement("label", {
      htmlFor: "password"
    }, "Password"), React.createElement("input", {
      name: "password",
      type: "password",
      required: "required"
    }), React.createElement("br", null), React.createElement("button", {
      type: "submit"
    }, "Submit")));
  }

}