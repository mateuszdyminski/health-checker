/* jshint devel:true */

(function(checker) {

    var loc = window.location, new_uri;
    if (loc.protocol === "https:") {
        new_uri = "wss:";
    } else {
        new_uri = "ws:";
    }
    new_uri += "//" + loc.hostname + ":8090/wsapi/ws";

    var socket = new WebSocket(new_uri);
    var id = -1;

    var data = {
    labels: [new Date().toTimeString().split(' ')[0]],
    datasets: [
        {
            label: "Responses",
            fillColor: "rgba(220,220,220,0.2)",
            strokeColor: "rgba(220,220,220,1)",
            pointColor: "rgba(220,220,220,1)",
            pointStrokeColor: "#fff",
            pointHighlightFill: "#fff",
            pointHighlightStroke: "rgba(220,220,220,1)",
            data: [0]
        }
    ]
    };

    var ctx = document.getElementById("status").getContext("2d");
    var myLineChart = new Chart(ctx).Line(data);

    socket.onopen = function(event) {
        console.log('Server connection open.');
    };

    socket.onmessage = function(msg) {
        var message = JSON.parse(msg.data);
        myLineChart.addData([message.Duration], new Date().toTimeString().split(' ')[0])
    }

    socket.onclose = function() {
        console.log('Server connection closed.');
        socket = undefined;
    }

    socket.onerror = function() {
        console.log('Server connection failure.');
        socket = undefined;
    }

})(this.chat = {});