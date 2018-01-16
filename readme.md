# Vault #
A file encryption service that's uses Go, Docker, MongoDB and Javascript.  
[**Click here for a demo**](https://file.darrenyxj.com/)   
**It does not store files**.
Files are encrypted upon upload and an encrypted file will be available for download.   
To decrypt, upload encrypted file and enter password.  


## Setup and Run ##
Using docker is recommended. But can also be done without.
* Docker is required.
* Modify docker-compose.yml, example provided with some elaboration.
* docker-compose up.

## Documentation ##
There are two services defined in the docker-compose. A webapp and a database(MongoDB).


## Future Development ##
This application was built mainly for practice. To test out deployment etc.
* Optimization
* Native application with GUI. Provides an option to run without a web browser.