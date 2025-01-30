# SYNCPLACE


In order to run everything in your local environment (You should have Docker + Docker Compose installed)

Run this command to create SyncPlace SERVER, UI, and DB instances that the application needs.

"docker-compose up -d"  
After some time,you should get :

[+] Running 3/3  
 ✔ Container syncplace-postgres-1       Started          2.1s  
 ✔ Container syncplace-syncplace-srv-1  Started          0.8s  
 ✔ Container syncplace-syncplace-ui-1   Started          1.0s  

After everything is up, connect through the browser : http://localhost:3000/login

