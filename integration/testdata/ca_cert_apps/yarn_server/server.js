const fs = require('fs');
const https = require('https');
const leftpad = require('leftpad');
const tls = require('tls');
const port = process.env.PORT || 8080;

const options = {
  cert: fs.readFileSync('cert.pem'),
  key: fs.readFileSync('key.pem'),
  requestCert: true,
  rejectUnauthorized: false,
};

const requestHandler = (request, response) => {
  response.end("Hello, World!")
}

const server = https.createServer(options, requestHandler)

server.listen(port, (err) => {
  if (err) {
    return console.log('something bad happened', err)
  }

  console.log(`server is listening on ${port}`)
})
