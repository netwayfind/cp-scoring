'use strict';

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _possibleConstructorReturn(self, call) { if (!self) { throw new ReferenceError("this hasn't been initialised - super() hasn't been called"); } return call && (typeof call === "object" || typeof call === "function") ? call : self; }

function _inherits(subClass, superClass) { if (typeof superClass !== "function" && superClass !== null) { throw new TypeError("Super expression must either be null or a function, not " + typeof superClass); } subClass.prototype = Object.create(superClass && superClass.prototype, { constructor: { value: subClass, enumerable: false, writable: true, configurable: true } }); if (superClass) Object.setPrototypeOf ? Object.setPrototypeOf(subClass, superClass) : subClass.__proto__ = superClass; }

var Plot = createPlotlyComponent(Plotly);

var App = function (_React$Component) {
  _inherits(App, _React$Component);

  function App() {
    _classCallCheck(this, App);

    return _possibleConstructorReturn(this, (App.__proto__ || Object.getPrototypeOf(App)).apply(this, arguments));
  }

  _createClass(App, [{
    key: "render",
    value: function render() {
      return React.createElement(
        "div",
        { className: "App" },
        React.createElement(Scoreboard, { scenarioID: "1" })
      );
    }
  }]);

  return App;
}(React.Component);

var Scoreboard = function (_React$Component2) {
  _inherits(Scoreboard, _React$Component2);

  function Scoreboard() {
    _classCallCheck(this, Scoreboard);

    var _this2 = _possibleConstructorReturn(this, (Scoreboard.__proto__ || Object.getPrototypeOf(Scoreboard)).call(this));

    _this2.state = {
      scores: []
    };
    return _this2;
  }

  _createClass(Scoreboard, [{
    key: "populateScores",
    value: function populateScores() {
      var id = this.props.scenarioID;
      var url = '/scenarios/' + id + '/scores';

      fetch(url).then(function (response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        return response.json();
      }).then(function (data) {
        this.setState({ scores: data });
      }.bind(this));
    }
  }, {
    key: "componentDidMount",
    value: function componentDidMount() {
      this.populateScores();
    }
  }, {
    key: "render",
    value: function render() {
      var body = [];
      for (var i in this.state.scores) {
        var entry = this.state.scores[i];
        body.push(React.createElement(
          "tr",
          { key: i },
          React.createElement(
            "td",
            null,
            entry.TeamName
          ),
          React.createElement(
            "td",
            null,
            entry.Score
          )
        ));
      }

      return React.createElement(
        "div",
        { className: "Scoreboard" },
        React.createElement(
          "strong",
          null,
          "Scoreboard"
        ),
        React.createElement("p", null),
        React.createElement(
          "table",
          null,
          React.createElement(
            "thead",
            null,
            React.createElement(
              "tr",
              null,
              React.createElement(
                "th",
                null,
                "Team"
              ),
              React.createElement(
                "th",
                null,
                "Score"
              )
            )
          ),
          React.createElement(
            "tbody",
            null,
            body
          )
        )
      );
    }
  }]);

  return Scoreboard;
}(React.Component);

ReactDOM.render(React.createElement(App, null), document.getElementById('app'));