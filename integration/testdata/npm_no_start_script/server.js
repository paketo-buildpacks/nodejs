const http = require('http');
const leftpad = require('leftpad');

const port = process.env.PORT || 8080;

const server = http.createServer((request, response) => {
  response.end(JSON.stringify(process.env))
});

server.listen(port, (err) => {
  if (err) {
    return console.log('something bad happened', err);
  }

  console.log(`NOT vendored server is listening on ${port}`);
});
