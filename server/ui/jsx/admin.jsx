'use strict';

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
    }

    this.authCallback = this.authCallback.bind(this);
    this.setPage = this.setPage.bind(this);
    this.logout = this.logout.bind(this);
  }

  authCallback(statusCode) {
    if (statusCode == 200) {
      this.setState({
        authenticated: true
      });
    }
    else {
      this.setState({
        authenticated: false
      })
    }    
  }

  logout() {
    let url = "/logout"
    fetch(url, {
      credentials: 'same-origin',
      method: "DELETE"
    })
    .then(function(_) {
      this.setState({
        authenticated: false
      })
    }.bind(this));
  }

  componentDidMount() {
    // check if logged in by visiting the following URL
    let url = "/templates";
    fetch(url, {
      credentials: 'same-origin'
    })
    .then(function(response) {
      this.authCallback(response.status);
    }.bind(this));

    // track which page to be on
    let i = window.location.hash.indexOf('#');
    let hash = window.location.hash.slice(i + 1);
    this.setPage(hash);
    // handle browser back/forward
    window.onhashchange = (e) => {
        let i = e.newURL.indexOf('#');
        let hash = e.newURL.slice(i + 1);
        this.setPage(hash);
    };
  }

  setPage(hash) {
    let parts = hash.split("/")
    let page = parts[0];
    let id = null;
    if (parts.length >= 2) {
      id = parts[1];
    }
    this.setState({
      page: page,
      id: id
    })
  }

  updateTeamCallback() {
    this.setState({
      lastUpdatedTeams: Date.now()
    })
  }

  updateHostCallback() {
    this.setState({
      lastUpdatedHosts: Date.now()
    })
  }

  updateTemplateCallback() {
    this.setState({
      lastUpdatedTemplates: Date.now()
    })
  }

  updateScenarioCallback() {
    this.setState({
      lastUpdatedScenarios: Date.now()
    })
  }

  updateAdministratorCallback() {
    this.setState({
      lastUpdatedAdministrators: Date.now()
    })
  }

  render() {
    if (!this.state.authenticated) {
      return (
        <div className="App">
          <Login callback={this.authCallback}/>
        </div>
      );
    }

    // reset links to available
    let classes_teams = ["nav-button"];
    let classes_hosts = ["nav-button"];
    let classes_templates = ["nav-button"];
    let classes_scenarios = ["nav-button"];
    let classes_administrators = ["nav-button"];

    // default page is empty
    let page = (<React.Fragment></React.Fragment>)
    let content = (<React.Fragment></React.Fragment>)
    if (this.state.page == "teams") {
      classes_teams.push("nav-button-selected");
      page = (<Teams lastUpdated={this.state.lastUpdatedTeams} selected={this.state.id}/>);
      content = (<TeamEntry id={this.state.id} updateCallback={this.updateTeamCallback.bind(this)}/>);
    }
    else if (this.state.page == "hosts") {
      classes_hosts.push("nav-button-selected");
      page = (<Hosts lastUpdated={this.state.lastUpdatedHosts} selected={this.state.id}/>);
      content = (<HostEntry id={this.state.id} updateCallback={this.updateHostCallback.bind(this)}/>);
    }
    else if (this.state.page == "templates") {
      classes_templates.push("nav-button-selected");
      page = (<Templates lastUpdated={this.state.lastUpdatedTemplates} selected={this.state.id}/>);
      content = (<TemplateEntry id={this.state.id} updateCallback={this.updateTemplateCallback.bind(this)}/>);
    }
    else if (this.state.page == "scenarios") {
      classes_scenarios.push("nav-button-selected");
      page = (<Scenarios lastUpdated={this.state.lastUpdatedScenarios} selected={this.state.id}/>);
      content = (<ScenarioEntry id={this.state.id} updateCallback={this.updateScenarioCallback.bind(this)}/>);
    }
    else if (this.state.page == "administrators") {
      classes_administrators.push("nav-button-selected");
      page = (<Administrators lastUpdated={this.state.lastUpdatedAdministrators} selected={this.state.id}/>);
      content = (<AdministratorEntry username={this.state.id} updateCallback={this.updateAdministratorCallback.bind(this)}/>);
    }

    return (
      <div className="App">
        <div className="heading">
          <h1>cp-scoring</h1>
        </div>
        <div className="navbar">
          <a className={classes_teams.join(" ")} href="#teams">Teams</a>
          <a className={classes_hosts.join(" ")} href="#hosts">Hosts</a>
          <a className={classes_templates.join(" ")} href="#templates">Templates</a>
          <a className={classes_scenarios.join(" ")} href="#scenarios">Scenarios</a>
          <a className={classes_administrators.join(" ")} href="#administrators">Administrators</a>
          <div className="right">
            <button onClick={this.logout}>Logout</button>
          </div>
        </div>
        <div className="toc">
          {page}
        </div>
        <div className="content">
          {content}
        </div>
      </div>
    );
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
    }

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
    })
    .then(async function(response) {
      if (response.status === 200) {
        let data = await response.json();
        return {
          error: null,
          items: data
        }
      }
      let text = await response.text();
      return {
        error: text
      }
    })
    .then(function(s) {
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
      rows.push(
        <li key={i}>
          <a className={classes.join(" ")} href={"#administrators/" + administrator}>{administrator}</a>
        </li>
      );
    }

    return (
      <div className="Admins">
        <strong>Administrators</strong>
        <Error message={this.state.error} />
        <ul>{rows}</ul>
      </div>
    );
  }
}

