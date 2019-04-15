'use strict';

const Plot = createPlotlyComponent(Plotly);

class App extends React.Component {
  constructor() {
    super();
    this.state = {
      authenticated: false,
      args: {}
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

  analysisRequestCallback(documentType, args) {
    this.setState({
      documentType: documentType,
      args: args
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

    return (
      <div className="App">
        <div className="heading">
          <h1>cp-scoring Insight</h1>
        </div>
        <div className="navbar">
          <button className="right" onClick={this.logout}>Logout</button>
        </div>
        <div className="toc">
          <Analysis requestCallback={this.analysisRequestCallback.bind(this)}/>
        </div>
        <div className="content">
          <AnalysisItem documentType={this.state.documentType} args={this.state.args}/>
        </div>
      </div>
    );
  }
}

class Analysis extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      documentType: "",
      scenarios: [],
      scenarioID: 0,
      teams: [],
      teamID: 0,
      hosts: [],
      hostname: "",
      timeStart: Date.now() * 1000,
      timeEnd: Date.now() * 1000
    }

    this.selectDocumentTypeCallback = this.selectDocumentTypeCallback.bind(this);
    this.selectScenarioCallback = this.selectScenarioCallback.bind(this);
    this.selectTeamCallback = this.selectTeamCallback.bind(this);
    this.selectHostCallback = this.selectHostCallback.bind(this);
    this.selectTimeStartCallback = this.selectTimeStartCallback.bind(this);
    this.selectTimeEndCallback = this.selectTimeEndCallback.bind(this);
    this.updateTime = this.updateTime.bind(this);
    this.submit = this.submit.bind(this);
  }
  
  componentDidMount() {
    this.populateSelectors();
  }

  componentWillReceiveProps(_) {
    this.populateSelectors();
  }

  selectDocumentTypeCallback(event) {
    event.preventDefault();

    this.setState({
      documentType: event.target.value
    })
  }

  selectScenarioCallback(event) {
    event.preventDefault();

    this.setState({
      scenarioID: event.target.value
    })
  }

  selectTeamCallback(event) {
    event.preventDefault();

    this.setState({
      teamID: event.target.value
    })
  }

  selectHostCallback(event) {
    event.preventDefault();

    this.setState({
      hostname: event.target.value
    })
  }

  selectTimeStartCallback(event) {
    event.preventDefault();

    let updated = this.updateTime(event, this.state.timeStart);
    if (updated === null) {
      return;
    }

    this.setState({
      timeStart: updated
    })
  }

  selectTimeEndCallback(event) {
    event.preventDefault();

    let updated = this.updateTime(event, this.state.timeEnd);
    if (updated === null) {
      return;
    }

    this.setState({
      timeEnd: updated
    })
  }

  updateTime(event, original) {
    let current = new Date(Math.trunc(original * 1000));


    if (event.target.type === "date") {
      let parts = event.target.value.split("-");
      if (parts.length != 3) {
        return null;
      }
      current.setFullYear(parts[0]);
      // months start counting at 0
      current.setMonth(parts[1] - 1);
      current.setDate(parts[2]);
    }
    else if (event.target.type === "time") {
      let parts = event.target.value.split(":");
      if (parts.length != 3) {
        return null;
      }
      current.setHours(parts[0]);
      current.setMinutes(parts[1]);
      current.setSeconds(parts[2]);
    }

    let value = Math.trunc(current.getTime() / 1000);
    if (Number.isNaN(value)) {
      return null;
    }

    return value;
  }

  populateSelectors() {
    fetch('/scenarios', {
      credentials: 'same-origin'
    })
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json();
    })
    .then(function(data) {
      this.setState({scenarios: data})
    }.bind(this));

