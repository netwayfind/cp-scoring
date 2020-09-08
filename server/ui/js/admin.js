'use strict';

function ownKeys(object, enumerableOnly) { var keys = Object.keys(object); if (Object.getOwnPropertySymbols) { var symbols = Object.getOwnPropertySymbols(object); if (enumerableOnly) symbols = symbols.filter(function (sym) { return Object.getOwnPropertyDescriptor(object, sym).enumerable; }); keys.push.apply(keys, symbols); } return keys; }

function _objectSpread(target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i] != null ? arguments[i] : {}; if (i % 2) { ownKeys(Object(source), true).forEach(function (key) { _defineProperty(target, key, source[key]); }); } else if (Object.getOwnPropertyDescriptors) { Object.defineProperties(target, Object.getOwnPropertyDescriptors(source)); } else { ownKeys(Object(source)).forEach(function (key) { Object.defineProperty(target, key, Object.getOwnPropertyDescriptor(source, key)); }); } } return target; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

class App extends React.Component {
  constructor() {
    super();
    this.state = {
      authenticated: false,
      page: null,
      id: null,
      lastUpdatedTeams: 0,
      lastUpdatedHosts: 0,
      lastUpdatedTemplates: 0,
      lastUpdatedScenarios: 0,
      lastUpdatedAdministrators: 0
    };
    this.authCallback = this.authCallback.bind(this);
    this.setPage = this.setPage.bind(this);
    this.logout = this.logout.bind(this);
  }

  authCallback(statusCode) {
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

  logout() {
    let url = "/logout";
    fetch(url, {
      credentials: 'same-origin',
      method: "DELETE"
    }).then(function (_) {
      this.setState({
        authenticated: false
      });
    }.bind(this));
  }

  componentDidMount() {
    // check if logged in by visiting the following URL
    let url = "/templates";
    fetch(url, {
      credentials: 'same-origin'
    }).then(function (response) {
      this.authCallback(response.status);
    }.bind(this)); // track which page to be on

    let i = window.location.hash.indexOf('#');
    let hash = window.location.hash.slice(i + 1);
    this.setPage(hash); // handle browser back/forward

    window.onhashchange = e => {
      let i = e.newURL.indexOf('#');
      let hash = e.newURL.slice(i + 1);
      this.setPage(hash);
    };
  }

  setPage(hash) {
    let parts = hash.split("/");
    let page = parts[0];
    let id = null;

    if (parts.length >= 2) {
      id = parts[1];
    }

    this.setState({
      page: page,
      id: id
    });
  }

  updateTeamCallback() {
    this.setState({
      lastUpdatedTeams: Date.now()
    });
  }

  updateHostCallback() {
    this.setState({
      lastUpdatedHosts: Date.now()
    });
  }

  updateTemplateCallback() {
    this.setState({
      lastUpdatedTemplates: Date.now()
    });
  }

  updateScenarioCallback() {
    this.setState({
      lastUpdatedScenarios: Date.now()
    });
  }

  updateAdministratorCallback() {
    this.setState({
      lastUpdatedAdministrators: Date.now()
    });
  }

  render() {
    if (!this.state.authenticated) {
      return /*#__PURE__*/React.createElement("div", {
        className: "App"
      }, /*#__PURE__*/React.createElement(Login, {
        callback: this.authCallback
      }));
    } // reset links to available


    let classes_teams = ["nav-button"];
    let classes_hosts = ["nav-button"];
    let classes_templates = ["nav-button"];
    let classes_scenarios = ["nav-button"];
    let classes_administrators = ["nav-button"]; // default page is empty

    let page = /*#__PURE__*/React.createElement(React.Fragment, null);
    let content = /*#__PURE__*/React.createElement(React.Fragment, null);

    if (this.state.page == "teams") {
      classes_teams.push("nav-button-selected");
      page = /*#__PURE__*/React.createElement(Teams, {
        lastUpdated: this.state.lastUpdatedTeams,
        selected: this.state.id
      });
      content = /*#__PURE__*/React.createElement(TeamEntry, {
        id: this.state.id,
        updateCallback: this.updateTeamCallback.bind(this)
      });
    } else if (this.state.page == "hosts") {
      classes_hosts.push("nav-button-selected");
      page = /*#__PURE__*/React.createElement(Hosts, {
        lastUpdated: this.state.lastUpdatedHosts,
        selected: this.state.id
      });
      content = /*#__PURE__*/React.createElement(HostEntry, {
        id: this.state.id,
        updateCallback: this.updateHostCallback.bind(this)
      });
    } else if (this.state.page == "templates") {
      classes_templates.push("nav-button-selected");
      page = /*#__PURE__*/React.createElement(Templates, {
        lastUpdated: this.state.lastUpdatedTemplates,
        selected: this.state.id
      });
      content = /*#__PURE__*/React.createElement(TemplateEntry, {
        id: this.state.id,
        updateCallback: this.updateTemplateCallback.bind(this)
      });
    } else if (this.state.page == "scenarios") {
      classes_scenarios.push("nav-button-selected");
      page = /*#__PURE__*/React.createElement(Scenarios, {
        lastUpdated: this.state.lastUpdatedScenarios,
        selected: this.state.id
      });
      content = /*#__PURE__*/React.createElement(ScenarioEntry, {
        id: this.state.id,
        updateCallback: this.updateScenarioCallback.bind(this)
      });
    } else if (this.state.page == "administrators") {
      classes_administrators.push("nav-button-selected");
      page = /*#__PURE__*/React.createElement(Administrators, {
        lastUpdated: this.state.lastUpdatedAdministrators,
        selected: this.state.id
      });
      content = /*#__PURE__*/React.createElement(AdministratorEntry, {
        username: this.state.id,
        updateCallback: this.updateAdministratorCallback.bind(this)
      });
    }

    return /*#__PURE__*/React.createElement("div", {
      className: "App"
    }, /*#__PURE__*/React.createElement("div", {
      className: "heading"
    }, /*#__PURE__*/React.createElement("h1", null, "cp-scoring")), /*#__PURE__*/React.createElement("div", {
      className: "navbar"
    }, /*#__PURE__*/React.createElement("a", {
      className: classes_teams.join(" "),
      href: "#teams"
    }, "Teams"), /*#__PURE__*/React.createElement("a", {
      className: classes_hosts.join(" "),
      href: "#hosts"
    }, "Hosts"), /*#__PURE__*/React.createElement("a", {
      className: classes_templates.join(" "),
      href: "#templates"
    }, "Templates"), /*#__PURE__*/React.createElement("a", {
      className: classes_scenarios.join(" "),
      href: "#scenarios"
    }, "Scenarios"), /*#__PURE__*/React.createElement("a", {
      className: classes_administrators.join(" "),
      href: "#administrators"
    }, "Administrators"), /*#__PURE__*/React.createElement("div", {
      className: "right"
    }, /*#__PURE__*/React.createElement("button", {
      onClick: this.logout
    }, "Logout"))), /*#__PURE__*/React.createElement("div", {
      className: "toc"
    }, page), /*#__PURE__*/React.createElement("div", {
      className: "content"
    }, content));
  }

}

class Listing extends React.Component {
  constructor(props) {
    super(props);
    this.itemUrl = props.itemUrl;
    this.state = {
      error: null,
      selected: null,
      items: []
    };
    this.saveSelected = this.saveSelected.bind(this);
  }

  componentDidMount() {
    this.populate();
    this.saveSelected(this.props.selected);
  }

  componentWillReceiveProps(newProps) {
    this.populate();
    this.saveSelected(newProps.selected);
  }

  saveSelected(selected) {
    this.setState({
      selected: selected
    });
  }

  populate() {
    fetch(this.itemUrl, {
      credentials: 'same-origin'
    }).then(async function (response) {
      if (response.status === 200) {
        let data = await response.json();
        return {
          error: null,
          items: data
        };
      }

      let text = await response.text();
      return {
        error: text
      };
    }).then(function (s) {
      this.setState(s);
    }.bind(this));
  }

}

class Administrators extends Listing {
  constructor(props) {
    props.itemUrl = "/admins";
    super(props);
  }

