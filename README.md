To host the project :
clone the project

````
git clone https://github.com/zeeenku/go-chat.git
```


first make the deploy.sh file executable via 
```
chmod +x deploy.sh
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


enable and start service
```
sudo systemctl daemon-reload
```

```
sudo systemctl enable go-chat.service
```

```
sudo systemctl start go-chat.service
```



make the nginx file
```
sudo nano /etc/nginx/sites-available/go-chat
```

and put 
```
# HTTP server block to redirect traffic to HTTPS
server {
    listen 80;
    server_name go-chat.zeenku.com;  # Replace with your domain name

    # Redirect all HTTP traffic to HTTPS
    return 301 https://$server_name$request_uri;
}

# HTTPS server block with SSL configuration
server {
    listen 443 ssl;
    server_name go-chat.zeenku.com;  # Replace with your domain name


    # WebSocket handling
    location /ws {
        proxy_pass http://127.0.0.1:7777;  # Go WebSocket server address
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_read_timeout 3600s;  # Increase WebSocket connection timeout
    }

    # Static file handling
    location / {
        root /var/www/go-chat/static;  # Correct path to your static files
        index index.html index.htm;
    }

    # Log files for debugging (optional)
    error_log /var/log/nginx/go-chat_error.log;
    access_log /var/log/nginx/go-chat_access.log;
}
```

then

```
sudo ln -s /etc/nginx/sites-available/go-chat /etc/nginx/sites-enabled/
```

```
sudo nginx -t
```

I recommand adding ssl via certbot


```
sudo certbot --nginx -d yourdomain.com
```



then run the deploy.sh file

```
./deploy.sh
```
