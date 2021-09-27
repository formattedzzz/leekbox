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

// etcdctl --endpoints="$ENDPOINTS" --write-out="json" get project
// {
//   header: {
//     cluster_id: 15851974733440446356,
//     member_id: 4767751730662354927,
//     revision: 7,
//     raft_term: 15
//   },
//   kvs: [
//     {
//       key: 'cHJvamVjdA==',
//       create_revision: 5,
//       mod_revision: 7,
//       version: 3,
//       value: 'bGVla2JveA=='
//     }
//   ],
//   count: 1
// }
