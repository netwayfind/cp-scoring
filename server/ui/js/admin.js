'use strict';

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _toConsumableArray(arr) { if (Array.isArray(arr)) { for (var i = 0, arr2 = Array(arr.length); i < arr.length; i++) { arr2[i] = arr[i]; } return arr2; } else { return Array.from(arr); } }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _possibleConstructorReturn(self, call) { if (!self) { throw new ReferenceError("this hasn't been initialised - super() hasn't been called"); } return call && (typeof call === "object" || typeof call === "function") ? call : self; }

function _inherits(subClass, superClass) { if (typeof superClass !== "function" && superClass !== null) { throw new TypeError("Super expression must either be null or a function, not " + typeof superClass); } subClass.prototype = Object.create(superClass && superClass.prototype, { constructor: { value: subClass, enumerable: false, writable: true, configurable: true } }); if (superClass) Object.setPrototypeOf ? Object.setPrototypeOf(subClass, superClass) : subClass.__proto__ = superClass; }

var Plot = createPlotlyComponent(Plotly);

var App = function (_React$Component) {
  _inherits(App, _React$Component);

  function App() {
    _classCallCheck(this, App);

    var _this = _possibleConstructorReturn(this, (App.__proto__ || Object.getPrototypeOf(App)).call(this));

    _this.state = {
      authenticated: false
    };

    _this.authCallback = _this.authCallback.bind(_this);
    _this.logout = _this.logout.bind(_this);
    return _this;
  }

  _createClass(App, [{
    key: 'authCallback',
    value: function authCallback(statusCode) {
      if (statusCode == 200) {
        this.setState({
          authenticated: true
        });
      } else {
        this.setState({
          authenticated: false
        });
      }
    }
  }, {
    key: 'logout',
    value: function logout() {
      var url = "/logout";
      fetch(url, {
        credentials: 'same-origin',
        method: "DELETE"
      }).then(function (_) {
        this.setState({
          authenticated: false
        });
      }.bind(this));
    }
  }, {
    key: 'componentDidMount',
    value: function componentDidMount() {
      // check if logged in by visiting the following URL
      var url = "/templates";
      fetch(url, {
        credentials: 'same-origin'
      }).then(function (response) {
        this.authCallback(response.status);
      }.bind(this));
    }
  }, {
    key: 'render',
    value: function render() {
      if (!this.state.authenticated) {
        return React.createElement(
          'div',
          { className: 'App' },
          React.createElement(Login, { callback: this.authCallback })
        );
      }
      return React.createElement(
        'div',
        { className: 'App' },
        React.createElement(
          'button',
          { onClick: this.logout },
          'Logout'
        ),
        React.createElement('p', null),
        React.createElement(Teams, null),
        React.createElement(Hosts, null),
        React.createElement(Templates, null),
        React.createElement(Scenarios, null)
      );
    }
  }]);

  return App;
}(React.Component);

