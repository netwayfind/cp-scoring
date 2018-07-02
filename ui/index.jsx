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
    this.state = {hosts: []};
  }

  componentDidMount() {    
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

  render() {
    return (
      <div className="Hosts">
        <strong>Hosts</strong>
        <ul>
          {this.state.hosts.map(host => {
            return <li>{host.ID} - {host.Hostname} - {host.OS}</li>
          })}
        </ul>
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
      selectedTemplateID: null
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
        <ul>{rows}</ul>
        <p />
        <button onClick={this.createTemplate.bind(this)}>Create Template</button>
        <TemplateModal templateID={this.state.selectedTemplateID} template={this.state.selectedTemplate} show={this.state.showModal} onClose={this.toggleModal} submit={this.handleSubmit}/>
      </div>
    );
  }
}

class TemplateModal extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      template: {}
    };

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    this.setState({
      template: {
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
    }.bind(this));
  }

  render() {
    if (!this.props.show) {
      return null;
    }

    var template = {};
    if (this.props.templateID != null) {
      template = this.props.template;
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
            <label htmlFor="Name">Name</label>
            <input name="Name" defaultValue={template.Name}></input>
            <br />
            <button type="submit">Submit</button>
            <button type="button" onClick={this.props.onClose}>Cancel</button>
          </form>
        </div>
      </div>
    );
  }
}

ReactDOM.render(<App />, document.getElementById('app'));