const webSocket = new WebSocket('ws://192.168.1.8:8081/status');

webSocket.onmessage = function(event) {
  const data = JSON.parse(event.data);
 // console.log('Received data:', data);

  if (data.status) {
    const activePeers = data.status.active_peers;
   // document.getElementById('numOnline').innerHTML = activePeers;

    const inactivePeers = data.status.inactive_peers;
    document.getElementById('numOffline').innerHTML = inactivePeers;
  } else {
    console.log('No "status" data found in received message');
  }

  if (data.calls) {
    const activeCalls = data.calls.active_calls;
    document.getElementById('numActiveCalls').innerHTML = activeCalls;
  } else {
    console.log('No "calls" data found in received message');
  }
};