var Login = function (_React$Component2) {
  _inherits(Login, _React$Component2);

  function Login() {
    _classCallCheck(this, Login);

    var _this2 = _possibleConstructorReturn(this, (Login.__proto__ || Object.getPrototypeOf(Login)).call(this));

    _this2.state = {
      username: "",
      password: "",
      messages: ""
    };

    _this2.handleChange = _this2.handleChange.bind(_this2);
    _this2.handleSubmit = _this2.handleSubmit.bind(_this2);
    return _this2;
  }

  _createClass(Login, [{
    key: 'handleChange',
    value: function handleChange(event) {
      var value = event.target.value;
      this.setState(Object.assign({}, this.state.credentials, _defineProperty({}, event.target.name, value)));
    }
  }, {
    key: 'handleSubmit',
    value: function handleSubmit(event) {
      event.preventDefault();

      if (this.state.username.length == 0 || this.state.password.length == 0) {
        return;
      }

      var url = "/login";

      fetch(url, {
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
  }, {
    key: 'render',
    value: function render() {
      return React.createElement(
        'div',
        { className: 'login' },
        this.state.messages,
        React.createElement(
          'form',
          { onChange: this.handleChange, onSubmit: this.handleSubmit },
          React.createElement(
            'label',
            { htmlFor: 'username' },
            'Username'
          ),
          React.createElement('input', { name: 'username', required: 'required' }),
          React.createElement('br', null),
          React.createElement(
            'label',
            { htmlFor: 'password' },
            'Password'
          ),
          React.createElement('input', { name: 'password', type: 'password', required: 'required' }),
          React.createElement('br', null),
          React.createElement(
            'button',
            { type: 'submit' },
            'Submit'
          )
        )
      );
    }
  }]);

  return Login;
}(React.Component);

var backgroundStyle = {
  position: 'fixed',
  top: 0,
  bottom: 0,
  left: 0,
  right: 0,
  backgroundColor: 'rgba(0,0,0,0.5)',
  padding: 50
};

var modalStyle = {
  backgroundColor: 'white',
  padding: 30,
  maxHeight: '100%',
  overflowY: 'auto'
};

var BasicModal = function (_React$Component3) {
  _inherits(BasicModal, _React$Component3);

  function BasicModal(props) {
    _classCallCheck(this, BasicModal);

    var _this3 = _possibleConstructorReturn(this, (BasicModal.__proto__ || Object.getPrototypeOf(BasicModal)).call(this, props));

    _this3.state = _this3.defaultState();

    _this3.handleChange = _this3.handleChange.bind(_this3);
    _this3.handleSubmit = _this3.handleSubmit.bind(_this3);
    _this3.handleClose = _this3.handleClose.bind(_this3);
    _this3.setValue = _this3.setValue.bind(_this3);
    return _this3;
  }

  _createClass(BasicModal, [{
    key: 'defaultState',
    value: function defaultState() {
      return {
        subject: {}
      };
    }
  }, {
    key: 'setValue',
    value: function setValue(key, value) {
      this.setState({
        subject: Object.assign({}, this.props.subject, this.state.subject, _defineProperty({}, key, value))
      });
    }
  }, {
    key: 'handleChange',
    value: function handleChange(event) {
      var value = event.target.value;
      if (event.target.type == 'checkbox') {
        value = event.target.checked;
      }
      this.setState({
        subject: Object.assign({}, this.props.subject, this.state.subject, _defineProperty({}, event.target.name, value))
      });
    }
  }, {
    key: 'handleSubmit',
    value: function handleSubmit(event) {
      event.preventDefault();

      if (Object.keys(this.state.subject) == 0) {
        return;
      }

      var url = "/" + this.props.subjectClass;
      if (this.props.subjectID != null) {
        url += "/" + this.props.subjectID;
      }

      fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(this.state.subject)
      }).then(function (response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        this.props.submit();
        this.setState(this.defaultState());
      }.bind(this));
    }
  }, {
    key: 'handleClose',
    value: function handleClose() {
      this.props.onClose();
      this.setState(this.defaultState());
    }
  }, {
    key: 'render',
    value: function render() {
      if (!this.props.show) {
        return null;
      }

      return React.createElement(
        'div',
        { className: 'background', style: backgroundStyle },
        React.createElement(
          'div',
          { className: 'modal', style: modalStyle },
          React.createElement(
            'label',
            { htmlFor: 'ID' },
            'ID'
          ),
          React.createElement('input', { name: 'ID', defaultValue: this.props.subjectID, disabled: true }),
          React.createElement('br', null),
          React.createElement(
            'form',
            { onChange: this.handleChange, onSubmit: this.handleSubmit },
            this.props.children,
            React.createElement('br', null),
            React.createElement(
              'button',
              { type: 'submit' },
              'Submit'
            ),
            React.createElement(
              'button',
              { type: 'button', onClick: this.handleClose },
              'Cancel'
            )
          )
        )
      );
    }
  }]);

  return BasicModal;
}(React.Component);

var Teams = function (_React$Component4) {
  _inherits(Teams, _React$Component4);

  function Teams() {
    _classCallCheck(this, Teams);

    var _this4 = _possibleConstructorReturn(this, (Teams.__proto__ || Object.getPrototypeOf(Teams)).call(this));

    _this4.toggleModal = function () {
      _this4.setState({
        showModal: !_this4.state.showModal
      });
    };

    _this4.state = {
      teams: [],
      showModal: false,
      selectedTeamID: null,
      selectedTeam: {}
    };

    _this4.modal = React.createRef();
    _this4.handleSubmit = _this4.handleSubmit.bind(_this4);
    _this4.regenKey = _this4.regenKey.bind(_this4);
    return _this4;
  }

  _createClass(Teams, [{
    key: 'componentDidMount',
    value: function componentDidMount() {
      this.populateTeams();
    }
  }, {
    key: 'populateTeams',
    value: function populateTeams() {
      var url = '/teams';

      fetch(url).then(function (response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        return response.json();
      }).then(function (data) {
        this.setState({ teams: data });
      }.bind(this));
    }
  }, {
    key: 'newKey',
    value: function newKey() {
      return Math.random().toString(16).substring(7);
    }
  }, {
    key: 'createTeam',
    value: function createTeam() {
      this.setState({
        selectedTeamID: null,
        selectedTeam: {
          Enabled: true,
          Key: this.newKey()
        }
      });
      this.toggleModal();
    }
  }, {
    key: 'editTeam',
    value: function editTeam(id) {
      var url = "/teams/" + id;

      fetch(url).then(function (response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        return response.json();
      }).then(function (data) {
        this.setState({
          selectedTeamID: id,
          selectedTeam: data
        });
        this.toggleModal();
      }.bind(this));
    }
  }, {
    key: 'deleteTeam',
    value: function deleteTeam(id) {
      var url = "/teams/" + id;

      fetch(url, {
        method: 'DELETE'
      }).then(function (response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        this.populateTeams();
      }.bind(this));
    }
  }, {
    key: 'handleSubmit',
    value: function handleSubmit() {
      this.populateTeams();
      this.toggleModal();
    }
  }, {
    key: 'regenKey',
    value: function regenKey() {
      var key = this.newKey();
      this.setState({
        selectedTeam: Object.assign({}, this.state.selectedTeam, {
          Key: key
        })
      });
      this.modal.current.setValue("Key", key);
    }
  }, {
    key: 'render',
    value: function render() {
      var rows = [];
      for (var i = 0; i < this.state.teams.length; i++) {
        var team = this.state.teams[i];
        rows.push(React.createElement(
          'li',
          { key: team.ID },
          team.ID,
          ' - ',
          team.Name,
          React.createElement(
            'button',
            { type: 'button', onClick: this.editTeam.bind(this, team.ID) },
            'Edit'
          ),
          React.createElement(
            'button',
            { type: 'button', onClick: this.deleteTeam.bind(this, team.ID) },
            '-'
          )
        ));
      }

      return React.createElement(
        'div',
        { className: 'Teams' },
        React.createElement(
          'strong',
          null,
          'Teams'
        ),
        React.createElement('p', null),
        React.createElement(
          'button',
          { onClick: this.createTeam.bind(this) },
          'Add Team'
        ),
        React.createElement(
          BasicModal,
          { ref: this.modal, subjectClass: 'teams', subjectID: this.state.selectedTeamID, subject: this.state.selectedTeam, show: this.state.showModal, onClose: this.toggleModal, submit: this.handleSubmit },
          React.createElement(Item, { name: 'Name', defaultValue: this.state.selectedTeam.Name }),
          React.createElement(Item, { name: 'POC', defaultValue: this.state.selectedTeam.POC }),
          React.createElement(Item, { name: 'Email', type: 'email', defaultValue: this.state.selectedTeam.Email }),
          React.createElement(
            'label',
            { htmlFor: 'Enabled' },
            'Enabled'
          ),
          React.createElement('input', { name: 'Enabled', type: 'checkbox', defaultChecked: !!this.state.selectedTeam.Enabled }),
          React.createElement('br', null),
          React.createElement(
            'label',
            { htmlFor: 'Key' },
            'Key'
          ),
          React.createElement(
            'ul',
            null,
            React.createElement(
              'li',
              null,
              this.state.selectedTeam.Key
            ),
            React.createElement(
              'button',
              { type: 'button', onClick: this.regenKey.bind(this) },
              'Regenerate'
            )
          )
        ),
        React.createElement(
          'ul',
          null,
          rows
        )
      );
    }
  }]);

  return Teams;
}(React.Component);

var Scenarios = function (_React$Component5) {
  _inherits(Scenarios, _React$Component5);

  function Scenarios() {
    _classCallCheck(this, Scenarios);

    var _this5 = _possibleConstructorReturn(this, (Scenarios.__proto__ || Object.getPrototypeOf(Scenarios)).call(this));

    _this5.toggleModal = function () {
      _this5.setState({
        showModal: !_this5.state.showModal
      });
    };

    _this5.state = {
      scenarios: [],
      showModal: false,
      selectedScenario: {}
    };
    _this5.modal = React.createRef();

    _this5.handleSubmit = _this5.handleSubmit.bind(_this5);
    _this5.handleCallback = _this5.handleCallback.bind(_this5);
    _this5.mapItems = _this5.mapItems.bind(_this5);
    _this5.listItems = _this5.listItems.bind(_this5);
    return _this5;
  }

  _createClass(Scenarios, [{
    key: 'componentDidMount',
    value: function componentDidMount() {
      this.populateScenarios();
    }
  }, {
    key: 'populateScenarios',
    value: function populateScenarios() {
      var url = '/scenarios';

      fetch(url).then(function (response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        return response.json();
      }).then(function (data) {
        this.setState({ scenarios: data });
      }.bind(this));
    }
  }, {
    key: 'createScenario',
    value: function createScenario() {
      this.setState({
        selectedScenarioID: null,
        selectedScenario: { Enabled: true }
      });
      this.toggleModal();
    }
  }, {
    key: 'editScenario',
    value: function editScenario(id) {
      var url = "/scenarios/" + id;

      fetch(url).then(function (response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        return response.json();
      }).then(function (data) {
        this.setState({
          selectedScenarioID: id,
          selectedScenario: data
        });
        this.toggleModal();
      }.bind(this));
    }
  }, {
    key: 'deleteScenario',
    value: function deleteScenario(id) {
      var url = "/scenarios/" + id;

      fetch(url, {
        method: 'DELETE'
      }).then(function (response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        this.populateScenarios();
      }.bind(this));
    }
  }, {
    key: 'handleSubmit',
    value: function handleSubmit() {
      this.populateScenarios();
      this.toggleModal();
    }
  }, {
    key: 'handleCallback',
    value: function handleCallback(key, value) {
      this.modal.current.setValue(key, value);
    }
  }, {
    key: 'mapItems',
    value: function mapItems(callback) {
      var url = "/hosts";

      fetch(url).then(function (response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        return response.json();
      }.bind(this)).then(function (data) {
        var items = data.map(function (host) {
          return {
            ID: host.ID,
            Display: host.Hostname
          };
        });
        callback(items);
      });
    }
  }, {
    key: 'listItems',
    value: function listItems(callback) {
      var url = "/templates";

      fetch(url).then(function (response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        return response.json();
      }.bind(this)).then(function (data) {
        var items = data.map(function (template) {
          return {
            ID: template.ID,
            Display: template.Name
          };
        });
        callback(items);
      });
    }
  }, {
    key: 'render',
    value: function render() {
      var rows = [];
      for (var i = 0; i < this.state.scenarios.length; i++) {
        var scenario = this.state.scenarios[i];
        rows.push(React.createElement(
          'li',
          { key: scenario.ID },
          scenario.ID,
          ' - ',
          scenario.Name,
          React.createElement(
            'button',
            { type: 'button', onClick: this.editScenario.bind(this, scenario.ID) },
            'Edit'
          ),
          React.createElement(
            'button',
            { type: 'button', onClick: this.deleteScenario.bind(this, scenario.ID) },
            '-'
          )
        ));
      }

      return React.createElement(
        'div',
        { className: 'Scenarios' },
        React.createElement(
          'strong',
          null,
          'Scenarios'
        ),
        React.createElement('p', null),
        React.createElement(
          'button',
          { onClick: this.createScenario.bind(this) },
          'Add Scenario'
        ),
        React.createElement(
          BasicModal,
          { ref: this.modal, subjectClass: 'scenarios', subjectID: this.state.selectedScenarioID, subject: this.state.selectedScenario, show: this.state.showModal, onClose: this.toggleModal, submit: this.handleSubmit },
          React.createElement(Item, { name: 'Name', defaultValue: this.state.selectedScenario.Name }),
          React.createElement(Item, { name: 'Description', defaultValue: this.state.selectedScenario.Description }),
          React.createElement(Item, { name: 'Enabled', type: 'checkbox', defaultChecked: !!this.state.selectedScenario.Enabled }),
          React.createElement(ItemMap, { name: 'HostTemplates', label: 'Hosts', listLabel: 'Templates', defaultValue: this.state.selectedScenario.HostTemplates, callback: this.handleCallback, mapItems: this.mapItems, listItems: this.listItems })
        ),
        React.createElement(
          'ul',
          null,
          rows
        )
      );
    }
  }]);

  return Scenarios;
}(React.Component);

var Hosts = function (_React$Component6) {
  _inherits(Hosts, _React$Component6);

  function Hosts() {
    _classCallCheck(this, Hosts);

    var _this6 = _possibleConstructorReturn(this, (Hosts.__proto__ || Object.getPrototypeOf(Hosts)).call(this));

    _this6.toggleModal = function () {
      _this6.setState({
        showModal: !_this6.state.showModal
      });
    };

    _this6.state = {
      hosts: [],
      showModal: false,
      selectedHost: {}
    };

    _this6.handleSubmit = _this6.handleSubmit.bind(_this6);
    return _this6;
  }

  _createClass(Hosts, [{
    key: 'componentDidMount',
    value: function componentDidMount() {
      this.populateHosts();
    }
  }, {
    key: 'populateHosts',
    value: function populateHosts() {
      var url = '/hosts';

      fetch(url).then(function (response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        return response.json();
      }).then(function (data) {
        this.setState({ hosts: data });
      }.bind(this));
    }
  }, {
    key: 'createHost',
    value: function createHost() {
      this.setState({
        selectedHostID: null,
        selectedHost: {}
      });
      this.toggleModal();
    }
  }, {
    key: 'editHost',
    value: function editHost(id, host) {
      this.setState({
        selectedHostID: id,
        selectedHost: host
      });
      this.toggleModal();
    }
  }, {
    key: 'deleteHost',
    value: function deleteHost(id) {
      var url = "/hosts/" + id;

      fetch(url, {
        method: 'DELETE'
      }).then(function (response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        this.populateHosts();
      }.bind(this));
    }
  }, {
    key: 'handleSubmit',
    value: function handleSubmit() {
      this.populateHosts();
      this.toggleModal();
    }
  }, {
    key: 'render',
    value: function render() {
      var rows = [];
      for (var i = 0; i < this.state.hosts.length; i++) {
        var host = this.state.hosts[i];
        rows.push(React.createElement(
          'li',
          { key: host.ID },
          host.ID,
          ' - ',
          host.Hostname,
          ' - ',
          host.OS,
          React.createElement(
            'button',
            { type: 'button', onClick: this.editHost.bind(this, host.ID, host) },
            'Edit'
          ),
          React.createElement(
            'button',
            { type: 'button', onClick: this.deleteHost.bind(this, host.ID) },
            '-'
          )
        ));
      }

      return React.createElement(
        'div',
        { className: 'Hosts' },
        React.createElement(
          'strong',
          null,
          'Hosts'
        ),
        React.createElement('p', null),
        React.createElement(
          'button',
          { onClick: this.createHost.bind(this) },
          'Add Host'
        ),
        React.createElement(
          BasicModal,
          { subjectClass: 'hosts', subjectID: this.state.selectedHostID, subject: this.state.selectedHost, show: this.state.showModal, onClose: this.toggleModal, submit: this.handleSubmit },
          React.createElement(Item, { name: 'Hostname', type: 'text', defaultValue: this.state.selectedHost.Hostname }),
          React.createElement(Item, { name: 'OS', type: 'text', defaultValue: this.state.selectedHost.OS })
        ),
        React.createElement(
          'ul',
          null,
          rows
        )
      );
    }
  }]);

  return Hosts;
}(React.Component);

var Templates = function (_React$Component7) {
  _inherits(Templates, _React$Component7);

  function Templates() {
    _classCallCheck(this, Templates);

    var _this7 = _possibleConstructorReturn(this, (Templates.__proto__ || Object.getPrototypeOf(Templates)).call(this));

    _this7.toggleModal = function () {
      _this7.setState({
        showModal: !_this7.state.showModal
      });
    };

    _this7.state = {
      templates: [],
      showModal: false,
      selectedTemplate: {
        Template: {}
      }
    };
    _this7.modal = React.createRef();

    _this7.handleSubmit = _this7.handleSubmit.bind(_this7);
    _this7.handleCallback = _this7.handleCallback.bind(_this7);
    return _this7;
  }

  _createClass(Templates, [{
    key: 'componentDidMount',
    value: function componentDidMount() {
      this.populateTemplates();
    }
  }, {
    key: 'populateTemplates',
    value: function populateTemplates() {
      var url = "/templates";

      fetch(url).then(function (response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        return response.json();
      }).then(function (data) {
        this.setState({ templates: data });
      }.bind(this));
    }
  }, {
    key: 'createTemplate',
    value: function createTemplate() {
      this.setState({
        selectedTemplateID: null,
        selectedTemplate: {
          Template: {}
        }
      });
      this.toggleModal();
    }
  }, {
    key: 'editTemplate',
    value: function editTemplate(id) {
      var url = "/templates/" + id;

      fetch(url).then(function (response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        return response.json();
      }).then(function (data) {
        this.setState({
          selectedTemplateID: id,
          selectedTemplate: data
        });
        this.toggleModal();
      }.bind(this));
    }
  }, {
    key: 'deleteTemplate',
    value: function deleteTemplate(id) {
      var url = "/templates/" + id;

      fetch(url, {
        method: 'DELETE'
      }).then(function (response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        this.populateTemplates();
      }.bind(this));
    }
  }, {
    key: 'handleSubmit',
    value: function handleSubmit() {
      this.populateTemplates();
      this.toggleModal();
    }
  }, {
    key: 'handleCallback',
    value: function handleCallback(key, value) {
      var template = Object.assign({}, this.state.selectedTemplate.Template, _defineProperty({}, key, value));
      this.setState({
        selectedTemplate: Object.assign({}, this.state.selectedTemplate, {
          Template: template
        })
      });
      this.modal.current.setValue("Template", template);
    }
  }, {
    key: 'render',
    value: function render() {
      var rows = [];
      for (var i = 0; i < this.state.templates.length; i++) {
        var template = this.state.templates[i];
        rows.push(React.createElement(
          'li',
          { key: template.ID },
          template.ID,
          ' - ',
          template.Name,
          React.createElement(
            'button',
            { type: 'button', onClick: this.editTemplate.bind(this, template.ID) },
            'Edit'
          ),
          React.createElement(
            'button',
            { type: 'button', onClick: this.deleteTemplate.bind(this, template.ID) },
            '-'
          )
        ));
      }

      return React.createElement(
        'div',
        { className: 'Templates' },
        React.createElement(
          'strong',
          null,
          'Templates'
        ),
        React.createElement('p', null),
        React.createElement(
          'button',
          { onClick: this.createTemplate.bind(this) },
          'Create Template'
        ),
        React.createElement(
          BasicModal,
          { ref: this.modal, subjectClass: 'templates', subjectID: this.state.selectedTemplateID, subject: this.state.selectedTemplate, show: this.state.showModal, onClose: this.toggleModal, submit: this.handleSubmit },
          React.createElement(Item, { name: 'Name', type: 'text', defaultValue: this.state.selectedTemplate.Name }),
          React.createElement(Users, { users: this.state.selectedTemplate.Template.Users, callback: this.handleCallback }),
          React.createElement(Groups, { name: 'GroupMembersAdd', label: 'Group members to add', groups: this.state.selectedTemplate.Template.GroupMembersAdd, callback: this.handleCallback }),
          React.createElement(Groups, { name: 'GroupMembersKeep', label: 'Group members to keep', groups: this.state.selectedTemplate.Template.GroupMembersKeep, callback: this.handleCallback }),
          React.createElement(Groups, { name: 'GroupMembersRemove', label: 'Group members to remove', groups: this.state.selectedTemplate.Template.GroupMembersRemove, callback: this.handleCallback })
        ),
        React.createElement(
          'ul',
          null,
          rows
        )
      );
    }
  }]);

  return Templates;
}(React.Component);

var Users = function (_React$Component8) {
  _inherits(Users, _React$Component8);

  function Users(props) {
    _classCallCheck(this, Users);

    var _this8 = _possibleConstructorReturn(this, (Users.__proto__ || Object.getPrototypeOf(Users)).call(this, props));

    var users = props.users;
    if (users === undefined || users === null) {
      users = [];
    }
    _this8.state = {
      users: users
    };

    _this8.addUser = _this8.addUser.bind(_this8);
    _this8.removeUser = _this8.removeUser.bind(_this8);
    _this8.updateUser = _this8.updateUser.bind(_this8);
    return _this8;
  }

  _createClass(Users, [{
    key: 'addUser',
    value: function addUser() {
      var empty = {
        Name: "",
        AccountPresent: true,
        AccountActive: true,
        PasswordExpires: true,
        // unix timestamp in seconds
        PasswordLastSet: Date.now() / 1000
      };
      var users = [].concat(_toConsumableArray(this.state.users), [empty]);
      this.setState({
        users: users
      });
      this.props.callback("Users", users);
    }
  }, {
    key: 'removeUser',
    value: function removeUser(id) {
      var users = this.state.users.filter(function (_, index) {
        return index != id;
      });
      this.setState({
        users: users
      });
      this.props.callback("Users", users);
    }
  }, {
    key: 'updateUser',
    value: function updateUser(id, field, event) {
      var updated = this.state.users;
      var value = event.target.value;
      if (event.target.type === "checkbox") {
        if (event.target.checked) {
          value = true;
        } else {
          value = false;
        }
      }
      if (event.target.type === "date") {
        value = new Date(event.target.value).getTime() / 1000;
        if (Number.isNaN(value)) {
          return;
        }
      }
      updated[id] = Object.assign({}, updated[id], _defineProperty({}, field, value));
      this.setState({
        users: updated
      });
      this.props.callback("Users", updated);
    }
  }, {
    key: 'render',
    value: function render() {
      var _this9 = this;

      var users = [];

      var _loop = function _loop(i) {
        var user = _this9.state.users[i];
        var d = new Date(user.PasswordLastSet * 1000);
        var passwordLastSet = ("000" + d.getUTCFullYear()).slice(-4);
        passwordLastSet += "-";
        passwordLastSet += ("0" + (d.getUTCMonth() + 1)).slice(-2);
        passwordLastSet += "-";
        passwordLastSet += ("0" + d.getUTCDate()).slice(-2);
        users.push(React.createElement(
          'li',
          { key: "user" + i },
          user.Name,
          React.createElement(
            'button',
            { type: 'button', onClick: _this9.removeUser.bind(_this9, i) },
            '-'
          ),
          React.createElement(
            'ul',
            null,
            React.createElement(
              'li',
              null,
              React.createElement(
                'label',
                null,
                'Name'
              ),
              React.createElement('input', { type: 'text', value: user.Name, onChange: function onChange(event) {
                  return _this9.updateUser(i, "Name", event);
                } })
            ),
            React.createElement(
              'li',
              null,
              React.createElement(
                'label',
                null,
                'Present'
              ),
              React.createElement('input', { type: 'checkbox', checked: user.AccountPresent, onChange: function onChange(event) {
                  return _this9.updateUser(i, "AccountPresent", event);
                } })
            ),
            React.createElement(
              'li',
              null,
              React.createElement(
                'label',
                null,
                'Active'
              ),
              React.createElement('input', { type: 'checkbox', checked: user.AccountActive, onChange: function onChange(event) {
                  return _this9.updateUser(i, "AccountActive", event);
                } })
            ),
            React.createElement(
              'li',
              null,
              React.createElement(
                'label',
                null,
                'Password Expires'
              ),
              React.createElement('input', { type: 'checkbox', checked: user.PasswordExpires, onChange: function onChange(event) {
                  return _this9.updateUser(i, "PasswordExpires", event);
                } })
            ),
            React.createElement(
              'li',
              null,
              React.createElement(
                'label',
                null,
                'Password Last Set'
              ),
              React.createElement('input', { type: 'date', value: passwordLastSet, onChange: function onChange(event) {
                  return _this9.updateUser(i, "PasswordLastSet", event);
                } })
            )
          )
        ));
      };

      for (var i = 0; i < this.state.users.length; i++) {
        _loop(i);
      }

      return React.createElement(
        'div',
        null,
        React.createElement(
          'label',
          { htmlFor: 'Users' },
          'Users'
        ),
        React.createElement('p', null),
        React.createElement(
          'button',
          { type: 'button', onClick: this.addUser.bind(this) },
          'Add User'
        ),
        React.createElement(
          'ul',
          null,
          users
        )
      );
    }
  }]);

  return Users;
}(React.Component);

var Groups = function (_React$Component9) {
  _inherits(Groups, _React$Component9);

  function Groups(props) {
    _classCallCheck(this, Groups);

    var _this10 = _possibleConstructorReturn(this, (Groups.__proto__ || Object.getPrototypeOf(Groups)).call(this, props));

    var groups = props.groups;
    if (groups === undefined || groups === null) {
      groups = {};
    }
    _this10.state = {
      groups: groups
    };

    _this10.newGroupName = React.createRef();

    _this10.addGroup = _this10.addGroup.bind(_this10);
    _this10.removeGroup = _this10.removeGroup.bind(_this10);
    _this10.updateGroup = _this10.updateGroup.bind(_this10);
    return _this10;
  }

  _createClass(Groups, [{
    key: 'addGroup',
    value: function addGroup() {
      if (this.newGroupName.current === null) {
        return;
      }
      var groups = Object.assign({}, this.state.groups, _defineProperty({}, this.newGroupName.current.value, []));
      this.setState({
        groups: groups
      });
      this.props.callback(this.props.name, groups);
    }
  }, {
    key: 'removeGroup',
    value: function removeGroup(name) {
      var groups = this.state.groups;
      delete groups[name];
      this.setState({
        groups: groups
      });
      this.props.callback(this.props.name, groups);
    }
  }, {
    key: 'updateGroup',
    value: function updateGroup(name, members) {
      var groups = Object.assign({}, this.state.groups, _defineProperty({}, name, members));
      this.setState({
        groups: groups
      });
      this.props.callback(this.props.name, groups);
    }
  }, {
    key: 'render',
    value: function render() {
      var groups = [];
      for (var groupName in this.state.groups) {
        var members = this.state.groups[groupName];
        groups.push(React.createElement(
          'li',
          { key: groupName },
          groupName,
          React.createElement(
            'button',
            { type: 'button', onClick: this.removeGroup.bind(this, groupName) },
            '-'
          ),
          React.createElement(ItemList, { name: groupName, defaultValue: members, callback: this.updateGroup })
        ));
      }

      return React.createElement(
        'div',
        null,
        React.createElement(
          'label',
          { htmlFor: this.props.name },
          this.props.label
        ),
        React.createElement('p', null),
        React.createElement('input', { ref: this.newGroupName }),
        React.createElement(
          'button',
          { type: 'button', onClick: this.addGroup.bind(this) },
          'Add Group'
        ),
        React.createElement(
          'ul',
          null,
          groups
        )
      );
    }
  }]);

  return Groups;
}(React.Component);

var Item = function (_React$Component10) {
  _inherits(Item, _React$Component10);

  function Item(props) {
    _classCallCheck(this, Item);

    return _possibleConstructorReturn(this, (Item.__proto__ || Object.getPrototypeOf(Item)).call(this, props));
  }

  _createClass(Item, [{
    key: 'render',
    value: function render() {
      return React.createElement(
        'div',
        null,
        React.createElement(
          'label',
          { htmlFor: this.props.name },
          this.props.name
        ),
        React.createElement('input', { name: this.props.name, type: this.props.type, defaultValue: this.props.defaultValue, defaultChecked: this.props.defaultChecked })
      );
    }
  }]);

  return Item;
}(React.Component);

var ItemMap = function (_React$Component11) {
  _inherits(ItemMap, _React$Component11);

  function ItemMap(props) {
    _classCallCheck(this, ItemMap);

    var _this12 = _possibleConstructorReturn(this, (ItemMap.__proto__ || Object.getPrototypeOf(ItemMap)).call(this, props));

    _this12.state = {
      item: "",
      value: _this12.props.defaultValue,
      mapItems: [],
      listItems: []
    };

    _this12.add = _this12.add.bind(_this12);
    _this12.remove = _this12.remove.bind(_this12);
    _this12.handleChange = _this12.handleChange.bind(_this12);
    _this12.handleCallback = _this12.handleCallback.bind(_this12);
    return _this12;
  }

  _createClass(ItemMap, [{
    key: 'handleChange',
    value: function handleChange(event) {
      var value = Number(event.target.value);
      this.setState({
        item: value
      });
    }
  }, {
    key: 'handleCallback',
    value: function handleCallback(key, value) {
      var v = Object.assign({}, this.state.value, _defineProperty({}, key, value));
      this.setState({
        value: v
      });
      this.props.callback(this.props.name, v);
    }
  }, {
    key: 'add',
    value: function add() {
      if (!this.state.item) {
        return;
      }
      if (this.state.value && this.state.value[this.state.item] != null) {
        return;
      }

      var value = Object.assign({}, this.state.value, _defineProperty({}, this.state.item, []));
      this.setState({
        value: value
      });
      this.props.callback(this.props.name, value);
    }
  }, {
    key: 'remove',
    value: function remove(id) {
      if (this.state.value == null) {
        return;
      }

      var value = Object.assign({}, this.state.value, _defineProperty({}, id, undefined));
      this.setState({
        value: value
      });
      this.props.callback(this.props.name, value);
    }
  }, {
    key: 'componentWillMount',
    value: function componentWillMount() {
      var _this13 = this;

      this.props.mapItems(function (items) {
        _this13.setState({
          mapItems: items
        });
      });
      this.props.listItems(function (items) {
        _this13.setState({
          listItems: items
        });
      });
    }
  }, {
    key: 'render',
    value: function render() {
      var _this14 = this;

      var rows = [];
      if (this.state.value) {
        var _loop2 = function _loop2(i) {
          if (_this14.state.value[i] === undefined) {
            return 'continue';
          }
          var text = i;
          var matches = _this14.state.mapItems.filter(function (obj) {
            return obj.ID == i;
          });
          if (matches.length > 0) {
            text = matches[0].Display;
          }
          rows.push(React.createElement(
            'li',
            { key: i },
            text,
            React.createElement(
              'button',
              { type: 'button', onClick: _this14.remove.bind(_this14, i) },
              '-'
            ),
            React.createElement(ItemList, { name: i, label: _this14.props.listLabel, type: 'select', listItems: _this14.state.listItems, defaultValue: _this14.state.value[i], callback: _this14.handleCallback })
          ));
        };

        for (var i in this.state.value) {
          var _ret2 = _loop2(i);

          if (_ret2 === 'continue') continue;
        }
      }

      var optionsMap = [];
      // empty selection
      optionsMap.push(React.createElement('option', { disabled: true, key: '', value: '' }));
      for (var i in this.state.mapItems) {
        var option = this.state.mapItems[i];
        // skip already selected
        if (this.state.value && this.state.value[option.ID] != null) {
          continue;
        }
        optionsMap.push(React.createElement(
          'option',
          { key: option.ID, value: option.ID },
          option.Display
        ));
      }

      return React.createElement(
        'div',
        null,
        React.createElement(
          'label',
          null,
          this.props.label
        ),
        React.createElement(
          'ul',
          null,
          rows,
          React.createElement(
            'select',
            { value: this.state.item, onChange: this.handleChange },
            optionsMap
          ),
          React.createElement(
            'button',
            { type: 'button', onClick: this.add },
            '+'
          )
        )
      );
    }
  }]);

  return ItemMap;
}(React.Component);

var ItemList = function (_React$Component12) {
  _inherits(ItemList, _React$Component12);

  function ItemList(props) {
    _classCallCheck(this, ItemList);

    var _this15 = _possibleConstructorReturn(this, (ItemList.__proto__ || Object.getPrototypeOf(ItemList)).call(this, props));

    _this15.state = {
      item: "",
      value: _this15.props.defaultValue
    };

    _this15.add = _this15.add.bind(_this15);
    _this15.remove = _this15.remove.bind(_this15);
    _this15.handleChange = _this15.handleChange.bind(_this15);
    return _this15;
  }

  _createClass(ItemList, [{
    key: 'handleChange',
    value: function handleChange(event) {
      var value = event.target.value;
      if (this.props.type === "select") {
        value = Number(value);
      }
      this.setState({
        item: value
      });
    }
  }, {
    key: 'add',
    value: function add() {
      if (!this.state.item) {
        return;
      }
      if (this.state.value && this.state.value.includes(this.state.item)) {
        return;
      }

      var value = null;
      if (this.state.value == null) {
        value = [this.state.item];
      } else {
        value = [].concat(_toConsumableArray(this.state.value), [this.state.item]);
      }
      this.setState({
        value: value
      });
      this.props.callback(this.props.name, value);
    }
  }, {
    key: 'remove',
    value: function remove(id) {
      if (this.state.value == null) {
        return;
      }

      var value = this.state.value.filter(function (_, index) {
        return index != id;
      });
      this.setState({
        value: value
      });
      this.props.callback(this.props.name, value);
    }
  }, {
    key: 'render',
    value: function render() {
      var _this16 = this;

      var rows = [];
      if (this.state.value) {
        var _loop3 = function _loop3(i) {
          var text = _this16.state.value[i];
          if (_this16.props.type === "select") {
            var _matches = _this16.props.listItems.filter(function (obj) {
              return obj.ID == text;
            });
            if (_matches.length > 0) {
              text = _matches[0].Display;
            }
          }
          rows.push(React.createElement(
            'li',
            { key: i },
            text,
            React.createElement(
              'button',
              { type: 'button', onClick: _this16.remove.bind(_this16, i) },
              '-'
            )
          ));
        };

        for (var i in this.state.value) {
          _loop3(i);
        }
      }

      var input = React.createElement('input', { type: this.props.type, value: this.state.item, onChange: this.handleChange });
      if (this.props.type === "select") {
        var optionsList = [];
        // empty selection
        optionsList.push(React.createElement('option', { disabled: true, key: '', value: '' }));
        for (var i in this.props.listItems) {
          var option = this.props.listItems[i];
          // skip already selected
          if (this.state.value && this.state.value.indexOf(option.ID) != -1) {
            continue;
          }
          optionsList.push(React.createElement(
            'option',
            { key: option.ID, value: option.ID },
            option.Display
          ));
        }
        input = React.createElement(
          'select',
          { value: this.state.item, onChange: this.handleChange },
          optionsList
        );
      }

      return React.createElement(
        'div',
        null,
        React.createElement(
          'label',
          null,
          this.props.label
        ),
        React.createElement(
          'ul',
          null,
          rows,
          input,
          React.createElement(
            'button',
            { type: 'button', onClick: this.add },
            '+'
          )
        )
      );
    }
  }]);

  return ItemList;
}(React.Component);

ReactDOM.render(React.createElement(App, null), document.getElementById('app'));