  render() {
    let rows = [];

    for (let i = 0; i < this.state.items.length; i++) {
      let administrator = this.state.items[i];
      let classes = ["nav-button"];

      if (this.state.selected === administrator) {
        classes.push("nav-button-selected");
      }

      rows.push( /*#__PURE__*/React.createElement("li", {
        key: i
      }, /*#__PURE__*/React.createElement("a", {
        className: classes.join(" "),
        href: "#administrators/" + administrator
      }, administrator)));
    }

    return /*#__PURE__*/React.createElement("div", {
      className: "Admins"
    }, /*#__PURE__*/React.createElement("strong", null, "Administrators"), /*#__PURE__*/React.createElement(Error, {
      message: this.state.error
    }), /*#__PURE__*/React.createElement("ul", null, rows));
  }

}

class AdministratorEntry extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      user: {}
    };
  }

  componentDidMount() {
    if (this.props.username) {
      this.setState({
        user: {
          Username: this.props.username,
          Password: ""
        }
      });
    }
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.username != nextProps.username) {
      this.setState({
        user: {
          Username: nextProps.username,
          Password: ""
        }
      });
    }
  }

  newAdministrator() {
    this.setState({
      user: {
        Username: "",
        Password: ""
      }
    });
    window.location.href = "#administrators";
  }

  updateAdministrator(event) {
    let value = event.target.value;
    this.setState({
      user: _objectSpread(_objectSpread({}, this.state.user), {}, {
        [event.target.name]: value
      })
    });
  }

  saveAdministrator(event) {
    event.preventDefault();
    var url = "/admins"; // for updating existing user

    if (this.props.username != null) {
      url += "/" + this.props.username;
    }

    fetch(url, {
      credentials: 'same-origin',
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(this.state.user)
    }).then(async function (response) {
      if (response.status === 200) {
        this.props.updateCallback();
        let username = this.state.user.Username;

        if (username != undefined && username != null) {
          window.location.href = "#administrators/" + username;
        }

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
    ;
  }

  deleteAdministrator(username) {
    var url = "/admins/" + username;
    fetch(url, {
      credentials: 'same-origin',
      method: 'DELETE'
    }).then(async function (response) {
      if (response.status === 200) {
        this.props.updateCallback();
        window.location.href = "#administrators";
        return {
          error: null,
          user: {}
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
    let content = null;

    if (Object.entries(this.state.user).length != 0) {
      // disable change existing user name
      let existingUser = true;

      if (this.props.username === null) {
        existingUser = false;
      }

      content = /*#__PURE__*/React.createElement("form", {
        onChange: this.updateAdministrator.bind(this),
        onSubmit: this.saveAdministrator.bind(this)
      }, /*#__PURE__*/React.createElement(Item, {
        name: "Username",
        value: this.state.user.Username || "",
        disabled: existingUser
      }), /*#__PURE__*/React.createElement(Item, {
        name: "Password",
        type: "password",
        value: this.state.user.Password
      }), /*#__PURE__*/React.createElement("p", null), /*#__PURE__*/React.createElement("button", {
        type: "submit"
      }, "Submit"), /*#__PURE__*/React.createElement("button", {
        class: "right",
        type: "button",
        disabled: !this.state.user.Username,
        onClick: this.deleteAdministrator.bind(this, this.state.user.Username)
      }, "Delete"));
    }

    return /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("button", {
      type: "button",
      onClick: this.newAdministrator.bind(this)
    }, "New Administrator"), /*#__PURE__*/React.createElement("hr", null), /*#__PURE__*/React.createElement(Error, {
      message: this.state.error
    }), content);
  }

}

class Teams extends Listing {
  constructor(props) {
    props.itemUrl = "/teams";
    super(props);
  }

  render() {
    let rows = [];

    for (let i = 0; i < this.state.items.length; i++) {
      let team = this.state.items[i];
      let classes = ["nav-button"];

      if (this.state.selected === team.ID.toString()) {
        classes.push("nav-button-selected");
      }

      rows.push( /*#__PURE__*/React.createElement("li", {
        key: team.ID
      }, /*#__PURE__*/React.createElement("a", {
        className: classes.join(" "),
        href: "#teams/" + team.ID
      }, "[", team.ID, "] ", team.Name)));
    }

    return /*#__PURE__*/React.createElement("div", {
      className: "Teams"
    }, /*#__PURE__*/React.createElement("strong", null, "Teams"), /*#__PURE__*/React.createElement(Error, {
      message: this.state.error
    }), /*#__PURE__*/React.createElement("ul", null, rows));
  }

}

class TeamEntry extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      team: {}
    };
    this.getTeam = this.getTeam.bind(this);
  }

  static newKey() {
    let key = "";

    for (let i = 0; i < 8; i++) {
      key += Math.floor(Math.random() * 16).toString(16).toUpperCase();
    }

    return key;
  }

  componentDidMount() {
    if (this.props.id) {
      this.getTeam(this.props.id);
    }
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.id != nextProps.id) {
      this.getTeam(nextProps.id);
    }
  }

  newTeam() {
    this.setState({
      team: {
        Name: "",
        POC: "",
        Email: "",
        Enabled: true,
        Key: TeamEntry.newKey()
      }
    });
    window.location.href = "#teams";
  }

  getTeam(id) {
    if (id === null || id === undefined) {
      return;
    }

    let url = "/teams/" + id;
    fetch(url, {
      credentials: 'same-origin'
    }).then(async function (response) {
      if (response.status === 200) {
        let data = await response.json();
        return {
          error: null,
          team: data
        };
      }

      let text = await response.text();
      return {
        error: text
      };
    }).then(function (s) {
      this.setState(s);
    }.bind(this));
  }

  updateTeam(event) {
    let value = event.target.value;

    if (event.target.type == 'checkbox') {
      value = event.target.checked;
    }

    this.setState({
      team: _objectSpread(_objectSpread({}, this.state.team), {}, {
        [event.target.name]: value
      })
    });
  }

  deleteTeam(id) {
    var url = "/teams/" + id;
    fetch(url, {
      credentials: 'same-origin',
      method: 'DELETE'
    }).then(async function (response) {
      if (response.status === 200) {
        this.props.updateCallback();
        window.location.href = "#teams";
        return {
          error: null,
          team: {}
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

  saveTeam(event) {
    event.preventDefault();
    var url = "/teams";

    if (this.state.team.ID != null) {
      url += "/" + this.state.team.ID;
    }

    fetch(url, {
      credentials: 'same-origin',
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(this.state.team)
    }).then(async function (response) {
      if (response.status === 200) {
        this.props.updateCallback();

        if (this.state.team.ID === null || this.state.team.ID === undefined) {
          // for new teams, response should be team ID
          let id = await response.text();
          window.location.href = "#teams/" + id;
        }

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

  regenKey() {
    let key = TeamEntry.newKey();
    this.setState({
      team: _objectSpread(_objectSpread({}, this.state.team), {}, {
        Key: key
      })
    });
  }

  render() {
    let content = null;

    if (Object.entries(this.state.team).length != 0) {
      content = /*#__PURE__*/React.createElement("form", {
        onChange: this.updateTeam.bind(this),
        onSubmit: this.saveTeam.bind(this)
      }, /*#__PURE__*/React.createElement("label", {
        htmlFor: "ID"
      }, "ID"), /*#__PURE__*/React.createElement("input", {
        disabled: true,
        value: this.state.team.ID || ""
      }), /*#__PURE__*/React.createElement(Item, {
        name: "Name",
        value: this.state.team.Name
      }), /*#__PURE__*/React.createElement(Item, {
        name: "POC",
        value: this.state.team.POC
      }), /*#__PURE__*/React.createElement(Item, {
        name: "Email",
        type: "email",
        value: this.state.team.Email
      }), /*#__PURE__*/React.createElement(Item, {
        name: "Enabled",
        type: "checkbox",
        checked: !!this.state.team.Enabled
      }), /*#__PURE__*/React.createElement("br", null), /*#__PURE__*/React.createElement("details", null, /*#__PURE__*/React.createElement("summary", null, "Key"), /*#__PURE__*/React.createElement("ul", null, /*#__PURE__*/React.createElement("li", null, this.state.team.Key, /*#__PURE__*/React.createElement("br", null), /*#__PURE__*/React.createElement("button", {
        type: "button",
        onClick: this.regenKey.bind(this)
      }, "Regenerate")))), /*#__PURE__*/React.createElement("br", null), /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("button", {
        type: "submit"
      }, "Save"), /*#__PURE__*/React.createElement("button", {
        class: "right",
        type: "button",
        disabled: !this.state.team.ID,
        onClick: this.deleteTeam.bind(this, this.state.team.ID)
      }, "Delete")));
    }

    return /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("button", {
      type: "button",
      onClick: this.newTeam.bind(this)
    }, "New Team"), /*#__PURE__*/React.createElement("hr", null), /*#__PURE__*/React.createElement(Error, {
      message: this.state.error
    }), content);
  }

}

class Scenarios extends Listing {
  constructor(props) {
    props.itemUrl = "/scenarios";
    super(props);
  }

  render() {
    let rows = [];

    for (let i = 0; i < this.state.items.length; i++) {
      let scenario = this.state.items[i];
      let classes = ["nav-button"];

      if (this.state.selected === scenario.ID.toString()) {
        classes.push("nav-button-selected");
      }

      rows.push( /*#__PURE__*/React.createElement("li", {
        key: scenario.ID
      }, /*#__PURE__*/React.createElement("a", {
        className: classes.join(" "),
        href: "#scenarios/" + scenario.ID
      }, scenario.Name)));
    }

    return /*#__PURE__*/React.createElement("div", {
      className: "Scenarios"
    }, /*#__PURE__*/React.createElement("strong", null, "Scenarios"), /*#__PURE__*/React.createElement(Error, {
      message: this.state.error
    }), /*#__PURE__*/React.createElement("ul", null, rows));
  }

}

class ScenarioEntry extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      scenario: {}
    };
  }

  componentDidMount() {
    if (this.props.id) {
      this.getScenario(this.props.id);
    }
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.id != nextProps.id) {
      this.getScenario(nextProps.id);
    }
  }

  newScenario() {
    this.setState({
      scenario: {
        Name: "",
        Description: "",
        Enabled: true,
        HostTemplates: {}
      }
    });
    window.location.href = "#scenarios";
  }

  newScenarioFromExisting() {
    this.state.scenario.ID = null;
    this.state.scenario.Name += " (copy " + new Date().toLocaleString() + ")";
    this.saveScenario(null);
  }

  getScenario(id) {
    if (id === null || id === undefined) {
      return;
    }

    let url = "/scenarios/" + id;
    fetch(url, {
      credentials: 'same-origin'
    }).then(async function (response) {
      if (response.status === 200) {
        let data = await response.json();
        return {
          error: null,
          scenario: data
        };
      }

      let text = await response.text();
      return {
        error: text
      };
    }).then(function (s) {
      this.setState(s);
    }.bind(this));
  }

  updateScenario(event) {
    let value = event.target.value;

    if (event.target.type == 'checkbox') {
      value = event.target.checked;
    }

    this.setState({
      scenario: _objectSpread(_objectSpread({}, this.state.scenario), {}, {
        [event.target.name]: value
      })
    });
  }

  saveScenario(event) {
    if (event != null) {
      event.preventDefault();
    }

    var url = "/scenarios";

    if (this.state.scenario.ID != null) {
      url += "/" + this.state.scenario.ID;
    }

    fetch(url, {
      credentials: 'same-origin',
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(this.state.scenario)
    }).then(async function (response) {
      if (response.status === 200) {
        this.props.updateCallback();

        if (this.state.scenario.ID === null || this.state.scenario.ID === undefined) {
          // for new scenarios, response should be scenario ID
          let id = await response.text();
          window.location.href = "#scenarios/" + id;
        }

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

  deleteScenario(id) {
    var url = "/scenarios/" + id;
    fetch(url, {
      credentials: 'same-origin',
      method: 'DELETE'
    }).then(async function (response) {
      if (response.status === 200) {
        this.props.updateCallback();
        window.location.href = "#scenarios";
        return {
          error: null,
          scenario: {}
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

  handleCallback(key, value) {
    this.setState({
      scenario: _objectSpread(_objectSpread({}, this.state.scenario), {}, {
        [key]: value
      })
    });
  }

  mapItems(callback) {
    var url = "/hosts";
    fetch(url, {
      credentials: 'same-origin'
    }).then(async function (response) {
      if (response.status === 200) {
        let data = await response.json();
        let items = data.map(function (host) {
          return {
            ID: host.ID,
            Display: host.Hostname
          };
        });
        callback(items);
        return;
      }

      let text = await response.text();
      this.setState({
        error: text
      });
    }.bind(this));
  }

  listItems(callback) {
    var url = "/templates";
    fetch(url, {
      credentials: 'same-origin'
    }).then(async function (response) {
      if (response.status === 200) {
        let data = await response.json();
        let items = data.map(function (template) {
          return {
            ID: template.ID,
            Display: template.Name
          };
        });
        callback(items);
        return;
      }

      let text = await response.text();
      this.setState({
        error: text
      });
    }.bind(this));
  }

  render() {
    let content = null;

    if (Object.entries(this.state.scenario).length != 0) {
      content = /*#__PURE__*/React.createElement("form", {
        onChange: this.updateScenario.bind(this),
        onSubmit: this.saveScenario.bind(this)
      }, /*#__PURE__*/React.createElement("label", {
        htmlFor: "ID"
      }, "ID"), /*#__PURE__*/React.createElement("input", {
        disabled: true,
        value: this.state.scenario.ID || ""
      }), /*#__PURE__*/React.createElement(Item, {
        name: "Name",
        value: this.state.scenario.Name
      }), /*#__PURE__*/React.createElement("label", {
        htmlFor: "Description"
      }, "Description"), /*#__PURE__*/React.createElement("textarea", {
        name: "Description",
        rows: "10",
        cols: "80",
        value: this.state.scenario.Description
      }), /*#__PURE__*/React.createElement(Item, {
        name: "Enabled",
        type: "checkbox",
        checked: !!this.state.scenario.Enabled
      }), /*#__PURE__*/React.createElement(ItemMap, {
        name: "HostTemplates",
        label: "Hosts",
        listLabel: "Templates",
        value: this.state.scenario.HostTemplates,
        callback: this.handleCallback.bind(this),
        mapItems: this.mapItems,
        listItems: this.listItems
      }), /*#__PURE__*/React.createElement("br", null), /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("button", {
        type: "submit"
      }, "Save"), /*#__PURE__*/React.createElement("button", {
        class: "right",
        type: "button",
        disabled: !this.state.scenario.ID,
        onClick: this.deleteScenario.bind(this, this.state.scenario.ID)
      }, "Delete")));
    }

    return /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("button", {
      onClick: this.newScenario.bind(this)
    }, "New Scenario"), /*#__PURE__*/React.createElement("button", {
      onClick: this.newScenarioFromExisting.bind(this)
    }, "Clone"), /*#__PURE__*/React.createElement("hr", null), /*#__PURE__*/React.createElement(Error, {
      message: this.state.error
    }), content);
  }

}

class Hosts extends Listing {
  constructor(props) {
    props.itemUrl = "/hosts";
    super(props);
  }

  render() {
    let rows = [];

    for (let i = 0; i < this.state.items.length; i++) {
      let host = this.state.items[i];
      let classes = ["nav-button"];

      if (this.state.selected === host.ID.toString()) {
        classes.push("nav-button-selected");
      }

      rows.push( /*#__PURE__*/React.createElement("li", {
        key: host.ID
      }, /*#__PURE__*/React.createElement("a", {
        className: classes.join(" "),
        href: "#hosts/" + host.ID
      }, host.Hostname, " - ", host.OS)));
    }

    return /*#__PURE__*/React.createElement("div", {
      className: "Hosts"
    }, /*#__PURE__*/React.createElement("strong", null, "Hosts"), /*#__PURE__*/React.createElement(Error, {
      message: this.state.error
    }), /*#__PURE__*/React.createElement("ul", null, rows));
  }

}

class HostEntry extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      host: {}
    };
  }

  componentDidMount() {
    this.getHost(this.props.id);
  }

  componentWillReceiveProps(newProps) {
    if (this.props.id != newProps.id) {
      this.getHost(newProps.id);
    }
  }

  newHost() {
    this.setState({
      error: null,
      host: {
        Hostname: "",
        OS: ""
      }
    });
    window.location.href = "#hosts";
  }

  newHostFromExisting() {
    this.state.host.ID = null;
    this.state.host.Hostname += " (copy " + new Date().toLocaleString() + ")";
    this.saveHost(null);
  }

  getHost(id) {
    if (id === null || id === undefined) {
      return;
    }

    let url = "/hosts/" + id;
    fetch(url, {
      credentials: 'same-origin'
    }).then(async function (response) {
      if (response.status === 200) {
        let data = await response.json();
        return {
          error: null,
          host: data
        };
      }

      let text = await response.text();
      return {
        error: text
      };
    }).then(function (s) {
      this.setState(s);
    }.bind(this));
  }

  updateHost(event) {
    let value = event.target.value;

    if (event.target.type == 'checkbox') {
      value = event.target.checked;
    }

    this.setState({
      host: _objectSpread(_objectSpread({}, this.state.host), {}, {
        [event.target.name]: value
      })
    });
  }

  saveHost(event) {
    if (event != null) {
      event.preventDefault();
    }

    var url = "/hosts";

    if (this.state.host.ID != null) {
      url += "/" + this.state.host.ID;
    }

    fetch(url, {
      credentials: 'same-origin',
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(this.state.host)
    }).then(async function (response) {
      if (response.status === 200) {
        this.props.updateCallback();

        if (this.state.host.ID === null || this.state.host.ID === undefined) {
          // for new hosts, response should be host ID
          let id = await response.text();
          window.location.href = "#hosts/" + id;
        }

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

  deleteHost(id) {
    var url = "/hosts/" + id;
    fetch(url, {
      credentials: 'same-origin',
      method: 'DELETE'
    }).then(async function (response) {
      if (response.status === 200) {
        this.props.updateCallback();
        window.location.href = "#hosts";
        return {
          error: null,
          host: {}
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
    let content = null;

    if (Object.entries(this.state.host).length != 0) {
      content = /*#__PURE__*/React.createElement("form", {
        onChange: this.updateHost.bind(this),
        onSubmit: this.saveHost.bind(this)
      }, /*#__PURE__*/React.createElement("label", {
        htmlFor: "ID"
      }, "ID"), /*#__PURE__*/React.createElement("input", {
        disabled: true,
        value: this.state.host.ID || ""
      }), /*#__PURE__*/React.createElement(Item, {
        name: "Hostname",
        type: "text",
        value: this.state.host.Hostname
      }), /*#__PURE__*/React.createElement(Item, {
        name: "OS",
        type: "text",
        value: this.state.host.OS
      }), /*#__PURE__*/React.createElement("br", null), /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("button", {
        type: "submit"
      }, "Save"), /*#__PURE__*/React.createElement("button", {
        class: "right",
        type: "button",
        disabled: !this.state.host.ID,
        onClick: this.deleteHost.bind(this, this.state.host.ID)
      }, "Delete")));
    }

    return /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("button", {
      type: "button",
      onClick: this.newHost.bind(this)
    }, "New Host"), /*#__PURE__*/React.createElement("button", {
      type: "button",
      onClick: this.newHostFromExisting.bind(this)
    }, "Clone"), /*#__PURE__*/React.createElement("hr", null), /*#__PURE__*/React.createElement(Error, {
      message: this.state.error
    }), content);
  }

}

class Templates extends Listing {
  constructor(props) {
    props.itemUrl = "/templates";
    super(props);
  }

  render() {
    let rows = [];

    for (let i = 0; i < this.state.items.length; i++) {
      let template = this.state.items[i];
      let classes = ["nav-button"];

      if (this.state.selected === template.ID.toString()) {
        classes.push("nav-button-selected");
      }

      rows.push( /*#__PURE__*/React.createElement("li", {
        key: template.ID
      }, /*#__PURE__*/React.createElement("a", {
        className: classes.join(" "),
        href: "#templates/" + template.ID
      }, template.Name)));
    }

    return /*#__PURE__*/React.createElement("div", {
      className: "Templates"
    }, /*#__PURE__*/React.createElement("strong", null, "Templates"), /*#__PURE__*/React.createElement(Error, {
      message: this.state.error
    }), /*#__PURE__*/React.createElement("ul", null, rows));
  }

}

class TemplateEntry extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      template: {}
    };
  }

  componentDidMount() {
    if (this.props.id) {
      this.getTemplate(this.props.id);
    }
  }

  componentWillReceiveProps(newProps) {
    if (this.props.id != newProps.id) {
      this.getTemplate(newProps.id);
    }
  }

  newTemplate() {
    this.setState({
      template: {
        Name: "",
        State: {}
      }
    });
    window.location.href = "#templates";
  }

  newTemplateFromState() {
    this.setState({
      template: {
        TemplateName: "",
        StateID: ""
      }
    });
    window.location.href = "#templates";
  }

  newTemplateFromExisting() {
    this.state.template.ID = null;
    this.state.template.Name += " (copy " + new Date().toLocaleString() + ")";
    this.saveTemplate(null, "/templates");
  }

  getTemplate(id) {
    if (id === null || id === undefined) {
      return;
    }

    let url = "/templates/" + id;
    fetch(url, {
      credentials: 'same-origin'
    }).then(async function (response) {
      if (response.status === 200) {
        let data = await response.json();
        return {
          error: null,
          template: data
        };
      }

      let text = await response.text();
      return {
        error: text
      };
    }).then(function (s) {
      this.setState(s);
    }.bind(this));
  }

  updateTemplate(event) {
    let value = event.target.value;

    if (event.target.type == 'checkbox') {
      value = event.target.checked;
    }

    this.setState({
      template: _objectSpread(_objectSpread({}, this.state.template), {}, {
        [event.target.name]: value
      })
    });
  }

  saveRegularTemplate(event) {
    this.saveTemplate(event, "/templates");
  }

  saveStateTemplate(event) {
    this.saveTemplate(event, "/templates/state");
  }

  saveTemplate(event, url) {
    if (event != null) {
      event.preventDefault();
    }

    if (this.state.template.ID != null) {
      url += "/" + this.state.template.ID;
    }

    fetch(url, {
      credentials: 'same-origin',
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(this.state.template)
    }).then(async function (response) {
      if (response.status === 200) {
        this.props.updateCallback();

        if (this.state.template.ID === null || this.state.template.ID === undefined) {
          // for new templates, response should be template ID
          let id = await response.text();
          window.location.href = "#templates/" + id;
        }

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

  deleteTemplate(id) {
    var url = "/templates/" + id;
    fetch(url, {
      credentials: 'same-origin',
      method: 'DELETE'
    }).then(async function (response) {
      if (response.status === 200) {
        this.props.updateCallback();
        window.location.href = "#templates";
        return {
          error: null,
          template: {}
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

  handleCallback(key, value) {
    let state = _objectSpread(_objectSpread({}, this.state.template.State), {}, {
      [key]: value
    });

    this.setState({
      template: _objectSpread(_objectSpread({}, this.state.template), {}, {
        State: state
      })
    });
  }

  render() {
    let content = null; // template from state

    if (this.state.template.StateID != undefined) {
      content = /*#__PURE__*/React.createElement("form", {
        onChange: this.updateTemplate.bind(this),
        onSubmit: this.saveStateTemplate.bind(this)
      }, /*#__PURE__*/React.createElement("label", {
        htmlFor: "TemplateName"
      }, "Template Name"), /*#__PURE__*/React.createElement("input", {
        name: "TemplateName",
        value: this.state.template.TemplateName || ""
      }), /*#__PURE__*/React.createElement("br", null), /*#__PURE__*/React.createElement("label", {
        htmlFor: "StateID"
      }, "State ID"), /*#__PURE__*/React.createElement("input", {
        name: "StateID",
        value: this.state.template.StateID || ""
      }), /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("button", {
        type: "submit"
      }, "Save")));
    } // regular template
    else if (Object.entries(this.state.template).length != 0) {
        content = /*#__PURE__*/React.createElement("form", {
          onChange: this.updateTemplate.bind(this),
          onSubmit: this.saveRegularTemplate.bind(this)
        }, /*#__PURE__*/React.createElement("label", {
          htmlFor: "ID"
        }, "ID"), /*#__PURE__*/React.createElement("input", {
          disabled: true,
          value: this.state.template.ID || ""
        }), /*#__PURE__*/React.createElement(Item, {
          name: "Name",
          type: "text",
          value: this.state.template.Name
        }), /*#__PURE__*/React.createElement(Users, {
          users: this.state.template.State.Users,
          callback: this.handleCallback.bind(this)
        }), /*#__PURE__*/React.createElement(Groups, {
          groups: this.state.template.State.Groups,
          callback: this.handleCallback.bind(this)
        }), /*#__PURE__*/React.createElement(Processes, {
          processes: this.state.template.State.Processes,
          callback: this.handleCallback.bind(this)
        }), /*#__PURE__*/React.createElement(Software, {
          software: this.state.template.State.Software,
          callback: this.handleCallback.bind(this)
        }), /*#__PURE__*/React.createElement(NetworkConnections, {
          conns: this.state.template.State.NetworkConnections,
          callback: this.handleCallback.bind(this)
        }), /*#__PURE__*/React.createElement(ScheduledTasks, {
          tasks: this.state.template.State.ScheduledTasks,
          callback: this.handleCallback.bind(this)
        }), /*#__PURE__*/React.createElement(WindowsFirewallProfiles, {
          profiles: this.state.template.State.WindowsFirewallProfiles,
          callback: this.handleCallback.bind(this)
        }), /*#__PURE__*/React.createElement(WindowsFirewallRules, {
          rules: this.state.template.State.WindowsFirewallRules,
          callback: this.handleCallback.bind(this)
        }), /*#__PURE__*/React.createElement(WindowsSettings, {
          settings: this.state.template.State.WindowsSettings,
          callback: this.handleCallback.bind(this)
        }), /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("button", {
          type: "submit"
        }, "Save"), /*#__PURE__*/React.createElement("button", {
          class: "right",
          type: "button",
          disabled: !this.state.template.ID,
          onClick: this.deleteTemplate.bind(this, this.state.template.ID)
        }, "Delete")));
      }

    return /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("button", {
      type: "button",
      onClick: this.newTemplate.bind(this)
    }, "New Template"), /*#__PURE__*/React.createElement("button", {
      type: "button",
      onClick: this.newTemplateFromState.bind(this)
    }, "From State"), /*#__PURE__*/React.createElement("button", {
      type: "button",
      onClick: this.newTemplateFromExisting.bind(this)
    }, "Clone"), /*#__PURE__*/React.createElement("hr", null), /*#__PURE__*/React.createElement(Error, {
      message: this.state.error
    }), content);
  }

}

class ObjectState extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("select", {
      value: this.props.value,
      onChange: this.props.onChange
    }, /*#__PURE__*/React.createElement("option", null, "Add"), /*#__PURE__*/React.createElement("option", null, "Keep"), /*#__PURE__*/React.createElement("option", null, "Remove")));
  }

}

class Users extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      users: []
    };
  }

  componentDidMount() {
    this.setUsers(this.props.users);
  }

  componentWillReceiveProps(newProps) {
    this.setUsers(newProps.users);
  }

  setUsers(users) {
    if (users === undefined || users === null) {
      users = [];
    }

    this.setState({
      users: users
    });
  }

  addUser() {
    let empty = {
      Name: "",
      ObjectState: "Keep",
      AccountActive: true,
      PasswordExpires: true,
      // unix timestamp in seconds
      PasswordLastSet: Math.trunc(Date.now() / 1000)
    };
    let users = [...this.state.users, empty];
    this.setState({
      users: users
    });
    this.props.callback("Users", users);
  }

  removeUser(id) {
    let users = this.state.users.filter(function (_, index) {
      return index != id;
    });
    this.setState({
      users: users
    });
    this.props.callback("Users", users);
  }

  updateUser(id, field, event) {
    let updated = this.state.users;
    let value = event.target.value;

    if (event.target.type === "checkbox") {
      if (event.target.checked) {
        value = true;
      } else {
        value = false;
      }
    } else if (event.target.type === "date") {
      let parts = event.target.value.split("-");

      if (parts.length != 3) {
        return;
      }

      let current = new Date(Math.trunc(this.state.users[id].PasswordLastSet * 1000));
      current.setFullYear(parts[0]); // months start counting at 0

      current.setMonth(parts[1] - 1);
      current.setDate(parts[2]);
      value = Math.trunc(current.getTime() / 1000);

      if (Number.isNaN(value)) {
        return;
      }
    } else if (event.target.type === "time") {
      let parts = event.target.value.split(":");

      if (parts.length != 3) {
        return;
      }

      let current = new Date(Math.trunc(this.state.users[id].PasswordLastSet * 1000));
      current.setHours(parts[0]);
      current.setMinutes(parts[1]);
      current.setSeconds(parts[2]);
      value = Math.trunc(current.getTime() / 1000);

      if (Number.isNaN(value)) {
        return;
      }
    }

    updated[id] = _objectSpread(_objectSpread({}, updated[id]), {}, {
      [field]: value
    });
    this.setState({
      users: updated
    });
    this.props.callback("Users", updated);
  }

  render() {
    let users = [];

    for (let i = 0; i < this.state.users.length; i++) {
      let user = this.state.users[i];
      let d = new Date(user.PasswordLastSet * 1000);
      let passwordLastSetDate = ("000" + d.getFullYear()).slice(-4);
      passwordLastSetDate += "-";
      passwordLastSetDate += ("0" + (d.getMonth() + 1)).slice(-2);
      passwordLastSetDate += "-";
      passwordLastSetDate += ("0" + d.getDate()).slice(-2);
      let passwordLastSetTime = ("000" + d.getHours()).slice(-2);
      passwordLastSetTime += ":";
      passwordLastSetTime += ("000" + d.getMinutes()).slice(-2);
      passwordLastSetTime += ":";
      passwordLastSetTime += ("000" + d.getSeconds()).slice(-2);
      let userOptions = null;

      if (user.ObjectState != "Remove") {
        userOptions = /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
          type: "checkbox",
          checked: user.AccountActive,
          onChange: event => this.updateUser(i, "AccountActive", event)
        })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
          type: "checkbox",
          checked: user.PasswordExpires,
          onChange: event => this.updateUser(i, "PasswordExpires", event)
        })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
          type: "date",
          value: passwordLastSetDate,
          onChange: event => this.updateUser(i, "PasswordLastSet", event)
        }), /*#__PURE__*/React.createElement("input", {
          type: "time",
          value: passwordLastSetTime,
          onChange: event => this.updateUser(i, "PasswordLastSet", event)
        })));
      }

      users.push( /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
        type: "text",
        value: user.Name,
        onChange: event => this.updateUser(i, "Name", event)
      })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(ObjectState, {
        value: user.ObjectState,
        onChange: event => this.updateUser(i, "ObjectState", event)
      })), userOptions, /*#__PURE__*/React.createElement("button", {
        type: "button",
        onClick: this.removeUser.bind(this, i)
      }, "-")));
    }

    return /*#__PURE__*/React.createElement("details", null, /*#__PURE__*/React.createElement("summary", null, "Users"), /*#__PURE__*/React.createElement("button", {
      type: "button",
      onClick: this.addUser.bind(this)
    }, "Add User"), /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Name"), /*#__PURE__*/React.createElement("th", null, "State"), /*#__PURE__*/React.createElement("th", null, "Active"), /*#__PURE__*/React.createElement("th", null, "Password Expires"), /*#__PURE__*/React.createElement("th", null, "Password Last Set"), /*#__PURE__*/React.createElement("th", null)), users));
  }

}

