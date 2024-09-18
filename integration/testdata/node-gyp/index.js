const express = require('express')
const app = express()
const port = process.env.PORT || 8080

const features = require('cpu-features')();

app.get('/', (req, res) => {
    res.send(features)
})

app.listen(port, () => {
    console.log(`Example app listening on port ${port}`)
})
