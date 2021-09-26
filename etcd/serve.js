const http = require('http')

const app = http.createServer((req, res) => {
  console.log(req.url, req.query)
  res.writeHead(200, { 'content-type': 'application/json' })
  res.end('{"data":"hello world!"}')
})
const port = 5000
app.listen(port, () => {
  console.log('app running at :' + port)
})