class AdministratorEntry extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      user: {}
    }
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
      user: {
        ...this.state.user,
        [event.target.name]: value
      }
    })
  }

  saveAdministrator(event) {
    event.preventDefault();

    var url = "/admins";
    // for updating existing user
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
    })
    .then(async function(response) {
      if (response.status === 200) {
        this.props.updateCallback();
        let username = this.state.user.Username;
        if (username != undefined && username != null) {
          window.location.href = "#administrators/" + username;
        }
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
    }.bind(this));;
  }

  deleteAdministrator(username) {
    var url = "/admins/" + username;

    fetch(url, {
      credentials: 'same-origin',
      method: 'DELETE'
    })
    .then(async function(response) {
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
      }
    }.bind(this))
    .then(function(s) {
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
      content = (
        <form onChange={this.updateAdministrator.bind(this)} onSubmit={this.saveAdministrator.bind(this)}>
          <Item name="Username" value={this.state.user.Username || ""} disabled={existingUser}/>
          <Item name="Password" type="password" value={this.state.user.Password} />
          <p />
          <button type="submit">Submit</button>
          <button class="right" type="button" disabled={!this.state.user.Username} onClick={this.deleteAdministrator.bind(this, this.state.user.Username)}>Delete</button>
        </form>
      );
    }
    
    return (
      <React.Fragment>
        <button type="button" onClick={this.newAdministrator.bind(this)}>New Administrator</button>
        <hr />
        <Error message={this.state.error} />
        {content}
      </React.Fragment>
    );
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
      rows.push(
        <li key={team.ID}>
          <a className={classes.join(" ")} href={"#teams/" + team.ID}>[{team.ID}] {team.Name}</a>
        </li>
      );
    }

    return (
      <div className="Teams">
        <strong>Teams</strong>
        <Error message={this.state.error} />
        <ul>{rows}</ul>
      </div>
    );
  }
}

class TeamEntry extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      team: {}
    }

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
    })
    .then(async function(response) {
      if (response.status === 200) {
        let data = await response.json();
        return {
          error: null,
          team: data
        }
      }
      let text = await response.text();
      return {
        error: text
      }
    })
    .then(function(s) {
      this.setState(s);
    }.bind(this));
  }

  updateTeam(event) {
    let value = event.target.value;
    if (event.target.type == 'checkbox') {
      value = event.target.checked;
    }
    this.setState({
      team: {
        ...this.state.team,
        [event.target.name]: value
      }
    })
  }

  deleteTeam(id) {
    var url = "/teams/" + id;

    fetch(url, {
      credentials: 'same-origin',
      method: 'DELETE'
    })
    .then(async function(response) {
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
      }
    }.bind(this))
    .then(function(s) {
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
    })
    .then(async function(response) {
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
    }.bind(this))
    .then(function(s) {
      this.setState(s);
    }.bind(this));
  }

  regenKey() {
    let key = TeamEntry.newKey()
    this.setState({
      team: {
        ...this.state.team,
        Key: key
      }
    })
  }

  render() {
    let content = null;
    if (Object.entries(this.state.team).length != 0) {
      content = (
        <form onChange={this.updateTeam.bind(this)} onSubmit={this.saveTeam.bind(this)}>
          <label htmlFor="ID">ID</label>
          <input disabled value={this.state.team.ID || ""}/>
          <Item name="Name" value={this.state.team.Name}/>
          <Item name="POC" value={this.state.team.POC}/>
          <Item name="Email" type="email" value={this.state.team.Email}/>
          <Item name="Enabled" type="checkbox" checked={!!this.state.team.Enabled}/>
          <br />
          <details>
            <summary>Key</summary>
            <ul>
              <li>
                {this.state.team.Key}
                <br />
                <button type="button" onClick={this.regenKey.bind(this)}>Regenerate</button>
              </li>
            </ul>
          </details>
          <br />
          <div>
            <button type="submit">Save</button>
            <button class="right" type="button" disabled={!this.state.team.ID} onClick={this.deleteTeam.bind(this, this.state.team.ID)}>Delete</button>
          </div>
        </form>
      );
    }

    return (
      <React.Fragment>
        <button type="button" onClick={this.newTeam.bind(this)}>New Team</button>
        <hr />
        <Error message={this.state.error} />
        {content}
      </React.Fragment>
    );
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
      rows.push(
        <li key={scenario.ID}>
          <a className={classes.join(" ")} href={"#scenarios/" + scenario.ID}>{scenario.Name}</a>
        </li>
      );
    }

    return (
      <div className="Scenarios">
        <strong>Scenarios</strong>
        <Error message={this.state.error} />
        <ul>{rows}</ul>
      </div>
    );
  }
}

