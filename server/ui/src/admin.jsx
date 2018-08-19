'use strict';

const Plot = createPlotlyComponent(Plotly);

class App extends React.Component {
  constructor() {
    super();
    this.state = {
      authenticated: false
    }

    this.authCallback = this.authCallback.bind(this);
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
  }

  render() {
    if (!this.state.authenticated) {
      return (
        <div className="App">
          <Login callback={this.authCallback}/>
        </div>
      );
    }
    return (
      <div className="App">
        <button onClick={this.logout}>Logout</button>
        <p/>
        <Teams />

        <Hosts />

        <Templates />

        <Scenarios />
      </div>
    );
  }
}

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

const backgroundStyle = {
  position: 'fixed',
  top: 0,
  bottom: 0,
  left: 0,
  right: 0,
  backgroundColor: 'rgba(0,0,0,0.5)',
  padding: 50
}

const modalStyle = {
  backgroundColor: 'white',
  padding: 30,
  maxHeight: '100%',
  overflowY: 'auto',
}

class BasicModal extends React.Component {
  constructor(props) {
    super(props);
    this.state = this.defaultState();

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleClose = this.handleClose.bind(this);
    this.setValue = this.setValue.bind(this);
  }

  defaultState() {
    return {
      subject: {}
    }
  }

  setValue(key, value) {
    this.setState({
      subject: {
        ...this.props.subject,
        ...this.state.subject,
        [key]: value
      }
    });
  }

  handleChange(event) {
    let value = event.target.value;
    if (event.target.type == 'checkbox') {
      value = event.target.checked;
    }
    this.setState({
      subject: {
        ...this.props.subject,
        ...this.state.subject,
        [event.target.name]: value
      }
    });
  }

  handleSubmit(event) {
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
    })
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      this.props.submit();
      this.setState(this.defaultState());
    }.bind(this));
  }

  handleClose() {
    this.props.onClose();
    this.setState(this.defaultState());
  }
  
  render() {
    if (!this.props.show) {
      return null;
    }

    return (
      <div className="background" style={backgroundStyle}>
        <div className="modal" style={modalStyle}>
          <label htmlFor="ID">ID</label>
          <input name="ID" defaultValue={this.props.subjectID} disabled></input>
          <br />
          <form onChange={this.handleChange} onSubmit={this.handleSubmit}>
            {this.props.children}
            <br />
            <button type="submit">Submit</button>
            <button type="button" onClick={this.handleClose}>Cancel</button>
          </form>
        </div>
      </div>
    );
  }
}

class Teams extends React.Component {
  constructor() {
    super();
    this.state = {
      teams: [],
      showModal: false,
      selectedTeamID: null,
      selectedTeam: {}
    }

    this.modal = React.createRef()
    this.handleSubmit = this.handleSubmit.bind(this);
    this.regenKey = this.regenKey.bind(this);
  }

  componentDidMount() {
    this.populateTeams();
  }

  populateTeams() {
    var url = '/teams';
  
    fetch(url)
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json();
    })
    .then(function(data) {
      this.setState({teams: data})
    }.bind(this));
  }

  newKey() {
    return Math.random().toString(16).substring(7);
  }

  createTeam() {
    this.setState({
      selectedTeamID: null,
      selectedTeam: {
        Enabled: true,
        Key: this.newKey()
      }
    });
    this.toggleModal();
  }

  editTeam(id) {
    let url = "/teams/" + id;

    fetch(url)
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json()
    })
    .then(function(data) {
      this.setState({
        selectedTeamID: id,
        selectedTeam: data
      });
      this.toggleModal();
    }.bind(this));
  }

  deleteTeam(id) {
    var url = "/teams/" + id;

    fetch(url, {
      method: 'DELETE'
    })
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      this.populateTeams();
    }.bind(this));
  }

  handleSubmit() {
    this.populateTeams();
    this.toggleModal();
  }

  toggleModal = () => {
    this.setState({
      showModal: !this.state.showModal
    })
  }

  regenKey() {
    let key = this.newKey()
    this.setState({
      selectedTeam: {
        ...this.state.selectedTeam,
        Key: key
      }
    })
    this.modal.current.setValue("Key", key)
  }

  render() {
    let rows = [];
    for (let i = 0; i < this.state.teams.length; i++) {
      let team = this.state.teams[i];
      rows.push(
        <li key={team.ID}>
          {team.ID} - {team.Name}
          <button type="button" onClick={this.editTeam.bind(this, team.ID)}>Edit</button>
          <button type="button" onClick={this.deleteTeam.bind(this, team.ID)}>-</button>
        </li>
      );
    }

    return (
      <div className="Teams">
        <strong>Teams</strong>
        <p />
        <button onClick={this.createTeam.bind(this)}>Add Team</button>
        <BasicModal ref={this.modal} subjectClass="teams" subjectID={this.state.selectedTeamID} subject={this.state.selectedTeam} show={this.state.showModal} onClose={this.toggleModal} submit={this.handleSubmit}>
          <Item name="Name" defaultValue={this.state.selectedTeam.Name}/>
          <Item name="POC" defaultValue={this.state.selectedTeam.POC}/>
          <Item name="Email" type="email" defaultValue={this.state.selectedTeam.Email}/>
          <label htmlFor="Enabled">Enabled</label>
          <input name="Enabled" type="checkbox" defaultChecked={!!this.state.selectedTeam.Enabled}></input>
          <br />
          <details>
            <summary>Key</summary>
            <ul>
              <li>
                {this.state.selectedTeam.Key}
                <button type="button" onClick={this.regenKey.bind(this)}>Regenerate</button>
              </li>
            </ul>
          </details>
        </BasicModal>
        <ul>{rows}</ul>
      </div>
    );
  }
}

