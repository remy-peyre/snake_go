// TODO
// display several apples ?
// more generic for snake move (not SnakeList[0])
// moveAllSnakes ?
// snake move automatic

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
        
        if (msg.kind == "init") {
            console.log(msg);
            
        }
        else if (msg.kind == "update") {
            console.log(msg);
            var snake1 = msg.snakes[0];
            var snake2 = msg.snakes[1];
            console.log(snake1, snake2);
        
            createCanvas();
            displaySnake(snake1);
            displaySnake(snake2);
            
        }
        else {
            console.log(msg);
        }
    };
    
    var canvas = document.querySelector('canvas');
    var context = canvas.getContext('2d');
    var nbPixel = 500;
    var nbCellsByLine = 20;
    var cellsize = nbPixel / nbCellsByLine;

    var SnakeList = [];
    var AppleList = [];

/*********************** CANVAS *************************/

    function createCanvas() {
        var i = 0;
        for (var x = 0; x < canvas.width; x += cellsize) {
            for (var y = 0; y < canvas.height; y += cellsize)  {
                if (i % 2 == 0) {
                    context.fillStyle = 'ivory';
                }
                else {
                    context.fillStyle = 'black';
                }
                context.beginPath();
                context.fillRect(x, y, cellsize, cellsize); // x and y coordonnates, 50 length and width
                i ++;
            }
            i ++;
        }
    }

/*********************** SNAKES *************************/

    function displaySnake(snake) {
        for (var i = 0; i < snake.body.length; i ++) {
            context.beginPath();
            context.fillStyle = snake.color;
            context.fillRect(snake.body[i].x * cellsize, snake.body[i].y * cellsize, cellsize, cellsize);
            context.strokeStyle = 'black';
            context.strokeRect(snake.body[i].x * cellsize, snake.body[i].y * cellsize, cellsize, cellsize);
        }
        // display snake head
        context.fillStyle = 'blue';
        context.fillRect(snake.body[0].x * cellsize, snake.body[0].y * cellsize, cellsize, cellsize);
        context.strokeStyle = 'black';
        context.strokeRect(snake.body[0].x * cellsize, snake.body[0].y * cellsize, cellsize, cellsize);
    }
    
    
    /*
    function displayAllSnakes(SnakeList) {
        for (var snake in SnakeList) {
            displaySnake(SnakeList[snake]);
        }
    }
    */
    
    document.onkeydown = function(event) {
        var move = {
            kind : 'move',
            key : ''
        };
        
        if (event.keyCode === 37 || event.keyCode === 81) { // left arrow || Q
            move.key = 'left';
            ws.send(JSON.stringify(move));
        }
        if (event.keyCode === 38 || event.keyCode === 90) { // up arrow || Z
            move.key = 'up';
            ws.send(JSON.stringify(move));
        }
        if (event.keyCode === 39 || event.keyCode === 68) { // right arrow || D
            move.key = 'right';
            ws.send(JSON.stringify(move));
        }
        if (event.keyCode === 40 || event.keyCode === 83) { // down arrow || S
            move.key = 'down';
            ws.send(JSON.stringify(move));
        }
    };
    
/*********************** MAIN PROGRAM *************************/
    createCanvas();
    //displayAllSnakes(SnakeList);
};
