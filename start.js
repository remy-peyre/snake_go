window.onload = function() {
    var ws = new WebSocket('wss://golang-game-skipcat.c9users.io:8081');
    ws.onopen = function(event) {
        console.log('Connection with ws OK');  
    };
    ws.onerror = function(event) {
        console.log("ERR: " + event.data);
    };
    ws.onmessage = function(event) {
        var msg = JSON.parse(event.data);
        if (msg.kind == 'init' ) {
            console.log(msg);
        }
        else if (msg.kind == 'update') {
            console.log(msg);
        }
        else {
            if (msg == 'p1 selected!') {
                document.querySelector('#player1').disabled = true;
            }
            if (msg == 'p2 selected!') {
                document.querySelector('#player2').disabled = true;
            }
            console.log(msg);
        }
    };
    
    var btnStart = document.querySelector('#btn-start');
    //var players = document.querySelectorAll('.player');
    //var playerList = [];
    var player1 = document.querySelector('#player1');
    var player2 = document.querySelector('#player2');

    var add_player = {
        kind : "add_player"
    };
    
    player1.onclick = function() {
        var nbPlayer = {
            kind : "connect",
            slot : parseInt(player1.value, 10)
        };
        ws.send(JSON.stringify(nbPlayer));
        console.log(nbPlayer);
    };
    
    player2.onclick = function() {
        var nbPlayer = {
            kind : "connect",
            slot : parseInt(player2.value, 10)
        };
        ws.send(JSON.stringify(nbPlayer));
        console.log(nbPlayer);
    };
    
    /*
    for (var i = 0; i < players.length; i ++) {
        players[i].onclick = function() {
            playerList.push(this.value);
            if (playerList.length == 1) {
                var nbPlayer = {
                    kind : "connect",
                    slot : 1,
                };
            }
            else {
                var nbPlayer = {
                    kind : "connect",
                    slot : 2,
                };
            }
            ws.send(JSON.stringify(nbPlayer));
            console.log(nbPlayer);
        };
    }
    */
    
    btnStart.onclick = function() {
        if (playerList.length >= 1) {
            ws.send(JSON.stringify(add_player));
            window.location.replace('canvas.html');
        }
    };
};