class Scenarios extends React.Component {
  constructor() {
    super();
    this.state = {
      scenarios: [],
      showModal: false,
      selectedScenario: {}
    }
    this.modal = React.createRef();

    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleCallback = this.handleCallback.bind(this);
    this.mapItems = this.mapItems.bind(this);
    this.listItems = this.listItems.bind(this);
  }

  componentDidMount() {
    this.populateScenarios();
  }

  populateScenarios() {
    var url = '/scenarios';
  
    fetch(url)
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json();
    })
    .then(function(data) {
      this.setState({scenarios: data})
    }.bind(this));
  }

  createScenario() {
    this.setState({
      selectedScenarioID: null,
      selectedScenario: {Enabled: true}
    });
    this.toggleModal();
  }

  editScenario(id) {
    let url = "/scenarios/" + id;

    fetch(url)
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json()
    })
    .then(function(data) {
      this.setState({
        selectedScenarioID: id,
        selectedScenario: data
      });
      this.toggleModal();
    }.bind(this));
  }

  deleteScenario(id) {
    var url = "/scenarios/" + id;

    fetch(url, {
      method: 'DELETE'
    })
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      this.populateScenarios();
    }.bind(this));
  }

  handleSubmit() {
    this.populateScenarios();
    this.toggleModal();
  }

  toggleModal = () => {
    this.setState({
      showModal: !this.state.showModal
    })
  }

  handleCallback(key, value) {
    this.modal.current.setValue(key, value);
  }

  mapItems(callback) {
    var url = "/hosts";

    fetch(url)
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json()
    }.bind(this))
    .then(function(data) {
      let items = data.map(function(host) {
        return {
          ID: host.ID,
          Display: host.Hostname
        }
      });
      callback(items);
    });
  };

  listItems(callback) {
    var url = "/templates";

    fetch(url)
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json()
    }.bind(this))
    .then(function(data) {
      let items = data.map(function(template) {
        return {
          ID: template.ID,
          Display: template.Name
        }
      });
      callback(items);
    });
  };

  render() {
    let rows = [];
    for (let i = 0; i < this.state.scenarios.length; i++) {
      let scenario = this.state.scenarios[i];
      rows.push(
        <li key={scenario.ID}>
          {scenario.ID} - {scenario.Name}
          <button type="button" onClick={this.editScenario.bind(this, scenario.ID)}>Edit</button>
          <button type="button" onClick={this.deleteScenario.bind(this, scenario.ID)}>-</button>
        </li>
      );
    }

    return (
      <div className="Scenarios">
        <strong>Scenarios</strong>
        <p />
        <button onClick={this.createScenario.bind(this)}>Add Scenario</button>
        <BasicModal ref={this.modal} subjectClass="scenarios" subjectID={this.state.selectedScenarioID} subject={this.state.selectedScenario} show={this.state.showModal} onClose={this.toggleModal} submit={this.handleSubmit}>
          <Item name="Name" defaultValue={this.state.selectedScenario.Name}/>
          <Item name="Description" defaultValue={this.state.selectedScenario.Description}/>
          <Item name="Enabled" type="checkbox" defaultChecked={!!this.state.selectedScenario.Enabled}/>
          <ItemMap name="HostTemplates" label="Hosts" listLabel="Templates" defaultValue={this.state.selectedScenario.HostTemplates} callback={this.handleCallback} mapItems={this.mapItems} listItems={this.listItems}/>
        </BasicModal>
        <ul>{rows}</ul>
      </div>
    );
  }
}

