import { Component } from 'react';

class ScenarioChecks extends Component {
    constructor(props) {
        super(props);
        this.state = {
            scenarioID: props.scenarioID,
            checkMap: this.setupCheckMap(props)
        }

        this.handleAddCheck = this.handleAddCheck.bind(this);
        this.handleAddCheckArg = this.handleAddCheckArg.bind(this);
        this.handleUpdateCheck = this.handleUpdateCheck.bind(this);
        this.handleUpdateCheckArg = this.handleUpdateCheckArg.bind(this);
    }

    componentDidUpdate(prevProps) {
        if (this.props.scenarioID !== prevProps.scenarioID) {
            this.setState({
                scenarioID: this.props.scenarioID,
                checkMap: this.setupCheckMap(this.props)
            });
        }
    }
    
    setupCheckMap(props) {
        if (props.checkMap) {
            return props.checkMap;
        }
        return {};
    }

    handleAddCheck() {
        let checkMap = {
            ...this.state.checkMap
        }
        let hostname = "new hostname" + (Object.keys(checkMap).length + 1);
        checkMap[hostname] = [{
            Type: 'EXEC',
            Command: null,
            Args: []
        }]
        this.setState({
            checkMap: checkMap
        });
    }

    handleAddCheckArg(hostname, i) {
        console.log("add check arg " + hostname + " " + i);
    }

    handleUpdateCheck(hostname, i, name, event) {
        let value = event.target.value;
        let checkMap = {
            ...this.state.checkMap
        }
        checkMap[hostname][i][name] = value;
        this.setState({
            checkMap: checkMap
        });
    }

    handleUpdateCheckArg(hostname, i, j, event) {
        let value = event.target.value;
        let checkMap = {
            ...this.state.checkMap
        }
        checkMap[hostname][i]['Args'][j] = value;
        this.setState({
            checkMap: checkMap
        });
    }

    render() {
        let actionExecOptions = [
            <option>A</option>,
            <option>EXEC</option>,
            <option>FILE_EXISTS</option>
        ]

        let checkList = [];
        for (let hostname in this.state.checkMap) {
            let checks = this.state.checkMap[hostname];
            checks.forEach((check, i) => {
                let args = [];
                if (check.Args) {
                    check.Args.forEach((arg, j) => {
                        args.push(<input onChange={event => this.handleUpdateCheckArg(hostname, i, j, event)} value={arg}></input>);
                    });
                }
                checkList.push(
                    <li key={`${hostname}${i}`}>
                        <input value={hostname} />
                        <br />
                        <select onChange={event => this.handleUpdateCheck(hostname, i, "Type", event)} value={check.Type}>{actionExecOptions}</select>
                        <input onChange={event => this.handleUpdateCheck(hostname, i, "Command", event)} value={check.Command} />
                        {args}
                        <button type="button" onClick={() => this.handleAddCheckArg(hostname, i)}>+</button>
                    </li>
                );
            });
        }

        return (
            <ul>
                {checkList}
                <li key="new"><button type="button" onClick={this.handleAddCheck}>+</button></li>
            </ul>
        );
    }
}

export default ScenarioChecks;
