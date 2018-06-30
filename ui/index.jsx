'use strict';

class App extends React.Component {
  render() {
    return (
      <div>
        <Hosts />

        <Templates />

        <CreateTemplate />
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
    var t = this;
  
    fetch(url)
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json();
    })
    .then(function(data) {
      t.setState({hosts: data})
    });
  }

  render() {
    return (
      <div>
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
    this.state = {templates: []};
  }

  componentDidMount() {
    var url = '/templates';
    var t = this;
  
    fetch(url)
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json();
    })
    .then(function(data) {
      t.setState({templates: data})
    });
  }

  render() {
    return (
      <div>
        <strong>Templates</strong>
        <ul>
          {this.state.templates.map(template => {
            return <li>{template.Name}</li>
          })}
        </ul>
      </div>
    );
  }
}

class CreateTemplate extends React.Component {
  constructor(props) {
    super(props);
    this.state = {};

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    this.setState({
      [event.target.name]: event.target.value
    });
  }

  handleSubmit(event) {
    event.preventDefault();

    if (Object.keys(this.state) == 0) {
      return;
    }

    var url = "/templates";

    fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(this.state)
    })
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      window.location.reload();
    });
  }

  render() {
    return (
      <div>
        <strong>Create Template</strong>
        <form onChange={this.handleChange} onSubmit={this.handleSubmit}>
          <label for="Name">Name</label>
          <input name="Name" />
          <button type="submit">Submit</button>
        </form>
      </div>
    );
  }
}

ReactDOM.render(<App />, document.getElementById('app'));