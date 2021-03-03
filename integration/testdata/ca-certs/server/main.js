const fs = require('fs');
const https = require('https');
const tls = require('tls');
const port = process.env.PORT || 8080;

const options = {
  cert: fs.readFileSync('cert.pem'),
  key: fs.readFileSync('key.pem'),
  requestCert: true,
  rejectUnauthorized: false,
};

const handler = (req, res) => {
  if (!req.client.authorized) {
    res.writeHead(401);
    return res.end('Invalid client certificate authentication. ' + req.client.authorizationError);
  }

  res.writeHead(200);
  res.end('Hello, world!');
};

const server = https.createServer(options, handler);

process.once('SIGINT', function (code) {
  console.log('echo from SIGINT handler');
  server.close();
});

server.listen(port);