class Groups extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      groups: []
    };
    this.newGroupName = React.createRef();
    this.addGroup = this.addGroup.bind(this);
    this.addGroupMember = this.addGroupMember.bind(this);
    this.removeGroup = this.removeGroup.bind(this);
    this.removeGroupMember = this.removeGroupMember.bind(this);
    this.updateGroupMember = this.updateGroupMember.bind(this);
  }

  componentDidMount() {
    this.setGroups(this.props.groups);
  }

  componentWillReceiveProps(newProps) {
    this.setGroups(newProps.groups);
  }

  setGroups(groups) {
    if (groups === undefined || groups === null) {
      groups = {};
    }

    this.setState({
      groups: groups
    });
  }

  addGroup() {
    if (this.newGroupName.current === null) {
      return;
    } else if (this.newGroupName.current.value.length === 0) {
      return;
    }

    let groups = _objectSpread(_objectSpread({}, this.state.groups), {}, {
      [this.newGroupName.current.value]: []
    });

    this.setState({
      groups: groups
    });
    this.newGroupName.current.value = "";
    this.props.callback("Groups", groups);
  }

  addGroupMember(groupName) {
    let group = this.state.groups[groupName];
    group.push({
      Name: "",
      ObjectState: "Keep"
    });

    let groups = _objectSpread(_objectSpread({}, this.state.groups), {}, {
      [groupName]: group
    });

    this.setState({
      groups: groups
    });
    this.props.callback("Groups", groups);
  }

  removeGroup(groupName) {
    let groups = this.state.groups;
    delete groups[groupName];
    this.setState({
      groups: groups
    });
    this.props.callback("Groups", groups);
  }

  removeGroupMember(groupName, memberIndex) {
    let group = this.state.groups[groupName];
    group.splice(memberIndex, 1);

    let groups = _objectSpread(_objectSpread({}, this.state.groups), {}, {
      [groupName]: group
    });

    this.setState({
      groups: groups
    });
    this.props.callback("Groups", groups);
  }

  updateGroupMember(groupName, memberIndex, key, value) {
    let group = this.state.groups[groupName];
    let member = group[memberIndex];
    member[key] = value;

    let groups = _objectSpread(_objectSpread({}, this.state.groups), {}, {
      [groupName]: group
    });

    this.setState({
      groups: groups
    });
    this.props.callback("Groups", groups);
  }

  render() {
    let groups = [];

    for (let groupName in this.state.groups) {
      let groupMembers = [];

      for (let i in this.state.groups[groupName]) {
        let member = this.state.groups[groupName][i];
        groupMembers.push( /*#__PURE__*/React.createElement("details", {
          key: i
        }, /*#__PURE__*/React.createElement("summary", null, member.Name), /*#__PURE__*/React.createElement("button", {
          type: "button",
          onClick: this.removeGroupMember.bind(this, groupName, i)
        }, "-"), /*#__PURE__*/React.createElement("ul", null, /*#__PURE__*/React.createElement("li", null, /*#__PURE__*/React.createElement("label", null, "Name"), /*#__PURE__*/React.createElement("input", {
          type: "text",
          value: member.Name,
          onChange: event => this.updateGroupMember(groupName, i, "Name", event.target.value)
        })), /*#__PURE__*/React.createElement("li", null, /*#__PURE__*/React.createElement(ObjectState, {
          value: member.ObjectState,
          onChange: event => this.updateGroupMember(groupName, i, "ObjectState", event.target.value)
        })))));
      }

      groups.push( /*#__PURE__*/React.createElement("details", {
        key: groupName
      }, /*#__PURE__*/React.createElement("summary", null, groupName), /*#__PURE__*/React.createElement("button", {
        type: "button",
        onClick: this.removeGroup.bind(this, groupName)
      }, "Remove Group"), /*#__PURE__*/React.createElement("br", null), /*#__PURE__*/React.createElement("button", {
        type: "button",
        onClick: event => this.addGroupMember(groupName, event)
      }, "Add Group Member"), /*#__PURE__*/React.createElement("ul", null, groupMembers)));
    }

    return /*#__PURE__*/React.createElement("details", null, /*#__PURE__*/React.createElement("summary", null, "Groups"), /*#__PURE__*/React.createElement("input", {
      ref: this.newGroupName
    }), /*#__PURE__*/React.createElement("button", {
      type: "button",
      onClick: this.addGroup.bind(this)
    }, "Add Group"), /*#__PURE__*/React.createElement("ul", null, groups));
  }

}