class ScenarioEntry extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      scenario: {}
    }
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
    })
    .then(async function(response) {
      if (response.status === 200) {
        let data = await response.json();
        return {
          error: null,
          scenario: data
        }
      }
      let text = await response.text();
      return {
        error: text
      }
    })
    .then(function(s) {
      this.setState(s);
    }.bind(this));
  }

  updateScenario(event) {
    let value = event.target.value;
    if (event.target.type == 'checkbox') {
      value = event.target.checked;
    }
    this.setState({
      scenario: {
        ...this.state.scenario,
        [event.target.name]: value
      }
    })
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
    })
    .then(async function(response) {
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
    }.bind(this))
    .then(function(s) {
      this.setState(s);
    }.bind(this));
  }

  deleteScenario(id) {
    var url = "/scenarios/" + id;

    fetch(url, {
      credentials: 'same-origin',
      method: 'DELETE'
    })
    .then(async function(response) {
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
      }
    }.bind(this))
    .then(function(s) {
      this.setState(s);
    }.bind(this));
  }

  handleCallback(key, value) {
    this.setState({
      scenario: {
        ...this.state.scenario,
        [key]: value
      }
    })
  }

  mapItems(callback) {
    var url = "/hosts";

    fetch(url, {
      credentials: 'same-origin'
    })
    .then(async function(response) {
      if (response.status === 200) {
        let data = await response.json();
        let items = data.map(function(host) {
          return {
            ID: host.ID,
            Display: host.Hostname
          }
        });
        callback(items);
        return;
      }
      let text = await response.text();
      this.setState({
        error: text
      });
    }.bind(this));
  };

  listItems(callback) {
    var url = "/templates";

    fetch(url, {
      credentials: 'same-origin'
    })
    .then(async function(response) {
      if (response.status === 200) {
        let data = await response.json();
        let items = data.map(function(template) {
          return {
            ID: template.ID,
            Display: template.Name
          }
        });
        callback(items);
        return;
      }
      let text = await response.text();
      this.setState({
        error: text
      });
    }.bind(this));
  };

  render() {
    let content = null;
    if (Object.entries(this.state.scenario).length != 0) {
      content = (
        <form onChange={this.updateScenario.bind(this)} onSubmit={this.saveScenario.bind(this)}>
          <label htmlFor="ID">ID</label>
          <input disabled value={this.state.scenario.ID || ""}/>
          <Item name="Name" value={this.state.scenario.Name}/>
          <label htmlFor="Description">Description</label>
          <textarea name="Description" rows="10" cols="80" value={this.state.scenario.Description}></textarea>
          <Item name="Enabled" type="checkbox" checked={!!this.state.scenario.Enabled}/>
          <ItemMap name="HostTemplates" label="Hosts" listLabel="Templates" value={this.state.scenario.HostTemplates} callback={this.handleCallback.bind(this)} mapItems={this.mapItems} listItems={this.listItems}/>
          <br />
          <div>
            <button type="submit">Save</button>
            <button class="right" type="button" disabled={!this.state.scenario.ID} onClick={this.deleteScenario.bind(this, this.state.scenario.ID)}>Delete</button>
          </div>
        </form>
      );
    }

    return (
      <React.Fragment>
        <button onClick={this.newScenario.bind(this)}>New Scenario</button>
        <button onClick={this.newScenarioFromExisting.bind(this)}>Clone</button>
        <hr />
        <Error message={this.state.error} />
        {content}
      </React.Fragment>
    )
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
      rows.push(
        <li key={host.ID}>
          <a className={classes.join(" ")} href={"#hosts/" + host.ID}>{host.Hostname} - {host.OS}</a>
        </li>
      );
    }
  
    return (
      <div className="Hosts">
        <strong>Hosts</strong>
        <Error message={this.state.error} />
        <ul>{rows}</ul>
      </div>
    );
  }
}

class HostEntry extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      host: {}
    }
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
    })
    .then(async function(response) {
      if (response.status === 200) {
        let data = await response.json();
        return {
          error: null,
          host: data
        }
      }
      let text = await response.text();
      return {
        error: text
      }
    })
    .then(function(s) {
      this.setState(s);
    }.bind(this));
  }

  updateHost(event) {
    let value = event.target.value;
    if (event.target.type == 'checkbox') {
      value = event.target.checked;
    }
    this.setState({
      host: {
        ...this.state.host,
        [event.target.name]: value
      }
    })
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
    })
    .then(async function(response) {
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
    }.bind(this))
    .then(function(s) {
      this.setState(s);
    }.bind(this));
  }

  deleteHost(id) {
    var url = "/hosts/" + id;

    fetch(url, {
      credentials: 'same-origin',
      method: 'DELETE'
    })
    .then(async function(response) {
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
      }
    }.bind(this))
    .then(function(s) {
      this.setState(s);
    }.bind(this));
  }

  render() {
    let content = null;
    if (Object.entries(this.state.host).length != 0) {
      content = (
        <form onChange={this.updateHost.bind(this)} onSubmit={this.saveHost.bind(this)}>
          <label htmlFor="ID">ID</label>
          <input disabled value={this.state.host.ID || ""}/>
          <Item name="Hostname" type="text" value={this.state.host.Hostname}/>
          <Item name="OS" type="text" value={this.state.host.OS}/>
          <br />
          <div>
            <button type="submit">Save</button>
            <button class="right" type="button" disabled={!this.state.host.ID} onClick={this.deleteHost.bind(this, this.state.host.ID)}>Delete</button>
          </div>
        </form>
      );
    }

    return (
      <React.Fragment>
        <button type="button" onClick={this.newHost.bind(this)}>New Host</button>
        <button type="button" onClick={this.newHostFromExisting.bind(this)}>Clone</button>
        <hr />
        <Error message={this.state.error} />
        {content}
      </React.Fragment>
    );
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
      rows.push(
        <li key={template.ID}>
          <a className={classes.join(" ")} href={"#templates/" + template.ID}>{template.Name}</a>
        </li>
      );
    }

    return (
      <div className="Templates">
        <strong>Templates</strong>
        <Error message={this.state.error} />
        <ul>{rows}</ul>
      </div>
    );
  }
}

