import asyncio
import ssl
from aioquic.asyncio import serve
from aioquic.quic.configuration import QuicConfiguration
from aioquic.quic.events import DatagramFrameReceived, StreamDataReceived
from aioquic.h3.connection import H3_ALPN
import logging

logging.basicConfig(level=logging.DEBUG)

class WebTransportProtocol:
    def __init__(self, connection, stream_handler):
        self._connection = connection
        self._stream_handler = stream_handler

    def connection_made(self, transport):
        self._transport = transport
        logging.info("Connection made")

    def connection_lost(self, exc):
        logging.info("Connection lost")

    def quic_event_received(self, event):
        print("event", event)
        if isinstance(event, StreamDataReceived):
            logging.info(f"Stream data received: {event.data}")
            self._connection.send_stream_data(event.stream_id, event.data, event.end_stream)
        elif isinstance(event, DatagramFrameReceived):
            logging.info(f"Datagram received: {event.data}")
            # Echo the received datagram back to the client
            self._connection.send_datagram_frame(event.data)

    def datagram_received(self, data, addr):
        print(f"Datagram received from {addr}")
        # logging.info(f"Datagram received from {addr}: {data}")
        

    async def data_received(self, data):
        await super().data_received(data)
        if isinstance(self._quic, QuicProtocol):
            for stream_id, buffer in self._quic._events[DataReceived]:
                data = buffer.read()
                response = Response(status_code=200, headers=[('content-length', str(len(data)))], content=data)
                data = self.conn.send(response)
                self._quic.transmit_data(stream_id, data)

async def main():
    configuration = QuicConfiguration(
        alpn_protocols=H3_ALPN,
        is_client=False,
        max_datagram_frame_size=65536,
    )
    configuration.load_cert_chain("cert.pem", "key.pem")

    await serve(
        "0.0.0.0",
        4433,  # Use a different port than your HTTPS server
        configuration=configuration,
        create_protocol=WebTransportProtocol,
    )
    await asyncio.Future()  # run forever

if __name__ == "__main__":
    asyncio.run(main())
