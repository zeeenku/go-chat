To host the project first make the deploy.sh file executable via

```
chmod +x deploy.sh
```

then run it

```
./deploy.sh
```


working directory (example):

```
/var/www/go-chat
```
make a service file for example in this path:

```
sudo nano /etc/systemd/system/go-chat.service
```

and put 
```
[Unit]
Description=Go WebSocket Chat Application
After=network.target

[Service]
ExecStart=/var/www/go-chat/go-chat
WorkingDirectory=/var/www/go-chat
User=www-data
Group=www-data
Restart=always
Environment=GO_ENV=production
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target

```
make the nginx file
```
sudo nano /etc/nginx/sites-available/go-chat
```

and put 
```
server {
    listen 80;
    server_name yourdomain.com;  # Replace with your domain name or IP address

    location / {
        root /var/www/html;  # Path to your static files (index.html)
        index index.html index.htm;
    }

    location /ws {
        proxy_pass http://127.0.0.1:7777;  # Replace with your Go app's WebSocket address
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    error_log /var/log/nginx/go-chat_error.log;
    access_log /var/log/nginx/go-chat_access.log;
}
```

I recommand adding ssl via certbot
sudo certbot --nginx -d yourdomain.com
