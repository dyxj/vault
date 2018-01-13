sudo docker run -it --publish 80:80 -p 443:443 --name vaulttest --rm -v ~/dvolumes/vault/vault-autocert:/go/src/vault/vault-autocert  vaultgo:latest

sudo docker run -it --publish 8080:80 -p 6060:443 --name vaultexp --rm -v ~/dvolumes/vault/vault-autocert:/go/src/vault/vault-autocert vaultexp:latest

docker run -it --publish 80:80 -p 443:443 --name vaulttest --rm -v //E/GoPaths/GoPath1/src/vault/vault-autocert:/go/src/vault/vault-autocert  vaultgo:latest