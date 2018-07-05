'use strict';

class App extends React.Component {
  render() {
    return (
      <div className="App">
        <Hosts />

        <Templates />
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
      selectedHostID: null,
      selectedHost: null
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
      selectedHost: null
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
        <HostModal hostID={this.state.selectedHostID} host={this.state.selectedHost} show={this.state.showModal} onClose={this.toggleModal} submit={this.handleSubmit}/>
        <ul>{rows}</ul>
      </div>
    );
  }
}

class HostModal extends React.Component {
  constructor(props) {
    super(props);
    this.state = this.defaultState();

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleClose = this.handleClose.bind(this);
  }

  defaultState() {
    return {
      host: {}
    }
  }

  handleChange(event) {
    this.setState({
      host: {
        ...this.props.host,
        ...this.state.host,
        [event.target.name]: event.target.value
      }
    })
  }

  handleSubmit(event) {
    event.preventDefault();

    if (Object.keys(this.state.host) == 0) {
      return;
    }

    var url = "/hosts";
    if (this.props.hostID != null) {
      url += "/" + this.props.hostID;
    }

    fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(this.state.host)
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

    var host = {};
    if (this.props.hostID != null) {
      host = this.props.host;
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
      padding: 30
    }

    return (
      <div className="background" style={backgroundStyle}>
        <div className="modal" style={modalStyle}>
          <form onChange={this.handleChange} onSubmit={this.handleSubmit}>
            <label htmlFor="ID">ID</label>
            <input name="ID" defaultValue={host.ID} disabled="true"></input>
            <br />
            <label htmlFor="Hostname">Hostname</label>
            <input name="Hostname" defaultValue={host.Hostname}></input>
            <br />
            <label htmlFor="OS">OS</label>
            <input name="OS" defaultValue={host.OS}></input>
            <br />
            <button type="submit">Submit</button>
            <button type="button" onClick={this.handleClose}>Cancel</button>
          </form>
        </div>
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
      selectedTemplateID: null,
      selectedTemplate: null
    };

    this.handleSubmit = this.handleSubmit.bind(this);
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
      selectedTemplate: null
    });
    this.toggleModal();
  }

  editTemplate(id, template) {
    this.setState({
      selectedTemplateID: id,
      selectedTemplate: template
    });
    this.toggleModal();
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

  render() {
    let rows = [];
    for (let i = 0; i < this.state.templates.length; i++) {
      let entry = this.state.templates[i];
      Object.keys(entry).map(id => {
        rows.push(
          <li key={id}>
            {id} - {entry[id].Name}
            <button type="button" onClick={this.editTemplate.bind(this, id, entry[id])}>Edit</button>
            <button type="button" onClick={this.deleteTemplate.bind(this, id)}>-</button>
          </li>
        );
      });
    }
    return (
      <div className="Templates">
        <strong>Templates</strong>
        <p />
        <button onClick={this.createTemplate.bind(this)}>Create Template</button>
        <TemplateModal templateID={this.state.selectedTemplateID} template={this.state.selectedTemplate} show={this.state.showModal} onClose={this.toggleModal} submit={this.handleSubmit}/>
        <ul>{rows}</ul>
      </div>
    );
  }
}

class TemplateModal extends React.Component {
  constructor(props) {
    super(props);
    this.state = this.defaultState();

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleClose = this.handleClose.bind(this);
    this.callbackUsersAdd = this.callbackUsersAdd.bind(this);
    this.callbackUsersKeep = this.callbackUsersKeep.bind(this);
    this.callbackUsersRemove = this.callbackUsersRemove.bind(this);
  }

  defaultState() {
    return {
      template: {}
    }
  }

  handleChange(event) {
    this.setState({
      template: {
        ...this.props.template,
        ...this.state.template,
        [event.target.name]: event.target.value
      }
    })
  }

  handleSubmit(event) {
    event.preventDefault();

    if (Object.keys(this.state.template) == 0) {
      return;
    }

    var url = "/templates";
    if (this.props.templateID != null) {
      url += "/" + this.props.templateID;
    }

    fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(this.state.template)
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

  callbackUsersAdd(users) {
    this.setState({
      template: {
        ...this.props.template,
        ...this.state.template,
        UsersAdd: users
      }
    });
  }

  callbackUsersKeep(users) {
    this.setState({
      template: {
        ...this.props.template,
        ...this.state.template,
        UsersKeep: users
      }
    });
  }

  callbackUsersRemove(users) {
    this.setState({
      template: {
        ...this.props.template,
        ...this.state.template,
        UsersRemove: users
      }
    });
  }

  render() {
    if (!this.props.show) {
      return null;
    }

    var template = {};
    if (this.props.templateID != null) {
      template = this.props.template;
    }
    template = Object.assign({}, template, this.state.template);

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
      padding: 30
    }

    return (
      <div className="background" style={backgroundStyle}>
        <div className="modal" style={modalStyle}>
          <form onChange={this.handleChange} onSubmit={this.handleSubmit}>
            <label htmlFor="Name">Name</label>
            <input name="Name" defaultValue={template.Name}></input>
            <br />
            <ItemList label="Users to add" items={template.UsersAdd} callback={this.callbackUsersAdd}/>
            <br />
            <ItemList label="Users to keep" items={template.UsersKeep} callback={this.callbackUsersKeep}/>
            <br />
            <ItemList label="Users to remove" items={template.UsersRemove} callback={this.callbackUsersRemove}/>
            <br />
            <button type="submit">Submit</button>
            <button type="button" onClick={this.handleClose}>Cancel</button>
          </form>
        </div>
      </div>
    );
  }
}

class ItemList extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      item: ""
    }

    this.add = this.add.bind(this);
    this.remove = this.remove.bind(this);
    this.handleChange = this.handleChange.bind(this);
  }

  handleChange(event) {
    this.setState({
      item: event.target.value
    });
  }

  add() {
    if (!this.state.item) {
      return;
    }
    if (this.props.items && this.props.items.includes(this.state.item)) {
      return;
    }

    if (this.props.items == null) {
      this.props.callback([this.state.item]);
    }
    else {
      this.props.callback([...this.props.items, this.state.item]);
    }
  }

  remove(id) {
    if (this.props.items == null) {
      return;
    }

    let newItems = this.props.items.filter(function(_, index) {
      return index != id;
    });
    this.props.callback(newItems);
  }

  render() {
    let rows = [];
    if (this.props.items) {
      for (let i = 0; i < this.props.items.length; i++) {
        rows.push(
          <li key={i}>
            {this.props.items[i]}
            <button type="button" onClick={this.remove.bind(this, i)}>-</button>
          </li>
        );
      }
    }

    return (
      <div>
        <label>{this.props.label}</label>
        <input onChange={this.handleChange}></input>
        <button type="button" onClick={this.add}>+</button>
        <br />
        <ul>{rows}</ul>
      </div>
    );
  }
}

ReactDOM.render(<App />, document.getElementById('app'));