# Protocols

- SDP (Session Description Protocol): messages with key/value pairs for signaling information for facilitating WebRTC sessions. These messages are sent out-of-band, as in they are separate from WebRTC sessions and can be managed by existing infrastructure.
- ICE (Interactive Connectivity Establishment): enables connection between two agents without going through a server by using NAT Traversal, TURN Agents and the STUN Protocol
- DTLS (Datagram Transport Layer Security): TLS over UDP, but in the case of ICE, there are no central certificate authorities; instead it validates a fingerprint from signaling, which depends on ICE agents. DTLS is used for DataChannel messages (not audio/video).
- TLS (Transport Layer Security): Encrypts communication by handshaking on encryption and session key(s) used for encryption, given the server's TLS certificate is valid (the client checks the server's certificate against a trusted Certificate Authority).
- UDP (User Datagram Protocol): Like TCP but no handshake, which allows packets to be lost.
- SRTP (Secure Real-Time Transport Protocol): Encrypted RTP (over UDP), which uses a complementary controller (RTCP or STCP if secured) to identify packet loss and out-of-order delivery. SRTP is used for audio/video messages.

# JavaScript

- RTCPeerConnection
  - RTCSessionDescription
  - RTCIceCandidate
- MediaStream