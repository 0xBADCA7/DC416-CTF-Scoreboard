const fs = require('fs')
const express = require('express')

const app = express()


app.get('/', (req, res) => {
  const fileContent = fs.readFileSync('app/index.html')
  res.append('Content-Type', 'text/html')
  res.send(fileContent)
})

app.use(express.static('dist'))

app.listen(9001)