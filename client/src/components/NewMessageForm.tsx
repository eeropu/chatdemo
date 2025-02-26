import { useState } from "react"
import Form from "react-bootstrap/form"
import Button from "react-bootstrap/button"
import { Col, Row } from "react-bootstrap"

type NewMessageFormProps = {
  sendMessage: (content: string) => void
}

export const NewMessageForm = ({ sendMessage }: NewMessageFormProps) => {
  const [content, setContent] = useState<string>("")

  const handleSubmit = () => {
    sendMessage(content)
    setContent("")
  }

  return (
    <Form style={{padding:'20px', position:'sticky', bottom:0,backgroundColor:'white',borderTop:'solid 1px gray'}}>
      <Row>
        <Col>
          <Form.Control 
            as="textarea"
            rows={3}
            value={content}
            onChange={(event) => setContent(event.target.value)}
            placeholder="New message..."
            style={{background:'lightgray'}}
            />
        </Col>
        <Col xs={1} style={{display: 'flex', alignItems:'center'}}>
          <Button onClick={handleSubmit}>Send</Button>
        </Col>
      </Row>
    </Form>
  )
}