class Processes extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      processes: []
    };
    this.addProcess = this.addProcess.bind(this);
    this.removeProcess = this.removeProcess.bind(this);
    this.updateProcess = this.updateProcess.bind(this);
  }

  componentDidMount() {
    this.setProcesses(this.props.processes);
  }

  componentWillReceiveProps(newProps) {
    this.setProcesses(newProps.processes);
  }

  setProcesses(processes) {
    if (processes === undefined || processes === null) {
      processes = [];
    }

    this.setState({
      processes: processes
    });
  }

  addProcess() {
    let empty = {
      CommandLine: "",
      ObjectState: "Keep"
    };
    let processes = [...this.state.processes, empty];
    this.setState({
      processes: processes
    });
    this.props.callback("Processes", processes);
  }

  removeProcess(id) {
    let processes = this.state.processes.filter(function (_, index) {
      return index != id;
    });
    this.setState({
      processes: processes
    });
    this.props.callback("Processes", processes);
  }

  updateProcess(id, field, event) {
    let updated = this.state.processes;
    let value = event.target.value;
    updated[id] = _objectSpread(_objectSpread({}, updated[id]), {}, {
      [field]: value
    });
    this.setState({
      processes: updated
    });
    this.props.callback("Processes", updated);
  }

  render() {
    let processes = [];

    for (let i in this.state.processes) {
      let entry = this.state.processes[i];
      processes.push( /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
        type: "text",
        value: entry.CommandLine,
        onChange: event => this.updateProcess(i, "CommandLine", event)
      })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(ObjectState, {
        value: entry.ObjectState,
        onChange: event => this.updateProcess(i, "ObjectState", event)
      })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("button", {
        type: "button",
        onClick: this.removeProcess.bind(this, i)
      }, "-"))));
    }

    return /*#__PURE__*/React.createElement("details", null, /*#__PURE__*/React.createElement("summary", null, "Processes"), /*#__PURE__*/React.createElement("button", {
      type: "button",
      onClick: this.addProcess.bind(this)
    }, "Add Process"), /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Command Line"), /*#__PURE__*/React.createElement("th", null, "State"), /*#__PURE__*/React.createElement("th", null)), processes));
  }

}

