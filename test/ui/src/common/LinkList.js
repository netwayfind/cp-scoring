import { Component } from 'react';
import { Link } from 'react-router-dom';
import { withRouter } from "react-router-dom/cjs/react-router-dom.min";

class LinkList extends Component {
    render() {
        let items = [];
        this.props.items.forEach((item, i) => {
            items.push(
                <li key={i}><Link to={`${this.props.path}/${item.ID}`}>[{item.ID}] {item.Name}</Link></li>
            );
        });

        return (
            <ul>{items}</ul>
        );
    }
}

export default withRouter(LinkList);