class TemplateEntry extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      template: {}
    }
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
    })
    .then(async function(response) {
      if (response.status === 200) {
        let data = await response.json();
        return {
          error: null,
          template: data
        }
      }
      let text = await response.text();
      return {
        error: text
      }
    })
    .then(function(s) {
      this.setState(s);
    }.bind(this));
  }

  updateTemplate(event) {
    let value = event.target.value;
    if (event.target.type == 'checkbox') {
      value = event.target.checked;
    }
    this.setState({
      template: {
        ...this.state.template,
        [event.target.name]: value
      }
    })
  }

  saveRegularTemplate(event) {
    this.saveTemplate(event, "/templates")
  }

  saveStateTemplate(event) {
    this.saveTemplate(event, "/templates/state")
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
    })
    .then(async function(response) {
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
    }.bind(this))
    .then(function(s) {
      this.setState(s);
    }.bind(this));
  }

  deleteTemplate(id) {
    var url = "/templates/" + id;

    fetch(url, {
      credentials: 'same-origin',
      method: 'DELETE'
    })
    .then(async function(response) {
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
      }
    }.bind(this))
    .then(function(s) {
      this.setState(s);
    }.bind(this));
  }

  handleCallback(key, value) {
    let state = {
      ...this.state.template.State,
      [key]: value
    }
    this.setState({
      template: {
        ...this.state.template,
        State: state
      }
    })
  }

  render() {
    let content = null;
    // template from state
    if (this.state.template.StateID != undefined) {
      content = (
        <form onChange={this.updateTemplate.bind(this)} onSubmit={this.saveStateTemplate.bind(this)}>
          <label htmlFor="TemplateName">Template Name</label>
          <input name="TemplateName" value={this.state.template.TemplateName || ""}/>
          <br />
          <label htmlFor="StateID">State ID</label>
          <input name="StateID" value={this.state.template.StateID || ""}/>
          <div>
            <button type="submit">Save</button>
          </div>
        </form>
      )
    }
    // regular template
    else if (Object.entries(this.state.template).length != 0) {
      content = (
        <form onChange={this.updateTemplate.bind(this)} onSubmit={this.saveRegularTemplate.bind(this)}>
          <label htmlFor="ID">ID</label>
          <input disabled value={this.state.template.ID || ""}/>
          <Item name="Name" type="text" value={this.state.template.Name}/>
          <Users users={this.state.template.State.Users} callback={this.handleCallback.bind(this)}/>
          <Groups groups={this.state.template.State.Groups} callback={this.handleCallback.bind(this)}/>
          <Processes processes={this.state.template.State.Processes} callback={this.handleCallback.bind(this)}/>
          <Software software={this.state.template.State.Software} callback={this.handleCallback.bind(this)}/>
          <NetworkConnections conns={this.state.template.State.NetworkConnections} callback={this.handleCallback.bind(this)}/>
          <ScheduledTasks tasks={this.state.template.State.ScheduledTasks} callback={this.handleCallback.bind(this)}/>
          <WindowsFirewallProfiles profiles={this.state.template.State.WindowsFirewallProfiles} callback={this.handleCallback.bind(this)}/>
          <WindowsFirewallRules rules={this.state.template.State.WindowsFirewallRules} callback={this.handleCallback.bind(this)}/>
          <div>
            <button type="submit">Save</button>
            <button class="right" type="button" disabled={!this.state.template.ID} onClick={this.deleteTemplate.bind(this, this.state.template.ID)}>Delete</button>
          </div>
        </form>
      );
    }

    return (
      <React.Fragment>
        <button type="button" onClick={this.newTemplate.bind(this)}>New Template</button>
        <button type="button" onClick={this.newTemplateFromState.bind(this)}>From State</button>
        <button type="button" onClick={this.newTemplateFromExisting.bind(this)}>Clone</button>
        <hr />
        <Error message={this.state.error} />
        {content}
      </React.Fragment>
    );
  }
}

class ObjectState extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <React.Fragment>
        <label>State</label>
        <select value={this.props.value} onChange={this.props.onChange}>
          <option>Add</option>
          <option>Keep</option>
          <option>Remove</option>
        </select>
      </React.Fragment>
    )
  }
}

