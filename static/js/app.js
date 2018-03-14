(function () {
  'use strict';

  function connect() {
    var port = location.port === '' ? '' : ':' + location.port;
    var wsUrl = 'ws://' + document.domain + port + '/echows';
    if (/https/.test(location.href)) {
      wsUrl = 'wss://' + document.domain + '/echo';
    }

    var uuid;
    var connected = false;
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

      if (received.type === 'server:hello') {
        uuid = received.uuid;
        uuidElem.innerText = received.uuid;
        if (!connected) {
          var liElem = document.createElement('li');
          content = document.createTextNode(received.type + ': ' + received.data);
          liElem.appendChild(content);
          messagesUlElem.appendChild(liElem);
          connected = true;
        }
      }

      if (received.type === 'client') {
        var liElem = document.createElement('li');
        var content = document.createTextNode(received.uuid + ': ' + received.data);
        liElem.appendChild(content);
        messagesUlElem.appendChild(liElem);
      }
    };

    echoWs.onclose = function (event) {
      console.log('Socket is closed. Reconnect will be attempted in 1 second.');
      // setTimeout(function () {
      //   connect();
      // }, 1000);
    };

    messageFormElem.onsubmit = function (event) {
      event.preventDefault();
      var data = {
        uuid: uuid,
        type: 'client',
        data: messageInputElem.value
      };
      echoWs.send(JSON.stringify(data));
      messageInputElem.value = '';
    };
  }

  connect();
}());