    fetch('/teams', {
      credentials: 'same-origin'
    })
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json();
    })
    .then(function(data) {
      this.setState({teams: data})
    }.bind(this));

    fetch('/hosts', {
      credentials: 'same-origin'
    })
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

  submit() {
    let args = {
      'scenario_id': this.state.scenarioID,
      'team_id': this.state.teamID,
      'hostname': this.state.hostname,
      'time_start': this.state.timeStart,
      'time_end': this.state.timeEnd,
    }
    this.props.requestCallback(this.state.documentType, args);
  }

  render() {
    // form document type options
    let documentTypeOptions = [];
    documentTypeOptions.push(<option key="-1" value=""></option>);
    documentTypeOptions.push(<option key="0" value="reports">Reports</option>);
    documentTypeOptions.push(<option key="1" value="states">States</option>);

    // form scenario options
    let scenarioOptions = [];
    scenarioOptions.push(<option key="-1" value=""></option>);
    for (let i in this.state.scenarios) {
      let scenario = this.state.scenarios[i];
      scenarioOptions.push(<option key={i} value={scenario.ID}>{scenario.Name}</option>);
    }

    // form team options
    let teamOptions = [];
    teamOptions.push(<option key="-1" value=""></option>);
    for (let i in this.state.teams) {
      let team = this.state.teams[i];
      teamOptions.push(<option key={i} value={team.ID}>{team.Name}</option>);
    }

    // form host options
    let hostOptions = [];
    hostOptions.push(<option key="-1" value=""></option>);
    for (let i in this.state.hosts) {
      let host = this.state.hosts[i];
      hostOptions.push(<option key={i} value={host.Hostname}>{host.Hostname}</option>);
    }

    // form time start
    let d = new Date(this.state.timeStart * 1000);
    let startDate = ("000" + d.getFullYear()).slice(-4);
    startDate += "-";
    startDate += ("0" + (d.getMonth() + 1)).slice(-2);
    startDate += "-";
    startDate += ("0" + d.getDate()).slice(-2);
    let startTime = ("000" + d.getHours()).slice(-2);
    startTime += ":";
    startTime += ("000" + d.getMinutes()).slice(-2);
    startTime += ":";
    startTime += ("000" + d.getSeconds()).slice(-2);
    // form time end
    d = new Date(this.state.timeEnd * 1000);
    let endDate = ("000" + d.getFullYear()).slice(-4);
    endDate += "-";
    endDate += ("0" + (d.getMonth() + 1)).slice(-2);
    endDate += "-";
    endDate += ("0" + d.getDate()).slice(-2);
    let endTime = ("000" + d.getHours()).slice(-2);
    endTime += ":";
    endTime += ("000" + d.getMinutes()).slice(-2);
    endTime += ":";
    endTime += ("000" + d.getSeconds()).slice(-2);

    return (
      <React.Fragment>
        <label name="type">Document Type</label>
        <select value={this.state.documentType} onChange={this.selectDocumentTypeCallback}>
          {documentTypeOptions}
        </select>
        <br />
        <label name="scenarios">Scenarios</label>
        <select value={this.state.scenarioID} onChange={this.selectScenarioCallback}>
          {scenarioOptions}
        </select>
        <br />
        <label name="teams">Teams</label>
        <select value={this.state.teamID} onChange={this.selectTeamCallback}>
          {teamOptions}
        </select>
        <br />
        <label name="hosts">Hosts</label>
        <select value={this.state.hostname} onChange={this.selectHostCallback}>
          {hostOptions}
        </select>
        <br />
        <label name="timeStart">Time start</label>
        <input type="date" value={startDate} onChange={this.selectTimeStartCallback}/>
        <input type="time" value={startTime} onChange={this.selectTimeStartCallback}/>
        <br />
        <label name="timeEnd">Time end</label>
        <input type="date" value={endDate} onChange={this.selectTimeEndCallback}/>
        <input type="time" value={endTime} onChange={this.selectTimeEndCallback}/>
        <p />
        <button onClick={this.submit}>Submit</button>
      </React.Fragment>
    );
  }
}