class Users extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      users: []
    }
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
    let users = [
      ...this.state.users,
      empty
    ];
    this.setState({
      users: users
    });
    this.props.callback("Users", users)
  }

  removeUser(id) {
    let users = this.state.users.filter(function(_, index) {
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
      }
      else {
        value = false;
      }
    }
    else if (event.target.type === "date") {
      let parts = event.target.value.split("-");
      if (parts.length != 3) {
        return;
      }
      let current = new Date(Math.trunc(this.state.users[id].PasswordLastSet * 1000));
      current.setFullYear(parts[0]);
      // months start counting at 0
      current.setMonth(parts[1] - 1);
      current.setDate(parts[2]);
      value = Math.trunc(current.getTime() / 1000);
      if (Number.isNaN(value)) {
        return
      }
    }
    else if (event.target.type === "time") {
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
        return
      }
    }
    updated[id] = {
      ...updated[id],
      [field]: value
    }
    this.setState({
      users: updated
    })
    this.props.callback("Users", updated);
  }

  render() {
    let users = [];
    for (let i = 0; i < this.state.users.length; i++) {
      let user = this.state.users[i];
      let d = new Date(user.PasswordLastSet * 1000)
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
        userOptions = (
          <React.Fragment>
            <li>
              <label>Active</label>
              <input type="checkbox" checked={user.AccountActive} onChange={event=> this.updateUser(i, "AccountActive", event)}/>
            </li>
            <li>
              <label>Password Expires</label>
              <input type="checkbox" checked={user.PasswordExpires} onChange={event=> this.updateUser(i, "PasswordExpires", event)}/>
            </li>
            <li>
              <label>Password Last Set</label>
              <input type="date" value={passwordLastSetDate} onChange={event=> this.updateUser(i, "PasswordLastSet", event)}/>
              <input type="time" value={passwordLastSetTime} onChange={event=> this.updateUser(i, "PasswordLastSet", event)}/>
            </li>
          </React.Fragment>
        );
      }
      users.push(
        <details key={i}>
          <summary>{user.Name}</summary>
          <button type="button" onClick={this.removeUser.bind(this, i)}>-</button>
          <ul>
            <li>
              <label>Name</label>
              <input type="text" value={user.Name} onChange={event=> this.updateUser(i, "Name", event)}/>
            </li>
            <li>
              <ObjectState value={user.ObjectState} onChange={event=> this.updateUser(i, "ObjectState", event)} />
            </li>
            {userOptions}
          </ul>
        </details>
      );
    }

    return (
      <details>
        <summary>Users</summary>
        <button type="button" onClick={this.addUser.bind(this)}>Add User</button>
        <ul>
          {users}
        </ul>
      </details>
    )
  }
}

class Groups extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      groups: []
    }

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
    }
    else if (this.newGroupName.current.value.length === 0) {
      return;
    }
    let groups = {
      ...this.state.groups,
      [this.newGroupName.current.value]: []
    };
    this.setState({
      groups: groups
    });
    this.newGroupName.current.value = "";
    this.props.callback("Groups", groups)
  }

  addGroupMember(groupName) {
    let group = this.state.groups[groupName];
    group.push({
      Name: "",
      ObjectState: "Keep"
    });
    let groups = {
      ...this.state.groups,
      [groupName]: group
    }
    this.setState({
      groups: groups
    })
    this.props.callback("Groups", groups)
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
    let groups = {
      ...this.state.groups,
      [groupName]: group
    }
    this.setState({
      groups: groups
    });
    this.props.callback("Groups", groups);
  }

  updateGroupMember(groupName, memberIndex, key, value) {
    let group = this.state.groups[groupName];
    let member = group[memberIndex];
    member[key] = value
    let groups = {
      ...this.state.groups,
      [groupName]: group
    }
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
        groupMembers.push(
          <details key={i}>
            <summary>{member.Name}</summary>
            <button type="button" onClick={this.removeGroupMember.bind(this, groupName, i)}>-</button>
            <ul>
              <li>
                <label>Name</label>
                <input type="text" value={member.Name} onChange={event=> this.updateGroupMember(groupName, i, "Name", event.target.value)}/>
              </li>
              <li>
                <ObjectState value={member.ObjectState} onChange={event=> this.updateGroupMember(groupName, i, "ObjectState", event.target.value)} />
              </li>
            </ul>
          </details>
        );
      }
      groups.push(
        <details key={groupName}>
          <summary>{groupName}</summary>
          <button type="button" onClick={this.removeGroup.bind(this, groupName)}>Remove Group</button>
          <br />
          <button type="button" onClick={event => this.addGroupMember(groupName, event)}>Add Group Member</button>
          <ul>
            {groupMembers}
          </ul>
        </details>
      );
    }

    return (
      <details>
        <summary>Groups</summary>
        <input ref={this.newGroupName}></input>
        <button type="button" onClick={this.addGroup.bind(this)}>Add Group</button>
        <ul>
          {groups}
        </ul>
      </details>
    )
  }
}

class Processes extends React.Component {
  constructor(props) {
    super(props);
    
    this.state = {
      processes: []
    }

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
    let processes = [
      ...this.state.processes,
      empty
    ];
    this.setState({
      processes: processes
    });
    this.props.callback("Processes", processes)
  }