class Software extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      software: []
    };
    this.addSoftware = this.addSoftware.bind(this);
    this.removeSoftware = this.removeSoftware.bind(this);
    this.updateSoftware = this.updateSoftware.bind(this);
  }

  componentDidMount() {
    this.setSoftware(this.props.software);
  }

  componentWillReceiveProps(newProps) {
    this.setSoftware(newProps.software);
  }

  setSoftware(software) {
    if (software === undefined || software === null) {
      software = [];
    }

    this.setState({
      software: software
    });
  }

  addSoftware() {
    let empty = {
      Name: "",
      Version: "",
      ObjectState: "Keep"
    };
    let software = [...this.state.software, empty];
    this.setState({
      software: software
    });
    this.props.callback("Software", software);
  }

  removeSoftware(id) {
    let software = this.state.software.filter(function (_, index) {
      return index != id;
    });
    this.setState({
      software: software
    });
    this.props.callback("Software", software);
  }

  updateSoftware(id, field, event) {
    let updated = this.state.software;
    let value = event.target.value;
    updated[id] = _objectSpread(_objectSpread({}, updated[id]), {}, {
      [field]: value
    });
    this.setState({
      software: updated
    });
    this.props.callback("Software", updated);
  }

  render() {
    let software = [];

    for (let i in this.state.software) {
      let entry = this.state.software[i];
      software.push( /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
        type: "text",
        value: entry.Name,
        onChange: event => this.updateSoftware(i, "Name", event)
      })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
        type: "text",
        value: entry.Version,
        onChange: event => this.updateSoftware(i, "Version", event)
      })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(ObjectState, {
        value: entry.ObjectState,
        onChange: event => this.updateSoftware(i, "ObjectState", event)
      })), /*#__PURE__*/React.createElement("button", {
        type: "button",
        onClick: this.removeSoftware.bind(this, i)
      }, "-")));
    }

    return /*#__PURE__*/React.createElement("details", null, /*#__PURE__*/React.createElement("summary", null, "Software"), /*#__PURE__*/React.createElement("button", {
      type: "button",
      onClick: this.addSoftware.bind(this)
    }, "Add Software"), /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Name"), /*#__PURE__*/React.createElement("th", null, "Version"), /*#__PURE__*/React.createElement("th", null, "State"), /*#__PURE__*/React.createElement("th", null)), software));
  }

}

