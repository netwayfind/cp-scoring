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

  analysisRequestCallback(args) {
    this.setState({
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
          <AnalysisItem args={this.state.args}/>
        </div>
      </div>
    );
  }
}

class Analysis extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      scenarios: [],
      scenarioID: 0,
      teams: [],
      teamID: 0,
      timeStart: Date.now() * 1000,
      timeEnd: Date.now() * 1000
    }

    this.selectScenarioCallback = this.selectScenarioCallback.bind(this);
    this.selectTeamCallback = this.selectTeamCallback.bind(this);
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
  }

  submit() {
    let args = {
      'scenario_id': this.state.scenarioID,
      'team_id': this.state.teamID,
      'time_start': this.state.timeStart,
      'time_end': this.state.timeEnd,
    }
    this.props.requestCallback(args);
  }

  render() {
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
      reportTimeline: {},
      reportDiffs: {},
      stateTimeline: {},
      stateDiffs: {},
      selected: {}
    }

    this.plotClick = this.plotClick.bind(this);
  }
  
  componentDidMount() {
    this.getData(this.props);
  }

  componentWillReceiveProps(newProps) {
    this.getData(newProps);
  }

  getData(props) {
    if (props.args === null || props.args === undefined) {
      return;
    }

    let urlTimeline = null;
    let urlDiffs = null;

    let params = Object.entries(props.args).map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`).join('&');
    urlTimeline = urlTimeline + '?' + params;
    urlDiffs = urlDiffs + '?' + params;

    let requestReports = fetch('/analysis/reports/timeline?' + params, {
      credentials: 'same-origin',
    });
    let requestReportDiffs = fetch('/analysis/reports/diffs?' + params, {
      credentials: 'same-origin',
    });
    let requestStates = fetch('/analysis/states/timeline?' + params, {
      credentials: 'same-origin',
    });
    let requestStateDiffs = fetch('/analysis/states/diffs?' + params, {
      credentials: 'same-origin',
    });
    Promise.all([requestReports, requestReportDiffs, requestStates, requestStateDiffs])
    .then(function(responses) {
      let j = [];
      for (let r in responses) {
        let response = responses[r];
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        j.push(response.json());
      }
      return Promise.all(j);
    })
    .then(function(data) {
      this.setState({
        reportTimeline: data[0],
        reportDiffs: data[1],
        stateTimeline: data[2],
        stateDiffs: data[3],
      });
    }.bind(this)); 
  }

  plotClick(plotlyEvent) {
    // only accept left click
    if (plotlyEvent.event.buttons != 1) {
      return;
    }

    let index = plotlyEvent.points[0].pointIndex;
    let timestamp = Math.trunc(plotlyEvent.points[0].data.x[index] / 1000);

    let type = plotlyEvent.points[0].y;
    if (type.endsWith('reports diff') || type.endsWith('states diff')) {
      let diffs = null;
      if (type.endsWith('reports diff')) {
        diffs = this.state.reportDiffs;
      }
      else if (type.endsWith('states diff')) {
        diffs = this.state.stateDiffs;
      }
      else {
        return false;
      }

      // find diff that matches timestamp
      let selected = null;
      for (let i in diffs) {
        if (diffs[i].length <= index) {
          continue
        }
        if (diffs[i][index].Timestamp != timestamp) {
          continue
        }
        selected = diffs[i][index];
      }

      if (selected === null) {
        selected = {}
      }
      this.setState({
        selected: selected
      });
    }
    else if (type.endsWith('reports') || type.endsWith('states')) {
      let documentType = null;
      let timeline = null;
      if (type.endsWith('reports')) {
        documentType = 'reports';
        timeline = this.state.reportTimeline;
      }
      else if (type.endsWith('states')) {
        documentType = 'states';
        timeline = this.state.stateTimeline;
      }
      else {
        return false;
      }

      // get report/state ID that matches timestamp and position
      let id = null;
      for (let i in timeline) {
        if (timeline[i].length <= index) {
          continue;
        }
        if (timeline[i][index].Document === timestamp) {
          id = timeline[i][index].ID;
        }
      }
      if (id === null) {
        this.setState({selected: {}});
        return false;
      }

      let url = '/analysis/' + documentType + '?id=' + id;
      
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
        this.setState({selected: data});
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
        autorange: 'reversed',
        visible: false
      }
    }
    let traces = [];

    // states
    for (let i in this.state.stateTimeline) {
      let hostInstance = this.state.stateTimeline[i];
      let name = i + ' - A.states';
      let trace = {
        name: name,
        mode: 'markers',
        x: hostInstance.map(document => document.Document * 1000),
        y: hostInstance.map(_ => name)
      }
      traces.push(trace);
    }

    // state diffs
    for (let i in this.state.stateDiffs) {
      let hostInstance = this.state.stateDiffs[i];
      let name = i + ' - B.states diff';
      let trace = {
        name: name,
        mode: 'markers',
        x: hostInstance.map(diff => diff.Timestamp * 1000),
        y: hostInstance.map(_ => name)
      }
      traces.push(trace);
    }

    // reports
    for (let i in this.state.reportTimeline) {
      let hostInstance = this.state.reportTimeline[i];
      let name = i + ' - C.reports';
      let trace = {
        name: name,
        mode: 'markers',
        x: hostInstance.map(document => document.Document * 1000),
        y: hostInstance.map(_ => name)
      }
      traces.push(trace);
    }

    // reports diffs
    for (let i in this.state.reportDiffs) {
      let hostInstance = this.state.reportDiffs[i];
      let name = i + ' - D.reports diff';
      let trace = {
        name: name,
        mode: 'markers',
        x: hostInstance.map(diff => diff.Timestamp * 1000),
        y: hostInstance.map(_ => name)
      }
      traces.push(trace);
    }

    // sort traces by name
    traces.sort(function(a, b) {
      if (a.name < b.name) {
        return -1;
      }
      if (a.name > b.name) {
        return 1;
      }
      return 0;
    });

    let selected = null;
    if (!this.state.selected) {
      selected = (
        <React.Fragment>
          No result
        </React.Fragment>
      );
    }
    // diff
    else if (this.state.selected.Changes != undefined) {
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
    else if (this.state.selected.Findings != undefined) {
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
    else if (this.state.selected.Users != undefined) {
      let time = new Date(this.state.selected.Timestamp * 1000).toLocaleString();
      let errors = [];
      for (let i in this.state.selected.Errors) {
        let error = this.state.selected.Errors[i];
        errors.push(<li key={i}>{error}</li>);
      }
      let users = [];
      for (let i in this.state.selected.Users) {
        let user = this.state.selected.Users[i];
        let passwordLastSet = new Date(user.PasswordLastSet * 1000).toLocaleString();
        users.push(<li key={i}>{user.ID} - {user.Name} - {user.AccountActive} - {user.AccountExpires} - {passwordLastSet} - {user.PasswordExpires}</li>);
      }
      let groups = [];
      for (let group in this.state.selected.Groups) {
        let members = this.state.selected.Groups[group];
        if (members.length === 0) {
          groups.push(<li key={group}>{group}</li>)
        }
        else {
          let membersStr = members.map(member => member.Name).join(', ');
          groups.push(<li key={group}>{group} - [{membersStr}]</li>);
        }
      }
      let software = [];
      for (let i in this.state.selected.Software) {
        let sw = this.state.selected.Software[i];
        software.push(<li key={i}>{sw.Name} - {sw.Version}</li>);
      }
      let processes = [];
      for (let i in this.state.selected.Processes) {
        let process = this.state.selected.Processes[i];
        processes.push(<li key={i}>{process.PID} - {process.User} - {process.CommandLine}</li>);
      }
      let conns = [];
      for (let i in this.state.selected.NetworkConnections) {
        let conn = this.state.selected.NetworkConnections[i];
        conns.push(<li key={i}>{conn.Protocol} - {conn.LocalAddress}:{conn.LocalPort} - {conn.RemoteAddress}:{conn.RemotePort} - {conn.State}</li>);
      }
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