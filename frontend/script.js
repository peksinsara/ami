const webSocket = new WebSocket("ws://192.168.1.27:8081/ws");

webSocket.onmessage = function(event) {
    const status = JSON.parse(event.data);
    document.getElementById("numOnline").innerHTML = status.numOnline;
    document.getElementById("numTotal").innerHTML = status.numTotal;
    document.getElementById("numOffline").innerHTML = status.numOffline;
    document.getElementById("numActiveChannels").innerHTML = status.numActiveChannels;
    document.getElementById("numActiveCalls").innerHTML = status.numActiveCalls;
    document.getElementById("numCallsProcessed").innerHTML = status.numCallsProcessed;
    document.getElementById("lastUpdate").innerHTML = status.lastUpdate;
}