class NetworkConnections extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      conns: []
    };
    this.add = this.add.bind(this);
    this.remove = this.remove.bind(this);
    this.update = this.update.bind(this);
  }

  componentDidMount() {
    this.setConns(this.props.conns);
  }

  componentWillReceiveProps(newProps) {
    this.setConns(newProps.conns);
  }

  setConns(conns) {
    if (conns === undefined || conns === null) {
      conns = [];
    }

    this.setState({
      conns: conns
    });
  }

  add() {
    let empty = {
      Protocol: "",
      LocalAddress: "",
      LocalPort: "",
      RemoteAddress: "",
      RemotePort: "",
      ObjectState: "Keep"
    };
    let conns = [...this.state.conns, empty];
    this.setState({
      conns: conns
    });
    this.props.callback("NetworkConnections", conns);
  }

  remove(id) {
    let conns = this.state.conns.filter(function (_, index) {
      return index != id;
    });
    this.setState({
      conns: conns
    });
    this.props.callback("NetworkConnections", conns);
  }

  update(id, field, event) {
    let updated = this.state.conns;
    let value = event.target.value;
    updated[id] = _objectSpread(_objectSpread({}, updated[id]), {}, {
      [field]: value
    });
    this.setState({
      conns: updated
    });
    this.props.callback("NetworkConnections", updated);
  }

  render() {
    let conns = [];

    for (let i in this.state.conns) {
      let entry = this.state.conns[i];
      conns.push( /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("select", {
        value: entry.Protocol,
        onChange: event => this.update(i, "Protocol", event)
      }, /*#__PURE__*/React.createElement("option", {
        value: ""
      }), /*#__PURE__*/React.createElement("option", {
        value: "TCP"
      }, "TCP"), /*#__PURE__*/React.createElement("option", {
        value: "UDP"
      }, "UDP"))), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
        type: "text",
        value: entry.LocalAddress,
        onChange: event => this.update(i, "LocalAddress", event)
      })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
        type: "text",
        value: entry.LocalPort,
        onChange: event => this.update(i, "LocalPort", event)
      })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
        type: "text",
        value: entry.RemoteAddress,
        onChange: event => this.update(i, "RemoteAddress", event)
      })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
        type: "text",
        value: entry.RemotePort,
        onChange: event => this.update(i, "RemotePort", event)
      })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(ObjectState, {
        value: entry.ObjectState,
        onChange: event => this.update(i, "ObjectState", event)
      })), /*#__PURE__*/React.createElement("button", {
        type: "button",
        onClick: this.remove.bind(this, i)
      }, "-")));
    }

    return /*#__PURE__*/React.createElement("details", null, /*#__PURE__*/React.createElement("summary", null, "Network Connections"), /*#__PURE__*/React.createElement("button", {
      type: "button",
      onClick: this.add.bind(this)
    }, "Add Network Connection"), /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Protocol"), /*#__PURE__*/React.createElement("th", null, "Local Address"), /*#__PURE__*/React.createElement("th", null, "Local Port"), /*#__PURE__*/React.createElement("th", null, "Remote Address"), /*#__PURE__*/React.createElement("th", null, "Remote Port"), /*#__PURE__*/React.createElement("th", null, "State"), /*#__PURE__*/React.createElement("th", null)), conns));
  }

}

class ScheduledTasks extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      tasks: []
    };
    this.add = this.add.bind(this);
    this.remove = this.remove.bind(this);
    this.update = this.update.bind(this);
  }

  componentDidMount() {
    this.setTasks(this.props.tasks);
  }

  componentWillReceiveProps(newProps) {
    this.setTasks(newProps.tasks);
  }

  setTasks(tasks) {
    if (tasks === undefined || tasks === null) {
      tasks = [];
    }

    this.setState({
      tasks: tasks
    });
  }

  add() {
    let empty = {
      Name: "",
      Path: "",
      Enabled: true
    };
    let tasks = [...this.state.tasks, empty];
    this.setState({
      tasks: tasks
    });
    this.props.callback("ScheduledTasks", tasks);
  }

  remove(id) {
    let tasks = this.state.tasks.filter(function (_, index) {
      return index != id;
    });
    this.setState({
      tasks: tasks
    });
    this.props.callback("ScheduledTasks", tasks);
  }

  update(id, field, event) {
    let updated = this.state.tasks;
    let value = event.target.value;

    if (event.target.type === "checkbox") {
      if (event.target.checked) {
        value = true;
      } else {
        value = false;
      }
    }

    updated[id] = _objectSpread(_objectSpread({}, updated[id]), {}, {
      [field]: value
    });
    this.setState({
      tasks: updated
    });
    this.props.callback("ScheduledTasks", updated);
  }

  render() {
    let tasks = [];

    for (let i in this.state.tasks) {
      let entry = this.state.tasks[i];
      let enabledStr = "Enabled";

      if (!entry.Enabled) {
        enabledStr = "Disabled";
      }

      tasks.push( /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
        type: "text",
        value: entry.Name,
        onChange: event => this.update(i, "Name", event)
      })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
        type: "text",
        value: entry.Path,
        onChange: event => this.update(i, "Path", event)
      })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
        type: "checkbox",
        checked: entry.Enabled,
        onChange: event => this.update(i, "Enabled", event)
      })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(ObjectState, {
        value: entry.ObjectState,
        onChange: event => this.update(i, "ObjectState", event)
      })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("button", {
        type: "button",
        onClick: this.remove.bind(this, i)
      }, "-"))));
    }

    return /*#__PURE__*/React.createElement("details", null, /*#__PURE__*/React.createElement("summary", null, "Scheduled Tasks"), /*#__PURE__*/React.createElement("button", {
      type: "button",
      onClick: this.add.bind(this)
    }, "Add Scheduled Task"), /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Name"), /*#__PURE__*/React.createElement("th", null, "Path"), /*#__PURE__*/React.createElement("th", null, "Enabled"), /*#__PURE__*/React.createElement("th", null, "State"), /*#__PURE__*/React.createElement("th", null)), tasks));
  }

}

