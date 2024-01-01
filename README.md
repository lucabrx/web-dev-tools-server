# Web Dev Tools API 

## Project Environment Variables
```
DB_URL=postgres://postgres:postgres@localhost:5432/web_dev_tools
SERVER_ADDRESS=:8080
CLIENT_ADDRESS=http://localhost:3000
RESEND_API_KEY=
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
AWS_ACCESS_KEY=
AWS_SECRET_KEY=
```

## Resources
- [Go](https://golang.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://www.docker.com/)
- [AWS](https://aws.amazon.com/)
- [Resend](https://resend.com/)
- [migrate-tool](https://github.com/golang-migrate/migrate)
- [neon db](https://neon.tech/)
- [Vultr](https://www.vultr.com/?ref=9577975-8H)
- [Nginx](https://www.nginx.com/)

## Project outline
- users -> add tools to favorites, suggest tools
- tools -> paginated list of tools with search
- auth -> login with GitHub and magic link
- admin -> add tools to the database, approve suggested tools

## How to run
### Running PostgreSQL in Docker
```shell 
docker run --name web_dev_tools -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
### OR 
task postgres-container
```
or check neon db 

### Migrating the database
```shell
migrate -path ./migrations -database postgres://postgres:postgres@localhost:5432/web_dev_tools || <pg db> up
### OR
task migrate-up ### DB_URL from .env or $DB_URL
```

### Running the server
```shell
go run ./cmd/api
### OR
task run-server
```

## Installing and Running project on VPS (Vultr)
Free $100 credit for 30 days [here](https://www.vultr.com/?ref=9577975-8H)


### adding web address to VPS (porkbun)
- add A record to VPS ip address
- add api subdomain to VPS ip address (host)
- answer -> ip address of VPS

### Installing Go
```shell
Install go on VPS
```shell
# download go
wget https://golang.org/dl/go1.21.4.linux-amd64.tar.gz
# extract go
sudo tar -C /usr/local -xzf go1.21.4.linux-amd64.tar.gz
# set path in profile (~./bashrc)
export PATH=$PATH:/usr/local/go/bin
# refresh profile and check go version
source ~/.profile
go version # go1.21.4 linux/amd64

# download taskfile
wget https://github.com/go-task/task/releases/download/v3.4.3/task_linux_amd64.tar.gz
# extract taskfile
tar -xvf task_linux_amd64.tar.gz
```

### Setting up Nginx
```shell
sudo apt install nginx
sudo systemctl start nginx
```
Make proxy pass to the server (localhost:8080)
```shell
vi /etc/nginx/sites-available/default
```

```
server {
    listen 80;

    location / {
        proxy_pass http://localhost:8080;  # Port of your Go server
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```
```shell
# check nginx config
sudo nginx -t
# restart nginx
sudo systemctl restart nginx
```

### Clone project (don't forget to add .env file)
```shell
git clone <repo>
cd <repo>
go build -o /bin/server ./cmd/api
```

### Setup systemd service (run server always)
```shell
sudo vi /etc/systemd/system/web-dev-tools.service
```
```
[Unit]
Description=Web Dev Tools API

[Service]
ExecStart=/root/src/web-dev-tools-server/bin/api
WorkingDirectory=/root/src/web-dev-tools-server
User=root
Group=root
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```
enable service
```shell
 sudo systemctl enable /etc/systemd/system/web-dev-tools.service
```
start service
```shell
sudo systemctl start web-dev-tools.service
```
check status
```shell
sudo systemctl status web-dev-tools.service
```

### allow port 80 (ufw)
```
sudo ufw status
sudo ufw allow 80
``

### Installing Certbot (SSL) 
```
- add server_name to nginx config
```
server_name api.web-dev-tools.xyz;
```
```shell
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx

```
You should see a message that says:
`Successfully deployed certificate for api.web-dev-tools.xyz to /etc/nginx/sites-enabled/default
Congratulations! You have successfully enabled HTTPS on https://api.web-dev-tools.xyz`

### dry run
```shell
sudo certbot renew --dry-run
```
certbot should change and add some lines to nginx config for ssl but if not we need to expose port 443 

last thing is to  `ufw allow 443` and `sudo systemctl restart nginx` so we can use https