  removeProcess(id) {
    let processes = this.state.processes.filter(function(_, index) {
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
    updated[id] = {
      ...updated[id],
      [field]: value
    }
    this.setState({
      processes: updated
    })
    this.props.callback("Processes", updated);
  }

  render() {
    let processes = [];
    for (let i in this.state.processes) {
      let entry = this.state.processes[i];
      processes.push(
        <details key={i}>
          <summary>{entry.CommandLine}</summary>
          <button type="button" onClick={this.removeProcess.bind(this, i)}>-</button>
          <ul>
            <li>
              <label>Command line</label>
              <input type="text" value={entry.CommandLine} onChange={event=> this.updateProcess(i, "CommandLine", event)}></input>
            </li>
            <li>
              <ObjectState value={entry.ObjectState} onChange={event=> this.updateProcess(i, "ObjectState", event)} />
            </li>
          </ul>
        </details>
      );
    }

    return (
      <details>
        <summary>Processes</summary>
        <button type="button" onClick={this.addProcess.bind(this)}>Add Process</button>
        <ul>
          {processes}
        </ul>
      </details>
    )
  }
}

class Software extends React.Component {
  constructor(props) {
    super(props);
    
    this.state = {
      software: []
    }

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
    let software = [
      ...this.state.software,
      empty
    ];
    this.setState({
      software: software
    });
    this.props.callback("Software", software)
  }

  removeSoftware(id) {
    let software = this.state.software.filter(function(_, index) {
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
    updated[id] = {
      ...updated[id],
      [field]: value
    }
    this.setState({
      software: updated
    })
    this.props.callback("Software", updated);
  }

  render() {
    let software = [];
    for (let i in this.state.software) {
      let entry = this.state.software[i];
      software.push(
        <details key={i}>
          <summary>{entry.Name}</summary>
          <button type="button" onClick={this.removeSoftware.bind(this, i)}>-</button>
          <ul>
            <li>
              <label>Name</label>
              <input type="text" value={entry.Name} onChange={event=> this.updateSoftware(i, "Name", event)}></input>
            </li>
            <li>
              <label>Version</label>
              <input type="text" value={entry.Version} onChange={event=> this.updateSoftware(i, "Version", event)}></input>
            </li>
            <li>
              <ObjectState value={entry.ObjectState} onChange={event=> this.updateSoftware(i, "ObjectState", event)} />
            </li>
          </ul>
        </details>
      );
    }

    return (
      <details>
        <summary>Software</summary>
        <button type="button" onClick={this.addSoftware.bind(this)}>Add Software</button>
        <ul>
          {software}
        </ul>
      </details>
    )
  }
}

class NetworkConnections extends React.Component {
  constructor(props) {
    super(props);
    
    this.state = {
      conns: []
    }

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
    let conns = [
      ...this.state.conns,
      empty
    ];
    this.setState({
      conns: conns
    });
    this.props.callback("NetworkConnections", conns)
  }

  remove(id) {
    let conns = this.state.conns.filter(function(_, index) {
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
    updated[id] = {
      ...updated[id],
      [field]: value
    }
    this.setState({
      conns: updated
    })
    this.props.callback("NetworkConnections", updated);
  }

  render() {
    let conns = [];
    for (let i in this.state.conns) {
      let entry = this.state.conns[i];
      conns.push(
        <details key={i}>
          <summary>{entry.Protocol} {entry.LocalAddress} {entry.LocalPort} {entry.RemoteAddress} {entry.RemotePort}</summary>
          <button type="button" onClick={this.remove.bind(this, i)}>-</button>
          <ul>
            <li>
              <label>Protocol</label>
              <select value={entry.Protocol} onChange={event=> this.update(i, "Protocol", event)}>
                <option value=""></option>
                <option value="TCP">TCP</option>
                <option value="UDP">UDP</option>
              </select>
            </li>
            <li>
              <label>Local Address</label>
              <input type="text" value={entry.LocalAddress} onChange={event=> this.update(i, "LocalAddress", event)}></input>
            </li>
            <li>
              <label>Local Port</label>
              <input type="text" value={entry.LocalPort} onChange={event=> this.update(i, "LocalPort", event)}></input>
            </li>
            <li>
              <label>Remote Address</label>
              <input type="text" value={entry.RemoteAddress} onChange={event=> this.update(i, "RemoteAddress", event)}></input>
            </li>
            <li>
              <label>Remote Port</label>
              <input type="text" value={entry.RemotePort} onChange={event=> this.update(i, "RemotePort", event)}></input>
            </li>
            <li>
              <ObjectState value={entry.ObjectState} onChange={event=> this.update(i, "ObjectState", event)} />
            </li>
          </ul>
        </details>
      );
    }

    return (
      <details>
        <summary>Network Connections</summary>
        <button type="button" onClick={this.add.bind(this)}>Add Network Connection</button>
        <ul>
          {conns}
        </ul>
      </details>
    )
  }
}

class ScheduledTasks extends React.Component {
  constructor(props) {
    super(props);
    
    this.state = {
      tasks: []
    }

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
    let tasks = [
      ...this.state.tasks,
      empty
    ];
    this.setState({
      tasks: tasks
    });
    this.props.callback("ScheduledTasks", tasks)
  }

  remove(id) {
    let tasks = this.state.tasks.filter(function(_, index) {
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
      }
      else {
        value = false;
      }
    }
    updated[id] = {
      ...updated[id],
      [field]: value
    }
    this.setState({
      tasks: updated
    })
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
      tasks.push(
        <details key={i}>
          <summary>{entry.Name} {entry.Path} {enabledStr}</summary>
          <button type="button" onClick={this.remove.bind(this, i)}>-</button>
          <ul>
            <li>
              <label>Name</label>
              <input type="text" value={entry.Name} onChange={event=> this.update(i, "Name", event)}></input>
            </li>
            <li>
              <label>Path</label>
              <input type="text" value={entry.Path} onChange={event=> this.update(i, "Path", event)}></input>
            </li>
            <li>
              <label>Enabled</label>
              <input type="checkbox" checked={entry.Enabled} onChange={event=> this.update(i, "Enabled", event)}></input>
            </li>
            <li>
              <ObjectState value={entry.ObjectState} onChange={event=> this.update(i, "ObjectState", event)} />
            </li>
          </ul>
        </details>
      );
    }

    return (
      <details>
        <summary>Scheduled Tasks</summary>
        <button type="button" onClick={this.add.bind(this)}>Add Scheduled Task</button>
        <ul>
          {tasks}
        </ul>
      </details>
    )
  }
}

class WindowsFirewallProfiles extends React.Component {
  constructor(props) {
    super(props);
    
    this.state = {
      profiles: []
    }

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
    let profiles = [
      ...this.state.profiles,
      empty
    ];
    this.setState({
      profiles: profiles
    });
    this.props.callback("WindowsFirewallProfiles", profiles)
  }

  remove(id) {
    let profiles = this.state.profiles.filter(function(_, index) {
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
      }
      else {
        value = false;
      }
    }
    updated[id] = {
      ...updated[id],
      [field]: value
    }
    this.setState({
      profiles: updated
    })
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
      profiles.push(
        <details key={i}>
          <summary>Profile: {entry.Name} {enabledStr} Inbound: {entry.DefaultInboundAction} Outbound: {entry.DefaultOutboundAction}</summary>
          <button type="button" onClick={this.remove.bind(this, i)}>-</button>
          <ul>
            <li>
              <label>Name</label>
              <select value={entry.Name} onChange={event=> this.update(i, "Name", event)}>
                <option disabled key="" value="">
                </option>
                <option>Domain</option>
                <option>Public</option>
                <option>Private</option>
              </select>
            </li>
            <li>
              <label>Enabled</label>
              <input type="checkbox" checked={entry.Enabled} onChange={event=> this.update(i, "Enabled", event)}></input>
            </li>
            <li>
              <label>Inbound</label>
              <select value={entry.DefaultInboundAction} onChange={event=> this.update(i, "DefaultInboundAction", event)}>
                <option>Block</option>
                <option>Allow</option>
                <option>NotConfigured</option>
              </select>
            </li>
            <li>
              <label>Outbound</label>
              <select value={entry.DefaultOutboundAction} onChange={event=> this.update(i, "DefaultOutboundAction", event)}>
                <option>Block</option>
                <option>Allow</option>
                <option>NotConfigured</option>
              </select>
            </li>
          </ul>
        </details>
      );
    }

    return (
      <details>
        <summary>Windows Firewall Profiles</summary>
        <button type="button" onClick={this.add.bind(this)}>Add Windows Firewall profile</button>
        <ul>
          {profiles}
        </ul>
      </details>
    )
  }
}

class WindowsFirewallRules extends React.Component {
  constructor(props) {
    super(props);
    
    this.state = {
      rules: []
    }

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
    let rules = [
      ...this.state.rules,
      empty
    ];
    this.setState({
      rules: rules
    });
    this.props.callback("WindowsFirewallRules", rules)
  }

  remove(id) {
    let rules = this.state.rules.filter(function(_, index) {
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
      }
      else {
        value = false;
      }
    }
    updated[id] = {
      ...updated[id],
      [field]: value
    }
    this.setState({
      rules: updated
    })
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
        ruleOptions = (
          <React.Fragment>
            <li>
              <label>Enabled</label>
              <input type="checkbox" checked={entry.Enabled} onChange={event=> this.update(i, "Enabled", event)}></input>
            </li>
            <li>
              <label>Protocol</label>
              <input type="text" value={entry.Protocol} onChange={event=> this.update(i, "Protocol", event)}></input>
            </li>
            <li>
              <label>Local Port</label>
              <input type="text" value={entry.LocalPort} onChange={event=> this.update(i, "LocalPort", event)}></input>
            </li>
            <li>
              <label>Remote Address</label>
              <input type="text" value={entry.RemoteAddress} onChange={event=> this.update(i, "RemoteAddress", event)}></input>
            </li>
            <li>
              <label>Remote Port</label>
              <input type="text" value={entry.RemotePort} onChange={event=> this.update(i, "RemotePort", event)}></input>
            </li>
            <li>
              <label>Direction</label>
              <select value={entry.Direction} onChange={event=> this.update(i, "Direction", event)}>
                <option disabled key="" value="">
                </option>
                <option>Inbound</option>
                <option>Outbound</option>
              </select>
            </li>
            <li>
              <label>Action</label>
              <select value={entry.Action} onChange={event=> this.update(i, "Action", event)}>
                <option disabled key="" value="">
                </option>
                <option>Block</option>
                <option>Allow</option>
              </select>
            </li>
          </React.Fragment>
        );
      }
      rules.push(
        <details key={i}>
          <summary>Display Name: {entry.DisplayName}, {enabledStr}, Protocol: {entry.Protocol}, Local Port: {entry.LocalPort}, Remote Address: {entry.RemoteAddress}, Remote Port: {entry.RemotePort}, Direction: {entry.Direction}, Action: {entry.Action}</summary>
          <button type="button" onClick={this.remove.bind(this, i)}>-</button>
          <ul>
            <li>
              <label>DisplayName</label>
              <input type="text" value={entry.DisplayName} onChange={event=> this.update(i, "DisplayName", event)}></input>
            </li>
            <li>
              <ObjectState value={entry.ObjectState} onChange={event=> this.update(i, "ObjectState", event)} />
            </li>
            {ruleOptions}
          </ul>
        </details>
      );
    }

    return (
      <details>
        <summary>Windows Firewall rules</summary>
        <button type="button" onClick={this.add.bind(this)}>Add Windows Firewall rule</button>
        <ul>
          {rules}
        </ul>
      </details>
    )
  }
}

class Item extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div>
        <label htmlFor={this.props.name}>{this.props.name}</label>
        <input name={this.props.name} type={this.props.type} value={this.props.value} checked={this.props.checked} disabled={this.props.disabled}></input>
      </div>
    )
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
    }

    this.add = this.add.bind(this);
    this.remove = this.remove.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.handleCallback = this.handleCallback.bind(this);
  }

  componentWillReceiveProps(newProps) {
    if (this.props.value != newProps.value) {
      this.setState({
        value: newProps.value
      })
    }
  }

  handleChange(event) {
    let value = Number(event.target.value)
    this.setState({
      item: value
    });
  }

  handleCallback(key, value) {
    let v = {
      ...this.state.value,
      [key]: value
    };
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

    let value = {
      ...this.state.value,
      [this.state.item]: []
    }
    this.setState({
      item: "",
      value: value
    })
    this.props.callback(this.props.name, value);
  }

  remove(id) {
    if (this.state.value == null) {
      return;
    }

    let value = {
      ...this.state.value,
      [id]: undefined
    }
    this.setState({
      value: value
    });
    this.props.callback(this.props.name, value);
  }

  componentWillMount() {
    this.props.mapItems((items) => {
      this.setState({
        mapItems: items
      });
    });
    this.props.listItems((items) => {
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
        let matches = this.state.mapItems.filter((obj) => {
          return obj.ID == i;
        });
        if (matches.length > 0) {
          text = matches[0].Display;
        }
        rows.push(
          <details key={i}>
            <summary>{text}</summary>
            <button type="button" onClick={this.remove.bind(this, i)}>-</button>
            <ul>
              <ItemList name={i} label={this.props.listLabel} type="select" listItems={this.state.listItems} value={this.state.value[i]} callback={this.handleCallback}/>
            </ul>
          </details>
        );
      }
    }

    let optionsMap = [];
    // empty selection
    optionsMap.push(
      <option disabled key="" value="">
      </option>
    );
    for (let i in this.state.mapItems) {
      let option = this.state.mapItems[i];
      // skip already selected
      if (this.state.value && this.state.value[option.ID] != null) {
        continue;
      }
      optionsMap.push(
        <option key={option.ID} value={option.ID}>
          {option.Display}
        </option>
      );
    }

    return (
      <div>
        <label>{this.props.label}</label>
        <ul>
          {rows}
          <select value={this.state.item} onChange={this.handleChange}>{optionsMap}</select>
          <button type="button" onClick={this.add}>+</button>
        </ul>
      </div>
    );
  }
}

