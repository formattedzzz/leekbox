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
    'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaW5mbyI6eyJpZCI6MzUsInV1aWQiOiIiLCJuYW1lIjoibGl1ZnVsaW4iLCJuaWNrX25hbWUiOiLmlrDlkI3lrZflk6YiLCJkZXNjIjoiIiwiYXZhdGFyIjoiaHR0cHM6Ly90aGVzaHkuY2MvaW1nLzEucG5nIiwicGhvbmUiOiIxODYwMDAwOTk5OSIsImVtYWlsIjoibGVvb29AcXEuY29tIiwicmF0ZSI6NC42LCJiYWxhbmNlIjowLCJjcmVhdGVkX2F0IjoiMjAyMS0wOS0wOFQxNjozMDowNy43NTMrMDg6MDAifSwiZXhwIjoxNjMxMTc2MjU0LCJpc3MiOiJsZWVrYm94In0.j80f34LaHHNJuIJxRg01yZyX-J9EXqMCWIB48on5PaA'
  document.querySelector('body').onclick = () => {
    ws.send(JSON.stringify({ type: 'LOGIN', data: token }))
  }
})
ws.addEventListener('error', console.log)