class WindowsFirewallProfiles extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      profiles: []
    };
    this.add = this.add.bind(this);
    this.remove = this.remove.bind(this);
    this.update = this.update.bind(this);
  }

  componentDidMount() {
    this.setProfiles(this.props.profiles);
  }

  componentWillReceiveProps(newProps) {
    this.setProfiles(newProps.profiles);
  }

  setProfiles(profiles) {
    if (profiles === undefined || profiles === null) {
      profiles = [];
    }

    this.setState({
      profiles: profiles
    });
  }

  add() {
    let empty = {
      Name: "",
      Enabled: true,
      DefaultInboundAction: "Block",
      DefaultOutboundAction: "Allow"
    };
    let profiles = [...this.state.profiles, empty];
    this.setState({
      profiles: profiles
    });
    this.props.callback("WindowsFirewallProfiles", profiles);
  }

  remove(id) {
    let profiles = this.state.profiles.filter(function (_, index) {
      return index != id;
    });
    this.setState({
      profiles: profiles
    });
    this.props.callback("WindowsFirewallProfiles", profiles);
  }

  update(id, field, event) {
    let updated = this.state.profiles;
    let value = event.target.value;

    if (event.target.type === "checkbox") {
      if (event.target.checked) {
        value = true;
      } else {
        value = false;
      }
    }

    updated[id] = _objectSpread(_objectSpread({}, updated[id]), {}, {
      [field]: value
    });
    this.setState({
      profiles: updated
    });
    this.props.callback("WindowsFirewallProfiles", updated);
  }

  render() {
    let profiles = [];

    for (let i in this.state.profiles) {
      let entry = this.state.profiles[i];
      let enabledStr = "Enabled";

      if (!entry.Enabled) {
        enabledStr = "Disabled";
      }

      profiles.push( /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("select", {
        value: entry.Name,
        onChange: event => this.update(i, "Name", event)
      }, /*#__PURE__*/React.createElement("option", {
        disabled: true,
        key: "",
        value: ""
      }), /*#__PURE__*/React.createElement("option", null, "Domain"), /*#__PURE__*/React.createElement("option", null, "Public"), /*#__PURE__*/React.createElement("option", null, "Private"))), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
        type: "checkbox",
        checked: entry.Enabled,
        onChange: event => this.update(i, "Enabled", event)
      })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("select", {
        value: entry.DefaultInboundAction,
        onChange: event => this.update(i, "DefaultInboundAction", event)
      }, /*#__PURE__*/React.createElement("option", null, "Block"), /*#__PURE__*/React.createElement("option", null, "Allow"), /*#__PURE__*/React.createElement("option", null, "NotConfigured"))), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("select", {
        value: entry.DefaultOutboundAction,
        onChange: event => this.update(i, "DefaultOutboundAction", event)
      }, /*#__PURE__*/React.createElement("option", null, "Block"), /*#__PURE__*/React.createElement("option", null, "Allow"), /*#__PURE__*/React.createElement("option", null, "NotConfigured"))), /*#__PURE__*/React.createElement("button", {
        type: "button",
        onClick: this.remove.bind(this, i)
      }, "-")));
    }

    return /*#__PURE__*/React.createElement("details", null, /*#__PURE__*/React.createElement("summary", null, "Windows Firewall Profiles"), /*#__PURE__*/React.createElement("button", {
      type: "button",
      onClick: this.add.bind(this)
    }, "Add Windows Firewall profile"), /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Name"), /*#__PURE__*/React.createElement("th", null, "Enabled"), /*#__PURE__*/React.createElement("th", null, "Inbound"), /*#__PURE__*/React.createElement("th", null, "Outbound"), /*#__PURE__*/React.createElement("th", null)), profiles));
  }

}

class WindowsFirewallRules extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      rules: []
    };
    this.add = this.add.bind(this);
    this.remove = this.remove.bind(this);
    this.update = this.update.bind(this);
  }

  componentDidMount() {
    this.setRules(this.props.rules);
  }

  componentWillReceiveProps(newProps) {
    this.setRules(newProps.rules);
  }

  setRules(rules) {
    if (rules === undefined || rules === null) {
      rules = [];
    }

    this.setState({
      rules: rules
    });
  }

  add() {
    let empty = {
      DisplayName: "",
      Enabled: true,
      Direction: "",
      Action: "",
      ObjectState: "Keep"
    };
    let rules = [...this.state.rules, empty];
    this.setState({
      rules: rules
    });
    this.props.callback("WindowsFirewallRules", rules);
  }

  remove(id) {
    let rules = this.state.rules.filter(function (_, index) {
      return index != id;
    });
    this.setState({
      rules: rules
    });
    this.props.callback("WindowsFirewallRules", rules);
  }

  update(id, field, event) {
    let updated = this.state.rules;
    let value = event.target.value;

    if (event.target.type === "checkbox") {
      if (event.target.checked) {
        value = true;
      } else {
        value = false;
      }
    }

    updated[id] = _objectSpread(_objectSpread({}, updated[id]), {}, {
      [field]: value
    });
    this.setState({
      rules: updated
    });
    this.props.callback("WindowsFirewallRules", updated);
  }

  render() {
    let rules = [];

    for (let i in this.state.rules) {
      let entry = this.state.rules[i];
      let enabledStr = "Enabled";

      if (!entry.Enabled) {
        enabledStr = "Disabled";
      }

      let ruleOptions = null;

      if (entry.ObjectState != "Remove") {
        ruleOptions = /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
          type: "checkbox",
          checked: entry.Enabled,
          onChange: event => this.update(i, "Enabled", event)
        })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("select", {
          value: entry.Protocol,
          onChange: event => this.update(i, "Protocol", event)
        }, /*#__PURE__*/React.createElement("option", {
          value: ""
        }), /*#__PURE__*/React.createElement("option", {
          value: "TCP"
        }, "TCP"), /*#__PURE__*/React.createElement("option", {
          value: "UDP"
        }, "UDP"))), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
          type: "text",
          value: entry.LocalPort,
          onChange: event => this.update(i, "LocalPort", event)
        })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
          type: "text",
          value: entry.RemoteAddress,
          onChange: event => this.update(i, "RemoteAddress", event)
        })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
          type: "text",
          value: entry.RemotePort,
          onChange: event => this.update(i, "RemotePort", event)
        })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("select", {
          value: entry.Direction,
          onChange: event => this.update(i, "Direction", event)
        }, /*#__PURE__*/React.createElement("option", {
          disabled: true,
          key: "",
          value: ""
        }), /*#__PURE__*/React.createElement("option", null, "Inbound"), /*#__PURE__*/React.createElement("option", null, "Outbound"))), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("select", {
          value: entry.Action,
          onChange: event => this.update(i, "Action", event)
        }, /*#__PURE__*/React.createElement("option", {
          disabled: true,
          key: "",
          value: ""
        }), /*#__PURE__*/React.createElement("option", null, "Block"), /*#__PURE__*/React.createElement("option", null, "Allow"))));
      }

      rules.push( /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
        type: "text",
        value: entry.DisplayName,
        onChange: event => this.update(i, "DisplayName", event)
      })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(ObjectState, {
        value: entry.ObjectState,
        onChange: event => this.update(i, "ObjectState", event)
      })), ruleOptions, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("button", {
        type: "button",
        onClick: this.remove.bind(this, i)
      }, "-"))));
    }

    return /*#__PURE__*/React.createElement("details", null, /*#__PURE__*/React.createElement("summary", null, "Windows Firewall rules"), /*#__PURE__*/React.createElement("button", {
      type: "button",
      onClick: this.add.bind(this)
    }, "Add Windows Firewall rule"), /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Display name"), /*#__PURE__*/React.createElement("th", null, "State"), /*#__PURE__*/React.createElement("th", null, "Enabled"), /*#__PURE__*/React.createElement("th", null, "Protocol"), /*#__PURE__*/React.createElement("th", null, "Local Address"), /*#__PURE__*/React.createElement("th", null, "Local Port"), /*#__PURE__*/React.createElement("th", null, "Remote Address"), /*#__PURE__*/React.createElement("th", null, "Remote Port"), /*#__PURE__*/React.createElement("th", null, "Direction"), /*#__PURE__*/React.createElement("th", null, "Action"), /*#__PURE__*/React.createElement("th", null)), rules));
  }

}

