
var http = require('http');

var fs = require("fs");

var handleRequest = function(request:any, response:any) {
    console.log('Received request for URL: ' + request.url);
    response.writeHead(200);
    response.end('Hello World!');
  };
  var www = http.createServer(handleRequest);
  www.listen(8080);


 