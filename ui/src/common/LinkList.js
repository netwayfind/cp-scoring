import { Component } from "react";
import { Link } from "react-router-dom";
import { withRouter } from "react-router-dom/cjs/react-router-dom.min";

class LinkList extends Component {
  render() {
    let items = [];
    this.props.items.forEach((item, i) => {
      let classes = ["nav-button"];
      if (this.props.currentID === item.ID) {
        classes.push("nav-button-selected");
      }
      items.push(
        <li key={i}>
          <Link
            className={classes.join(" ")}
            to={`${this.props.path}/${item.ID}`}
          >
            [{item.ID}] {item[this.props.label]}
          </Link>
        </li>
      );
    });

    return <ul>{items}</ul>;
  }
}

export default withRouter(LinkList);