class Hosts extends React.Component {
  constructor() {
    super();
    this.state = {
      hosts: [],
      showModal: false,
      selectedHost: {}
    };

    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDidMount() {
    this.populateHosts();
  }

  populateHosts() {
    var url = '/hosts';
  
    fetch(url)
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json();
    })
    .then(function(data) {
      this.setState({hosts: data})
    }.bind(this));
  }

  createHost() {
    this.setState({
      selectedHostID: null,
      selectedHost: {}
    });
    this.toggleModal();
  }

  editHost(id, host) {
    this.setState({
      selectedHostID: id,
      selectedHost: host
    });
    this.toggleModal();
  }

  deleteHost(id) {
    var url = "/hosts/" + id;

    fetch(url, {
      method: 'DELETE'
    })
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      this.populateHosts();
    }.bind(this));
  }

  handleSubmit() {
    this.populateHosts();
    this.toggleModal();
  }

  toggleModal = () => {
    this.setState({
      showModal: !this.state.showModal
    })
  }

  render() {
    let rows = [];
    for (let i = 0; i < this.state.hosts.length; i++) {
      let host = this.state.hosts[i];
      rows.push(
        <li key={host.ID}>
          {host.ID} - {host.Hostname} - {host.OS}
          <button type="button" onClick={this.editHost.bind(this, host.ID, host)}>Edit</button>
          <button type="button" onClick={this.deleteHost.bind(this, host.ID)}>-</button>
        </li>
      );
    }
  
    return (
      <div className="Hosts">
        <strong>Hosts</strong>
        <p />
        <button onClick={this.createHost.bind(this)}>Add Host</button>
        <BasicModal subjectClass="hosts" subjectID={this.state.selectedHostID} subject={this.state.selectedHost} show={this.state.showModal} onClose={this.toggleModal} submit={this.handleSubmit}>
          <Item name="Hostname" type="text" defaultValue={this.state.selectedHost.Hostname}/>
          <Item name="OS" type="text" defaultValue={this.state.selectedHost.OS}/>
        </BasicModal>
        <ul>{rows}</ul>
      </div>
    );
  }
}

class Templates extends React.Component {
  constructor() {
    super();
    this.state = {
      templates: [],
      showModal: false,
      selectedTemplate: {
        Template: {}
      }
    };
    this.modal = React.createRef();

    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleCallback = this.handleCallback.bind(this);
  }

  componentDidMount() {
    this.populateTemplates();
  }

  populateTemplates() {
    var url = "/templates";
  
    fetch(url)
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json();
    })
    .then(function(data) {
      this.setState({templates: data})
    }.bind(this));
  }

  createTemplate() {
    this.setState({
      selectedTemplateID: null,
      selectedTemplate: {
        Template: {}
      }
    });
    this.toggleModal();
  }

  editTemplate(id) {
    let url = "/templates/" + id;

    fetch(url)
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json()
    })
    .then(function(data) {
      this.setState({
        selectedTemplateID: id,
        selectedTemplate: data
      });
      this.toggleModal();
    }.bind(this));
  }

  deleteTemplate(id) {
    var url = "/templates/" + id;

    fetch(url, {
      method: 'DELETE'
    })
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      this.populateTemplates();
    }.bind(this));
  }

  handleSubmit() {
    this.populateTemplates();
    this.toggleModal();
  }

  toggleModal = () => {
    this.setState({
      showModal: !this.state.showModal
    })
  }

  handleCallback(key, value) {
    let template = {
      ...this.state.selectedTemplate.Template,
      [key]: value
    }
    this.setState({
      selectedTemplate: {
        ...this.state.selectedTemplate,
        Template: template
      }
    })
    this.modal.current.setValue("Template", template);
  }

  render() {
    let rows = [];
    for (let i = 0; i < this.state.templates.length; i++) {
      let template = this.state.templates[i];
      rows.push(
        <li key={template.ID}>
          {template.ID} - {template.Name}
          <button type="button" onClick={this.editTemplate.bind(this, template.ID)}>Edit</button>
          <button type="button" onClick={this.deleteTemplate.bind(this, template.ID)}>-</button>
        </li>
      );
    }

    return (
      <div className="Templates">
        <strong>Templates</strong>
        <p />
        <button onClick={this.createTemplate.bind(this)}>Create Template</button>
        <BasicModal ref={this.modal} subjectClass="templates" subjectID={this.state.selectedTemplateID} subject={this.state.selectedTemplate} show={this.state.showModal} onClose={this.toggleModal} submit={this.handleSubmit}>
          <Item name="Name" type="text" defaultValue={this.state.selectedTemplate.Name}/>
          <Users users={this.state.selectedTemplate.Template.Users} callback={this.handleCallback}/>
          <Groups name="GroupMembersAdd" label="Group members to add" groups={this.state.selectedTemplate.Template.GroupMembersAdd} callback={this.handleCallback}/>
          <Groups name="GroupMembersKeep" label="Group members to keep" groups={this.state.selectedTemplate.Template.GroupMembersKeep} callback={this.handleCallback}/>
          <Groups name="GroupMembersRemove" label="Group members to remove" groups={this.state.selectedTemplate.Template.GroupMembersRemove} callback={this.handleCallback}/>
        </BasicModal>
        <ul>{rows}</ul>
      </div>
    );
  }
}

