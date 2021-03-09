const fs = require('fs');
const https = require('https');
const leftpad = require('leftpad');
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
  res.end(JSON.stringify(process.env));
};

const server = https.createServer(options, handler);

server.listen(port, (err) => {
  if (err) {
    return console.log('something bad happened', err);
  }
  console.log(`server is listening on ${port}`);
});