class AnalysisItem extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      timeline: {},
      diffs: {},
      selected: {}
    }

    this.plotClick = this.plotClick.bind(this);
  }
  
  componentDidMount() {
    this.getData();
  }

  componentWillReceiveProps(_) {
    this.getData();
  }

  getData() {
    if (this.props.args === null || this.props.args === undefined) {
      return;
    }

    let urlTimeline = null;
    let urlDiffs = null;
    if (this.props.documentType === 'reports') {
      urlTimeline = '/analysis/reports/timeline';
      urlDiffs = '/analysis/reports/diffs';
    }
    else if (this.props.documentType === 'states') {
      urlTimeline = '/analysis/states/timeline';
      urlDiffs = '/analysis/states/diffs';
    }
    else {
      return;
    }

    let params = Object.entries(this.props.args).map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`).join('&');
    urlTimeline = urlTimeline + '?' + params;
    urlDiffs = urlDiffs + '?' + params;

    fetch(urlTimeline, {
      credentials: 'same-origin',
    })
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json();
    })
    .then(function(data) {
      this.setState({timeline: data})
    }.bind(this));

    fetch(urlDiffs, {
      credentials: 'same-origin',
    })
    .then(function(response) {
      if (response.status >= 400) {
        throw new Error("Bad response from server");
      }
      return response.json();
    })
    .then(function(data) {
      this.setState({diffs: data})
    }.bind(this));
  }

  plotClick(plotlyEvent) {
    // only accept left click
    if (plotlyEvent.event.buttons != 1) {
      return;
    }

    let i = plotlyEvent.points[0].pointIndex;
    

    let type = plotlyEvent.points[0].y;
    if (type.endsWith('(diff)')) {
      // diffs should have been previously retrieved
      this.setState({
        // choose first instance
        selected: this.state.diffs[0][i]
      });
    } else if (type.endsWith('(reports)') || type.endsWith('(states)')) {
      // get report/state
      let args = {...this.props.args};
      delete args['time_start'];
      delete args['time_end'];
      args['timestamp'] = Math.trunc(plotlyEvent.points[0].data.x[i] / 1000);

      let params = Object.entries(args).map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`).join('&');
      let url = '/analysis/' + this.props.documentType + '?' + params;
      
      fetch(url, {
        credentials: 'same-origin',
      })
      .then(function(response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        return response.json();
      })
      .then(function(data) {
        // choose first instance
        this.setState({selected: data[0]})
      }.bind(this));
    }

    return false;
  }

  render() {
    let config = {
      displaylogo: false
    }

    let layout = {
      hovermode: 'closest',
      xaxis: {
        type: 'date'
      },
      yaxis: {
        autorange: 'reversed'
      }
    }
    let traces = [];

    let name = this.props.args.hostname + ' (' + this.props.documentType + ')';
    for (let i in this.state.timeline) {
      let hostInstance = this.state.timeline[i];
      let trace = {
        name: name,
        mode: 'markers',
        x: hostInstance.map(document => document.Document * 1000),
        y: hostInstance.map(_ => name)
      }
      traces.push(trace);
    }

    // plot diffs
    name = this.props.args.hostname + ' (diff)';
    for (let i in this.state.diffs) {
      let hostInstance = this.state.diffs[i];
      let trace = {
        name: name,
        mode: 'markers',
        x: hostInstance.map(diff => diff.Timestamp * 1000),
        y: hostInstance.map(_ => name)
      }
      traces.push(trace)
    }

    let selected = null;
    // diff
    if (this.state.selected.Changes != undefined) {
      let time = new Date(this.state.selected.Timestamp * 1000).toLocaleString();
      let changes = [];
      for (let i in this.state.selected.Changes) {
        let change = this.state.selected.Changes[i];
        changes.push(<li key={i}>{change.Type} - {change.Key} - {change.Item}</li>)
      }
      selected = (
        <React.Fragment>
          Time: {time}
          <br />
          Changes:
          <ul>
            {changes}
          </ul>
        </React.Fragment>
      );
    }
    // report
    else if (this.props.documentType === 'reports') {
      let time = new Date(this.state.selected.Timestamp * 1000).toLocaleString();
      let findings = [];
      for (let i in this.state.selected.Findings) {
        let finding = this.state.selected.Findings[i];
        findings.push(<li key={i}>{finding.Show} - {finding.Value} - {finding.Message}</li>)
      }
      selected = (
        <React.Fragment>
          Time: {time}
          <br />
          Findings:
          <ul>
            {findings}
          </ul>
        </React.Fragment>
      );
    }
    // state
    else if (this.props.documentType === 'states') {
      let time = new Date(this.state.selected.Timestamp * 1000).toLocaleString();
      let errors = [];
      for (let i in this.state.selected.Errors) {
        let error = this.state.selected.Errors[i];
        errors.push(<li key={i}>{error}</li>);
      }
      let users = [];
      for (let i in this.state.selected.Users) {
        let user = this.state.selected.Users[i];
        users.push(<li key={i}>{user.Name}</li>);
      }
      let groups = [];
      let software = [];
      let processes = [];
      let conns = [];
      selected = (
        <React.Fragment>
          Time: {time}
          <br />
          OS: {this.state.selected.OS}
          <br />
          Hostname: {this.state.selected.Hostname}
          <br />
          Errors:
          <ul>
            {errors}
          </ul>
          <br />
          Users:
          <ul>
            {users}
          </ul>
          <br />
          Groups:
          <ul>
            {groups}
          </ul>
          <br />
          Software:
          <ul>
            {software}
          </ul>
          <br />
          Processes:
          <ul>
            {processes}
          </ul>
          <br />
          Network connections:
          <ul>
            {conns}
          </ul>
          <br />
        </React.Fragment>
      );
    }

    return (
      <React.Fragment>
      Timeline
      <ul>
        <Plot data={traces} layout={layout} config={config} onClick={this.plotClick}/>
      </ul>
      <p />
      Selected
      <ul>
        {selected}
      </ul>
      </React.Fragment>
    );
  }
}

ReactDOM.render(<App />, document.getElementById('app'));