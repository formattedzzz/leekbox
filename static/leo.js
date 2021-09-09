console.log('leojs inited')
window.addEventListener('load', () => {
  console.log('loaded')
})

ws = new WebSocket('ws://localhost:7003/api/stream?room_id=1')
ws.addEventListener('message', ev => console.log('message:', ev.data))
ws.addEventListener('close', console.log)
ws.addEventListener('open', ev => {
  console.log('open', ev)
  ws.send(JSON.stringify({ type: 'PING', data: null }))
  const token =
    'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaW5mbyI6eyJpZCI6MzgsInV1aWQiOiIiLCJuYW1lIjoidXNlcl8xMjEzOCIsIm5pY2tfbmFtZSI6IumfreiPnOmdkuW5tCIsImRlc2MiOiIiLCJhdmF0YXIiOiIiLCJwaG9uZSI6IiIsImVtYWlsIjoiIiwicmF0ZSI6MCwiYmFsYW5jZSI6MCwiY3JlYXRlZF9hdCI6IjIwMjEtMDktMDhUMTY6MzA6MjIuNjA4KzA4OjAwIn0sImV4cCI6MTYzMTI1NTM0OSwiaXNzIjoibGVla2JveCJ9.QLmWFeKyuP0IvsngTulgW4VTfZ5WDkrK0SZ45k-FYYw'
  document.querySelector('body').onclick = () => {
    ws.send(JSON.stringify({ type: 'LOGIN', data: token }))
  }
})
ws.addEventListener('error', console.log)
