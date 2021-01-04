import '../App.css';

import { Component } from 'react';
import { withRouter } from 'react-router-dom';

class ScenarioHost extends Component {
    constructor(props) {
        super(props);
        this.state = {
            answers: props.answers,
            checks: props.checks,
            config: props.config,
            hostname: props.hostname,
        }

        this.handleAnswerUpdate = this.handleAnswerUpdate.bind(this);
        this.handleCheckAdd = this.handleCheckAdd.bind(this);
        this.handleCheckDelete = this.handleCheckDelete.bind(this);
        this.handleCheckUpdate = this.handleCheckUpdate.bind(this);
        this.handleCheckArgAdd = this.handleCheckArgAdd.bind(this);
        this.handleCheckArgDelete = this.handleCheckArgDelete.bind(this);
        this.handleCheckArgUpdate = this.handleCheckArgUpdate.bind(this);
        this.handleConfigAdd = this.handleConfigAdd.bind(this);
        this.handleSave = this.handleSave.bind(this);
    }

    componentDidUpdate(prevProps) {
        if (this.props.hostname !== prevProps.hostname) {
            this.setState({
                answers: this.props.answers,
                checks: this.props.checks,
                config: this.props.config,
                hostname: this.props.hostname,
            });
        }
    }
    
    handleAnswerUpdate(i, event) {
        let name = event.target.name;
        let value = event.target.value;
        let answers = [
            ...this.state.answers
        ]
        if (event.target.type === "number") {
            value = Number(value);
        }
        answers[i][name] = value;
        this.setState({
            answers: answers
        });
    }

    handleCheckAdd() {
        let answers = [
            ...this.state.answers
        ]
        let checks = [
            ...this.state.checks
        ]
        answers.push({
            Type: '',
            Value: ''
        });
        checks.push({
            Type: 'EXEC',
            Command: '',
            Args: []
        });
        this.setState({
            answers: answers,
            checks: checks
        });
    }

    handleCheckDelete(i) {
        let answers = [
            ...this.state.answers
        ]
        let checks = [
            ...this.state.checks
        ]
        answers.splice(i, 1);
        checks.splice(i, 1);
        this.setState({
            answers: answers,
            checks: checks
        });
    }

    handleCheckUpdate(i, event) {
        let name = event.target.name;
        let value = event.target.value;
        let checks = [
            ...this.state.checks
        ]
        checks[i][name] = value;
        this.setState({
            checks: checks
        });
    }

    handleCheckArgAdd(i) {
        let checks = [
            ...this.state.checks
        ]
        checks[i]['Args'].push('');
        this.setState({
            checks: checks
        });
    }

    handleCheckArgDelete(i, j) {
        let checks = [
            ...this.state.checks
        ]
        checks[i]['Args'].splice(j, 1);
        this.setState({
            checks: checks
        });
    }

    handleCheckArgUpdate(i, j, event) {
        let value = event.target.value;
        let checks = [
            ...this.state.checks
        ]
        checks[i]['Args'][j] = value;
        this.setState({
            checks: checks
        });
    }

    handleConfigAdd() {

    }

    handleSave(event) {
        if (event !== null) {
            event.preventDefault();
        }
        this.props.parentCallback(this.state.checks, this.state.answers, this.state.config);
    }

    render() {
        let actionOptions = [
            <option key="1">A</option>,
            <option key="2">EXEC</option>,
            <option key="3">FILE_EXISTS</option>
        ]
        let operatorOptions = [
            <option key="1" value='' />,
            <option key="2">EQUAL</option>,
            <option key="3">NOT_EQUAL</option>,
            <option key="4">NIL</option>,
            <option key="5">NOT_NIL</option>
        ]

        let checkList = [];
        let checks = this.state.checks;
        checks.forEach((check, i) => {
            let args = [];
            if (check.Args) {
                check.Args.forEach((arg, j) => {
                    args.push(
                        <li key={j}>
                            <input onChange={event => this.handleCheckArgUpdate(i, j, event)} value={arg}></input>
                            <button type="button" onClick={() => this.handleCheckArgDelete(i, j)}>-</button>
                        </li>
                    );
                });
            }
            args.push(
                <li key="arg_add"><button type="button" onClick={() => this.handleCheckArgAdd(i)}>Add Arg</button></li>
            );
            let answer = this.state.answers[i];
            checkList.push(
                <li key={i}>
                    <details>
                        <summary>Type: {check.Type}, Command: {check.Command}, Args: [{ check.Args.join(" ") || ""}]</summary>
                        <button type="button" onClick={() => this.handleCheckDelete(i)}>Delete Check</button>
                        <p />
                        <label htmlFor="Type">Type</label>
                        <select name="Type" onChange={event => this.handleCheckUpdate(i, event)} value={check.Type}>{actionOptions}</select>
                        <br />
                        <label htmlFor="Command">Command</label>
                        <input name="Command" onChange={event => this.handleCheckUpdate(i, event)} value={check.Command} />
                        <br />
                        <label htmlFor="Args">Args</label>
                        <ul>{args}</ul>
                        <label htmlFor="Answer">Answer</label>
                        <select name="Operator" onChange={event => this.handleAnswerUpdate(i, event)} value={answer.Operator}>{operatorOptions}</select>
                        <input name="Value" onChange={event => this.handleAnswerUpdate(i, event)} value={answer.Value} />
                        <input name="Description" onChange={event => this.handleAnswerUpdate(i, event)} value={answer.Description} />
                        <input name="Points" onChange={event => this.handleAnswerUpdate(i, event)} value={answer.Points} type="number" steps="1" />
                    </details>
                </li>
            );
        });
        checkList.push(
            <li key="check_add">
                <button type="button" onClick={this.handleCheckAdd}>Add Check</button>
            </li>
        );

        let configList = [];
        configList.push(
            <li key="config_add">
                <button type="button" onClick={this.handleConfigAdd}>Add Config</button>
            </li>
        );

        return (
            <form onSubmit={this.handleSave}>
                <p>Checks</p>
                <ul>{checkList}</ul>
                <p>Config</p>
                <ul>{configList}</ul>
                <button type="submit">Save Host</button>
            </form>
        );
    }
}

export default withRouter(ScenarioHost);
