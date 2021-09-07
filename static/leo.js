console.log('leojs inited')
window.addEventListener('load', () => {
  console.log('loaded')
})

ws = new WebSocket('ws://localhost:7003/api/stream?room_id=1')
ws.addEventListener('message', ev => console.log('message:', ev.data))
ws.addEventListener('close', console.log)
ws.addEventListener('open', ev => {
  console.log('open', ev)
  ws.send(JSON.stringify({ type: 'PING0', data: null }))
  window.wstimer = setInterval(() => {
    ws.send(JSON.stringify({ type: 'PING', data: null }))
  }, 15000)
  const token =
    'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaW5mbyI6eyJpZCI6MzUsInV1aWQiOiIiLCJvbWl0IjoiYjQ1Y2ZmZTA4NGRkM2QyMGQ5MjhiZWU4NWU3YjBmMjEiLCJuYW1lIjoibGl1ZnVsaW4iLCJuaWNrX25hbWUiOiLmlrDlkI3lrZflk6YiLCJkZXNjIjoiIiwiYXZhdGFyIjoiaHR0cHM6Ly90aGVzaHkuY2MvaW1nLzEucG5nIiwicGhvbmUiOiIxODYwMDAwOTk5OSIsImVtYWlsIjoibGVvb29AcXEuY29tIiwicmF0ZSI6NC42LCJiYWxhbmNlIjowLCJ1cGRhdGVkX2F0IjoiMjAyMS0wOS0wM1QxNjowNzowMC4yMTkrMDg6MDAiLCJjcmVhdGVkX2F0IjoiMjAyMS0wOS0wM1QxNjowNzowMC4yMiswODowMCJ9LCJleHAiOjE2MzEwOTk0MTAsImlzcyI6ImxlZWtib3gifQ.cnLfs6LQQVXhVYohaUfwQ0heIHtWRrFQcFxmfqd9qpU'
  document.querySelector('h2.title').addEventListener('click', () => {
    ws.send(JSON.stringify({ type: 'LOGIN', data: token }))
  })
})
ws.addEventListener('error', console.log)