class ItemList extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      item: "",
      value: this.props.value
    }

    this.add = this.add.bind(this);
    this.remove = this.remove.bind(this);
    this.handleChange = this.handleChange.bind(this);
  }

  componentWillReceiveProps(newProps) {
    if (this.props.value != newProps.value) {
      this.setState({
        value: newProps.value
      })
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
    }
    else  {
      value = [...this.state.value, this.state.item];
    }
    this.setState({
      item: "",
      value: value
    })
    this.props.callback(this.props.name, value);
  }

  remove(id) {
    if (this.state.value == null) {
      return;
    }

    let value = this.state.value.filter(function(_, index) {
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
          let matches = this.props.listItems.filter((obj) => {
            return obj.ID == text;
          });
          if (matches.length > 0) {
            text = matches[0].Display;
          }
        }
        rows.push(
          <li key={i}>
            {text}
            <button type="button" onClick={this.remove.bind(this, i)}>-</button>
          </li>
        );
      }
    }

    let input = (
      <input type={this.props.type} value={this.state.item} onChange={this.handleChange}></input>
    );
    if (this.props.type === "select") {
      let optionsList = [];
      // empty selection
      optionsList.push(
        <option disabled key="" value="">
        </option>
      );
      for (let i in this.props.listItems) {
        let option = this.props.listItems[i];
        // skip already selected
        if (this.state.value && this.state.value.indexOf(option.ID) != -1) {
          continue;
        }
        optionsList.push(
          <option key={option.ID} value={option.ID}>
            {option.Display}
          </option>
        );
      }
      input = (
        <select value={this.state.item} onChange={this.handleChange}>{optionsList}</select>
      );
    }

    return (
      <details>
        <summary>{this.props.label}</summary>
        <ul>
          {rows}
          {input}
          <button type="button" onClick={this.add}>+</button>
        </ul>
      </details>
    );
  }
}

ReactDOM.render(<App />, document.getElementById('app'));