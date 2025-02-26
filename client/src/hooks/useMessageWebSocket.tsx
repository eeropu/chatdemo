import { useEffect, useRef, useState } from "react";

type IncomingMessage = {
  type: string;
  data: Message[];
};

export type Message = {
  content: string;
  _id: string;
};

export const useMessageWebSocket = (url: string) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const ws = useRef<WebSocket | null>(null);

  useEffect(() => {
    ws.current = new WebSocket(url);
    ws.current.onopen = () => console.log("ws connection opened");
    ws.current.onclose = () => console.log("ws connection closed");
    ws.current.onmessage = (event) => {
      const incomingData: IncomingMessage = JSON.parse(event.data);
      if (incomingData.type === "snapshot") {
        setMessages(incomingData.data);
      } else {
        setMessages((m) => m.concat(incomingData.data));
      }
    };

    const wsCurrent = ws.current;
    return () => {
      wsCurrent.close();
    };
  }, [url]);

  const sendMessage = (content: string) => {
    if (ws.current) {
      ws.current.send(JSON.stringify({ content }));
    }
  };

  return { messages, sendMessage };
};
