'use strict';

class App extends React.Component {
  render() {
    return <Hosts />;
  }
}

class Hosts extends React.Component {
  render() {
    return (
      <ul>
        <li>host</li>
      </ul>
    )
  }
}

ReactDOM.render(<App />, document.getElementById('app'));