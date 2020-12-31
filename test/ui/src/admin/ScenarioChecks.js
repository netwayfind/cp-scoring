import { Component } from 'react';

class ScenarioChecks extends Component {
    constructor(props) {
        super(props);
        this.state = {
            scenarioID: props.scenarioID,
            checkMap: props.checkMap,
            currentHostname: '',
            newHostname: ''
        }

        this.handleCheckAdd = this.handleCheckAdd.bind(this);
        this.handleCheckDelete = this.handleCheckDelete.bind(this);
        this.handleCheckUpdate = this.handleCheckUpdate.bind(this);
        this.handleCheckArgAdd = this.handleCheckArgAdd.bind(this);
        this.handleCheckArgDelete = this.handleCheckArgDelete.bind(this);
        this.handleCheckArgUpdate = this.handleCheckArgUpdate.bind(this);
        this.handleHostnameAdd = this.handleHostnameAdd.bind(this);
        this.handleHostnameDelete = this.handleHostnameDelete.bind(this);
        this.handleHostnameSelect = this.handleHostnameSelect.bind(this);
        this.handleNewHostnameUpdate = this.handleNewHostnameUpdate.bind(this);
        this.handleSave = this.handleSave.bind(this);
    }

    componentDidUpdate(prevProps) {
        if (this.props.scenarioID !== prevProps.scenarioID) {
            this.setState({
                scenarioID: this.props.scenarioID,
                checkMap: this.props.checkMap,
                currentHostname: '',
                newHostname: ''
            });
        }
    }

    handleCheckAdd(hostname) {
        let checkMap = {
            ...this.state.checkMap
        }
        checkMap[hostname].push({
            Type: 'EXEC',
            Command: '',
            Args: []
        });
        this.setState({
            checkMap: checkMap
        });
    }

    handleCheckDelete(hostname, i) {
        let checkMap = {
            ...this.state.checkMap
        }
        checkMap[hostname].splice(i, 1);
        this.setState({
            checkMap: checkMap
        });
    }

    handleCheckUpdate(hostname, i, name, event) {
        let value = event.target.value;
        let checkMap = {
            ...this.state.checkMap
        }
        checkMap[hostname][i][name] = value;
        this.setState({
            checkMap: checkMap
        });
    }

    handleCheckArgAdd(hostname, i) {
        let checkMap = {
            ...this.state.checkMap
        }
        checkMap[hostname][i]['Args'].push('');
        this.setState({
            checkMap: checkMap
        });
    }

    handleCheckArgDelete(hostname, i, j) {
        let checkMap = {
            ...this.state.checkMap
        }
        checkMap[hostname][i]['Args'].splice(j, 1);
        this.setState({
            checkMap: checkMap
        });
    }

    handleCheckArgUpdate(hostname, i, j, event) {
        let value = event.target.value;
        let checkMap = {
            ...this.state.checkMap
        }
        checkMap[hostname][i]['Args'][j] = value;
        this.setState({
            checkMap: checkMap
        });
    }

    handleHostnameAdd() {
        let hostname = this.state.newHostname;
        if (!hostname) {
            return;
        }
        let checkMap = {
            ...this.state.checkMap
        }
        checkMap[hostname] = [{
            Type: 'EXEC',
            Command: '',
            Args: []
        }]
        this.setState({
            checkMap: checkMap,
            currentHostname: hostname,
            newHostname: ''
        });
    }

    handleHostnameDelete() {
        let hostname = this.state.currentHostname;
        if (!hostname) {
            return;
        }
        let checkMap = {
            ...this.state.checkMap
        }
        delete checkMap[hostname];
        this.setState({
            checkMap: checkMap,
            currentHostname: ''
        });
    }

    handleHostnameSelect(event) {
        let value = event.target.value;
        this.setState({
            currentHostname: value
        });
    }

    handleNewHostnameUpdate(event) {
        let newHostname = event.target.value;
        this.setState({
            newHostname: newHostname
        });
    }

    handleSave(event) {
        if (event !== null) {
            event.preventDefault();
        }
        this.props.parentCallback(this.state.checkMap);
    }

    render() {
        let actionExecOptions = [
            <option>A</option>,
            <option>EXEC</option>,
            <option>FILE_EXISTS</option>
        ]

        let hostnameList = [];
        hostnameList.push(<option value='' />);
        for (let hostname in this.state.checkMap) {
            hostnameList.push(<option>{hostname}</option>);
        }
        let checkList = [];
        if (this.state.currentHostname) {
            let hostname = this.state.currentHostname;
            let checks = this.state.checkMap[hostname];
            if (checks) {
                checks.forEach((check, i) => {
                    let args = [];
                    if (check.Args) {
                        check.Args.forEach((arg, j) => {
                            args.push(
                                <li key={j}>
                                    <input onChange={event => this.handleCheckArgUpdate(hostname, i, j, event)} value={arg}></input>
                                    <button type="button" onClick={() => this.handleCheckArgDelete(hostname, i, j)}>-</button>
                                </li>
                            );
                        });
                    }
                    checkList.push(
                        <li key={i}>
                            <details>
                                <summary>Type: {check.Type}, Command: {check.Command}, Args: [{ check.Args.join(" ") || ""}]</summary>
                                <button type="button" onClick={() => this.handleCheckDelete(hostname, i)}>Delete Check</button>
                                <p />
                                <label htmlFor="Type">Type</label>
                                <select onChange={event => this.handleCheckUpdate(hostname, i, "Type", event)} value={check.Type}>{actionExecOptions}</select>
                                <br />
                                <label htmlFor="Command">Command</label>
                                <input onChange={event => this.handleCheckUpdate(hostname, i, "Command", event)} value={check.Command} />
                                <br />
                                <label htmlFor="Args">Args</label>
                                <ul>
                                    {args}
                                    <li key="new_arg"><button type="button" onClick={() => this.handleCheckArgAdd(hostname, i)}>Add Arg</button></li>
                                </ul>
                            </details>
                        </li>
                    );
                });
            }
            checkList.push(
                <li key="new_check">
                    <button type="button" onClick={() => this.handleCheckAdd(hostname)}>Add Check</button>
                </li>
            );
        }

        return (
            <form onSubmit={this.handleSave}>
                <ul>
                    <input onChange={this.handleNewHostnameUpdate} value={this.state.newHostname} />
                    <button type="button" onClick={this.handleHostnameAdd}>Add Hostname</button>
                    <p />
                    <select onChange={this.handleHostnameSelect} value={this.state.currentHostname}>{hostnameList}</select>
                    <button type="button" disabled={!this.state.currentHostname} onClick={this.handleHostnameDelete}>Delete Hostname</button>
                    <p />
                    <ul>{checkList}</ul>
                </ul>
                <button type="submit" disabled={!this.state.currentHostname}>Save Checks</button>
            </form>
        );
    }
}

export default ScenarioChecks;
