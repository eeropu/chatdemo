import ListGroup from 'react-bootstrap/ListGroup'
import { Message } from '../hooks/useMessageWebSocket'

type MessageListProps = {
  messages: Message[]
}

export const MessageList = ({ messages } : MessageListProps) => {
  return (
    <div>
      <ListGroup variant="flush">
        { messages.map(msg => <ListGroup.Item key={msg._id}>{msg.content}</ListGroup.Item>)}
      </ListGroup>
    </div>
  )
}