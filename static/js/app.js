(function () {
  'use strict';

  function connect() {
    var port = location.port === '' ? '' : ':' + location.port;
    var wsUrl = 'ws://' + document.domain + port + '/echo';
    if (/https/.test(location.href)) {
      wsUrl = 'wss://' + document.domain + '/echo';
    }

    var uuid;
    // var echoWs = new WebSocket(wsUrl);
    var echoWs = new ReconnectingWebSocket(wsUrl);

    var messageFormElem = document.getElementById('message-form');
    var messagesUlElem = document.getElementById('messages');
    var messageInputElem = document.getElementById('message');
    var uuidElem = document.getElementById('uuid');

    echoWs.onopen = function () {
      // console.log('open ws');
    }

    echoWs.onmessage = function (event) {
      var received = JSON.parse((event.data))


      var liElem = document.createElement('li');
      var content = document.createTextNode(received.uuid + ': ' + received.data);
      if (received.type === 'server') {
        if (!uuid) {
          var uuidTextNode = document.createTextNode(received.uuid);
          uuidElem.appendChild(uuidTextNode)
          content = document.createTextNode(received.type + ': ' + received.data);
        }
        uuid = received.uuid;
      }
      liElem.appendChild(content);
      messagesUlElem.appendChild(liElem);
    };

    echoWs.onclose = function (event) {
      // console.log('Socket is closed. Reconnect will be attempted in 1 second.');
      // setTimeout(function () {
        //   connect();
        // }, 1000);
      };

    messageFormElem.onsubmit = function (event) {
      event.preventDefault();

      echoWs.send(JSON.stringify({
        uuid: uuid,
        type: 'client',
        data: messageInputElem.value
      }));
      messageInputElem.value = '';
    };
  }

  connect();
}());