class WindowsSettings extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      settings: []
    };
    this.add = this.add.bind(this);
    this.remove = this.remove.bind(this);
    this.update = this.update.bind(this);
  }

  componentDidMount() {
    this.setSettings(this.props.settings);
  }

  componentWillReceiveProps(newProps) {
    this.setSettings(newProps.settings);
  }

  setSettings(settings) {
    if (settings === undefined || settings === null) {
      settings = [];
    }

    this.setState({
      settings: settings
    });
  }

  add() {
    let empty = {
      Key: "",
      Value: ""
    };
    let settings = [...this.state.settings, empty];
    this.setState({
      settings: settings
    });
    this.props.callback("WindowsSettings", settings);
  }

  remove(id) {
    let settings = this.state.settings.filter(function (_, index) {
      return index != id;
    });
    this.setState({
      settings: settings
    });
    this.props.callback("WindowsSettings", settings);
  }

  update(id, field, event) {
    let updated = this.state.settings;
    let value = event.target.value;

    if (event.target.type === "checkbox") {
      if (event.target.checked) {
        value = true;
      } else {
        value = false;
      }
    }

    updated[id] = _objectSpread(_objectSpread({}, updated[id]), {}, {
      [field]: value
    });
    this.setState({
      settings: updated
    });
    this.props.callback("WindowsSettings", settings);
  }

  render() {
    let settings = [];

    for (let i in this.state.settings) {
      let entry = this.state.settings[i];
      settings.push( /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
        type: "text",
        value: entry.Key,
        onChange: event => this.update(i, "Key", event)
      })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("input", {
        type: "text",
        value: entry.Value,
        onChange: event => this.update(i, "Value", event)
      })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("button", {
        type: "button",
        onClick: this.remove.bind(this, i)
      }, "-"))));
    }

    return /*#__PURE__*/React.createElement("details", null, /*#__PURE__*/React.createElement("summary", null, "Windows Settings"), /*#__PURE__*/React.createElement("button", {
      type: "button",
      onClick: this.add.bind(this)
    }, "Add Windows setting"), /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Key"), /*#__PURE__*/React.createElement("th", null, "Value"), /*#__PURE__*/React.createElement("th", null)), settings));
  }

}

class Item extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("label", {
      htmlFor: this.props.name
    }, this.props.name), /*#__PURE__*/React.createElement("input", {
      name: this.props.name,
      type: this.props.type,
      value: this.props.value,
      checked: this.props.checked,
      disabled: this.props.disabled
    }));
  }

}

class ItemMap extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      item: "",
      value: this.props.value,
      mapItems: [],
      listItems: []
    };
    this.add = this.add.bind(this);
    this.remove = this.remove.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.handleCallback = this.handleCallback.bind(this);
  }

  componentWillReceiveProps(newProps) {
    if (this.props.value != newProps.value) {
      this.setState({
        value: newProps.value
      });
    }
  }

  handleChange(event) {
    let value = Number(event.target.value);
    this.setState({
      item: value
    });
  }

  handleCallback(key, value) {
    let v = _objectSpread(_objectSpread({}, this.state.value), {}, {
      [key]: value
    });

    this.setState({
      value: v
    });
    this.props.callback(this.props.name, v);
  }

  add() {
    if (!this.state.item) {
      return;
    }

    if (this.state.value && this.state.value[this.state.item] != null) {
      return;
    }

    let value = _objectSpread(_objectSpread({}, this.state.value), {}, {
      [this.state.item]: []
    });

    this.setState({
      item: "",
      value: value
    });
    this.props.callback(this.props.name, value);
  }

  remove(id) {
    if (this.state.value == null) {
      return;
    }

    let value = _objectSpread(_objectSpread({}, this.state.value), {}, {
      [id]: undefined
    });

    this.setState({
      value: value
    });
    this.props.callback(this.props.name, value);
  }

  componentWillMount() {
    this.props.mapItems(items => {
      this.setState({
        mapItems: items
      });
    });
    this.props.listItems(items => {
      this.setState({
        listItems: items
      });
    });
  }

  render() {
    let rows = [];

    if (this.state.value) {
      for (let i in this.state.value) {
        if (this.state.value[i] === undefined) {
          continue;
        }

        let text = i;
        let matches = this.state.mapItems.filter(obj => {
          return obj.ID == i;
        });

        if (matches.length > 0) {
          text = matches[0].Display;
        }

        rows.push( /*#__PURE__*/React.createElement("details", {
          key: i
        }, /*#__PURE__*/React.createElement("summary", null, text), /*#__PURE__*/React.createElement("button", {
          type: "button",
          onClick: this.remove.bind(this, i)
        }, "-"), /*#__PURE__*/React.createElement("ul", null, /*#__PURE__*/React.createElement(ItemList, {
          name: i,
          label: this.props.listLabel,
          type: "select",
          listItems: this.state.listItems,
          value: this.state.value[i],
          callback: this.handleCallback
        }))));
      }
    }

    let optionsMap = []; // empty selection

    optionsMap.push( /*#__PURE__*/React.createElement("option", {
      disabled: true,
      key: "",
      value: ""
    }));

    for (let i in this.state.mapItems) {
      let option = this.state.mapItems[i]; // skip already selected

      if (this.state.value && this.state.value[option.ID] != null) {
        continue;
      }

      optionsMap.push( /*#__PURE__*/React.createElement("option", {
        key: option.ID,
        value: option.ID
      }, option.Display));
    }

    return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("label", null, this.props.label), /*#__PURE__*/React.createElement("ul", null, rows, /*#__PURE__*/React.createElement("select", {
      value: this.state.item,
      onChange: this.handleChange
    }, optionsMap), /*#__PURE__*/React.createElement("button", {
      type: "button",
      onClick: this.add
    }, "+")));
  }

}

class ItemList extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      item: "",
      value: this.props.value
    };
    this.add = this.add.bind(this);
    this.remove = this.remove.bind(this);
    this.handleChange = this.handleChange.bind(this);
  }

  componentWillReceiveProps(newProps) {
    if (this.props.value != newProps.value) {
      this.setState({
        value: newProps.value
      });
    }
  }

  handleChange(event) {
    let value = event.target.value;

    if (this.props.type === "select") {
      value = Number(value);
    }

    this.setState({
      item: value
    });
  }

  add() {
    if (!this.state.item) {
      return;
    }

    if (this.state.value && this.state.value.includes(this.state.item)) {
      return;
    }

    let value = null;

    if (this.state.value == null) {
      value = [this.state.item];
    } else {
      value = [...this.state.value, this.state.item];
    }

    this.setState({
      item: "",
      value: value
    });
    this.props.callback(this.props.name, value);
  }

  remove(id) {
    if (this.state.value == null) {
      return;
    }

    let value = this.state.value.filter(function (_, index) {
      return index != id;
    });
    this.setState({
      value: value
    });
    this.props.callback(this.props.name, value);
  }

  render() {
    let rows = [];

    if (this.state.value) {
      for (let i in this.state.value) {
        let text = this.state.value[i];

        if (this.props.type === "select") {
          let matches = this.props.listItems.filter(obj => {
            return obj.ID == text;
          });

          if (matches.length > 0) {
            text = matches[0].Display;
          }
        }

        rows.push( /*#__PURE__*/React.createElement("li", {
          key: i
        }, text, /*#__PURE__*/React.createElement("button", {
          type: "button",
          onClick: this.remove.bind(this, i)
        }, "-")));
      }
    }

    let input = /*#__PURE__*/React.createElement("input", {
      type: this.props.type,
      value: this.state.item,
      onChange: this.handleChange
    });

    if (this.props.type === "select") {
      let optionsList = []; // empty selection

      optionsList.push( /*#__PURE__*/React.createElement("option", {
        disabled: true,
        key: "",
        value: ""
      }));

      for (let i in this.props.listItems) {
        let option = this.props.listItems[i]; // skip already selected

        if (this.state.value && this.state.value.indexOf(option.ID) != -1) {
          continue;
        }

        optionsList.push( /*#__PURE__*/React.createElement("option", {
          key: option.ID,
          value: option.ID
        }, option.Display));
      }

      input = /*#__PURE__*/React.createElement("select", {
        value: this.state.item,
        onChange: this.handleChange
      }, optionsList);
    }

    return /*#__PURE__*/React.createElement("details", null, /*#__PURE__*/React.createElement("summary", null, this.props.label), /*#__PURE__*/React.createElement("ul", null, rows, input, /*#__PURE__*/React.createElement("button", {
      type: "button",
      onClick: this.add
    }, "+")));
  }

}

ReactDOM.render( /*#__PURE__*/React.createElement(App, null), document.getElementById('app'));