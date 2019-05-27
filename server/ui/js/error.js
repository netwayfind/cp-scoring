'use strict';

class Error extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    if (this.props.message === null) {
      return null;
    }

    return React.createElement("div", {
      class: "error"
    }, this.props.message);
  }

}