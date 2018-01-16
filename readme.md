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

It would be good to create a compose file for prod and adding the -f flag to run. Check out [docker docs](https://docs.docker.com/compose/extends/#example-use-case) for more info.

## Documentation ##
There are two services defined in docker-compose.yml. A webapp and a database(MongoDB).  
* The webapp service is the Go web application.  
    * Can be ran in dev and prod mode.
    * When running in different modes do configure volumes accordingly.
* Database used to store number of files encrypted and decrypted.


## Future Development ##
This application was built mainly for practice. To test out deployment etc.
* Optimization
* Native application with GUI. Provides an option to run without a web browser.