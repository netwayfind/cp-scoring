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
      var scenario = "0";
      var teamKey = "key";
      var hostname = "hostname";
      var query = window.location.search.substring(1);
      var params = query.split("&");
      for (var i = 0; i < params.length; i++) {
        var param = params[i].split("=");
        if (param.length != 2) {
          continue;
        }
        if (param[0] === "scenario") {
          scenario = param[1];
        } else if (param[0] === "team_key") {
          teamKey = param[1];
        } else if (param[0] === "hostname") {
          hostname = param[1];
        }
      }
      return React.createElement(
        "div",
        { className: "App" },
        React.createElement(ScoreTimeline, { scenarioID: scenario, teamKey: teamKey, hostname: hostname })
      );
    }
  }]);

  return App;
}(React.Component);

var ScoreTimeline = function (_React$Component2) {
  _inherits(ScoreTimeline, _React$Component2);

  function ScoreTimeline() {
    _classCallCheck(this, ScoreTimeline);

    var _this2 = _possibleConstructorReturn(this, (ScoreTimeline.__proto__ || Object.getPrototypeOf(ScoreTimeline)).call(this));

    _this2.state = {
      timestamps: [],
      scores: [],
      report: {}
    };
    return _this2;
  }

  _createClass(ScoreTimeline, [{
    key: "populateScores",
    value: function populateScores() {
      var scenarioID = this.props.scenarioID;
      var teamKey = this.props.teamKey;
      var hostname = this.props.hostname;
      var url = "/reports/scenario/" + scenarioID + "/timeline?team_key=" + teamKey + "&hostname=" + hostname;

      fetch(url).then(function (response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        return response.json();
      }).then(function (data) {
        if (data) {
          // should only be one match
          this.setState({
            scores: data.Scores,
            // timestamps is seconds, need milliseconds
            timestamps: data.Timestamps.map(function (timestamp) {
              return timestamp * 1000;
            })
          });
        }
      }.bind(this));
    }
  }, {
    key: "populateReport",
    value: function populateReport() {
      var scenarioID = this.props.scenarioID;
      var teamKey = this.props.teamKey;
      var hostname = this.props.hostname;
      var url = '/reports/scenario/' + scenarioID + '?team_key=' + teamKey + '&hostname=' + hostname;

      fetch(url).then(function (response) {
        if (response.status >= 400) {
          throw new Error("Bad response from server");
        }
        return response.json();
      }).then(function (data) {
        this.setState({
          report: data
        });
      }.bind(this));
    }
  }, {
    key: "componentDidMount",
    value: function componentDidMount() {
      this.populateScores();
      this.populateReport();
    }
  }, {
    key: "render",
    value: function render() {
      var data = [{
        x: this.state.timestamps,
        y: this.state.scores,
        type: 'scatter',
        mode: 'lines+markers'
      }];

      var layout = {
        xaxis: {
          type: 'date'
        },
        yaxis: {
          fixedrange: true
        }
      };

      var config = {
        displayModeBar: false
      };

      var rows = [];
      if (this.state.report) {
        for (var i in this.state.report.Findings) {
          var finding = this.state.report.Findings[i];
          if (!finding.Hidden) {
            rows.push(React.createElement(
              "li",
              { key: i },
              finding.Value,
              " - ",
              finding.Message
            ));
          } else {
            rows.push(React.createElement(
              "li",
              { key: i },
              "?"
            ));
          }
        }
      }

      return React.createElement(
        "div",
        { className: "ScoreTimeline" },
        React.createElement(
          "strong",
          null,
          "Score Timeline"
        ),
        React.createElement("p", null),
        React.createElement(Plot, { data: data, layout: layout, config: config }),
        React.createElement(
          "ul",
          null,
          rows
        )
      );
    }
  }]);

  return ScoreTimeline;
}(React.Component);

ReactDOM.render(React.createElement(App, null), document.getElementById('app'));