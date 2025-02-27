import { useEffect, useRef } from 'react'
import { MessageList } from './components/MessageList'
import { NewMessageForm } from './components/NewMessageForm'
import 'bootstrap/dist/css/bootstrap.min.css';
import { Container } from 'react-bootstrap';
import { useMessageWebSocket } from './hooks/useMessageWebSocket';

const WS_URL = import.meta.env.MODE === 'development' 
  ? "ws://localhost:5000/api/ws"
  : `wss://${window.location.host}/api/ws`

console.log(window.location)

function App() {
  const {messages, sendMessage} = useMessageWebSocket(WS_URL)
  const ref = useRef(HTMLDivElement.prototype)

  useEffect(() => {
    ref.current.scrollIntoView()
  }, [messages])

  return (
    <Container >
      <MessageList messages={messages} />
      <div ref={ref}/> {/* Used to scroll down to the latest message */}
      <NewMessageForm sendMessage={sendMessage} />
    </Container>
  )
}

export default App
