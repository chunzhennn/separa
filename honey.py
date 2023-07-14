import socket
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect(("143.110.244.215", 2222))
banner = s.recv(1024)
print(banner)
s.send(b"SSH-1.9-OpenSSH_5.9p1\r\n")
data = s.recv(1024)
print(data)