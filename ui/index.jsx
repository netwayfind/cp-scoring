'use strict';

class App extends React.Component {
  render() {
    return <Hosts />;
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
      console.log(data);
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

ReactDOM.render(<App />, document.getElementById('app'));