class Users extends React.Component {
  constructor(props) {
    super(props);

    let users = props.users;
    if (users === undefined || users === null) {
      users = [];
    }
    this.state = {
      users: users
    }

    this.addUser = this.addUser.bind(this);
    this.removeUser = this.removeUser.bind(this);
    this.updateUser = this.updateUser.bind(this);
  }

  addUser() {
    let empty = {
      Name: "",
      AccountPresent: true,
      AccountActive: true,
      PasswordExpires: true,
      // unix timestamp in seconds
      PasswordLastSet: Date.now() / 1000
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
    if (event.target.type === "date") {
      value = new Date(event.target.value).getTime() / 1000
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
      let passwordLastSet = ("000" + d.getUTCFullYear()).slice(-4);
      passwordLastSet += "-";
      passwordLastSet += ("0" + (d.getUTCMonth() + 1)).slice(-2);
      passwordLastSet += "-";
      passwordLastSet += ("0" + d.getUTCDate()).slice(-2);
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
              <label>Present</label>
              <input type="checkbox" checked={user.AccountPresent} onChange={event=> this.updateUser(i, "AccountPresent", event)}/>
            </li>
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
              <input type="date" value={passwordLastSet} onChange={event=> this.updateUser(i, "PasswordLastSet", event)}/>
            </li>
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
    
    let groups = props.groups;
    if (groups === undefined || groups === null) {
      groups = {};
    }
    this.state = {
      groups: groups
    }

    this.newGroupName = React.createRef();

    this.addGroup = this.addGroup.bind(this);
    this.removeGroup = this.removeGroup.bind(this);
    this.updateGroup = this.updateGroup.bind(this);
  }

  addGroup() {
    if (this.newGroupName.current === null) {
      return;
    }
    let groups = {
      ...this.state.groups,
      [this.newGroupName.current.value]: []
    };
    this.setState({
      groups: groups
    });
    this.props.callback(this.props.name, groups)
  }

  removeGroup(name) {
    let groups = this.state.groups;
    delete groups[name];
    this.setState({
      groups: groups
    });
    this.props.callback(this.props.name, groups);
  }

  updateGroup(name, members) {
    let groups = {
      ...this.state.groups,
      [name]: members
    }
    this.setState({
      groups: groups
    });
    this.props.callback(this.props.name, groups);
  }

  render() {
    let groups = [];
    for (let groupName in this.state.groups) {
      let members = this.state.groups[groupName];
      groups.push(
        <details key={groupName}>
          <summary>{groupName}</summary>
          <button type="button" onClick={this.removeGroup.bind(this, groupName)}>-</button>
          <ItemList name={groupName} defaultValue={members} callback={this.updateGroup}/>
        </details>
      );
    }

    return (
      <details>
        <summary>{this.props.label}</summary>
        <input ref={this.newGroupName}></input>
        <button type="button" onClick={this.addGroup.bind(this)}>Add Group</button>
        <ul>
          {groups}
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
        <input name={this.props.name} type={this.props.type} defaultValue={this.props.defaultValue} defaultChecked={this.props.defaultChecked}></input>
      </div>
    )
  }
}

class ItemMap extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      item: "",
      value: this.props.defaultValue,
      mapItems: [],
      listItems: []
    }

    this.add = this.add.bind(this);
    this.remove = this.remove.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.handleCallback = this.handleCallback.bind(this);
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
          <li key={i}>
            {text}
            <button type="button" onClick={this.remove.bind(this, i)}>-</button>
            <ItemList name={i} label={this.props.listLabel} type="select" listItems={this.state.listItems} defaultValue={this.state.value[i]} callback={this.handleCallback}/>
          </li>
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
      value: this.props.defaultValue
    }

    this.add = this.add.bind(this);
    this.remove = this.remove.bind(this);
    this.handleChange = this.handleChange.bind(this);
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
      <div>
        <label>{this.props.label}</label>
        <ul>
          {rows}
          {input}
          <button type="button" onClick={this.add}>+</button>
        </ul>
      </div>
    );
  }
}

ReactDOM.render(<App />, document.getElementById('app'));