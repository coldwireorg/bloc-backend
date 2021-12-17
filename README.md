# Bloc backend
## Description

bloc is a project of federated cloud storage with strong encryption.<br>
The goal is to have a stable, simple to use and secure place to store our files and to not even worry about the loss of an instance of the network.

## Security

1. All the users have their own Private / Public keypair generated using secp256k1
2. Each files have their own encryption key, theses keys are encrypted usign the user public key.
3. The files are encrypted using Xchacha20-poly1305.
4. The private key is encrypted with Xchacha20-poly1305 and the encryption key is a 32bytes argon2 derivation of the user's password.

Currently all the opperations are done on the server side, we are considering the possibility to build a WASM library in Rust to makes all the operations on the user side, which would makes easier to trust the servers but would be impossible for users who are disabling javascript (with TOR for exemple)

## Features
*because we are currently in early developpement, the federation system is not implemented yet, the protocol will be written in Rust in another repo*

### Available features
- User creation / authentification (only using username/password system to avoid collecting metadatas)
- File uploading / downloading
- File deleting
- Favorite system (add the files to a list of favorites)
- Private file sharing (the file encryption key is re-encrypted using the receiver public key)

### Planned features
- Making possible to scale this server on nomad/kubernetes cluster 
- Public file sharing (sharing a file using a link, this one expose the file's encryption key!)
- Authentification from others server of the federated network
- Sharing files over the federated network
- Automated backup over the network (encrypted files are sliced into many little parts and sent to a swarm of servers)

## Installation

*note: before running the backend, setup a postgresql server, you can see the tables in `database/sql/tables.sql`*

### Environment variables:
```sh
SERVER_DOMAIN=coldwire.org # the domain name to use for the cookies 
SERVER_HOST=0.0.0.0 # the ip on which the web server listen
SERVER_PORT=3000 # the port of the web server
SERVER_HTTPS=false # set to true if the your server use SSL (IT MUST BE ACTIVATED IN PRODUCTION)
DB_ADDRESS=127.0.0.1 # address of the mongodb database
DB_PORT=27017 # port of the database
DB_NAME=bloc # name of the database
STORAGE_DIR=/tmp/bloc/files # path to your file server/storage dir where all the files will be put
STORAGE_QUOTA=1024 # Limit for users in Mb
```

### With docker
```sh
git clone https://github.com/coldwireorg/bloc-backend.git
cd bloc-backend

# build image
docker build -t coldwireorg\bloc-backend .

# run container
docker run -it --rm -p 1500:3000
# visit 127.0.0.1:3000
```

### Using git / For development

Before, be sure to have NodeJS and Go installed !

```sh
# clone the repository
git clone https://github.com/coldwireorg/bloc-backend.git
cd bloc-backend

go mod tidy # install modules
go run main.go # run server
```

## Building
### Docker

```sh
# Go to the root dir
cd bloc-backend
# build the image
docker build . -t coldwire/bloc-backend
```

### Manually
```sh
# Go to the root dir
cd bloc
# build program
go build main.go
# then you can upload your "main" binary to your server :)
# Don't forget to put the "view" directory in the same path as the binary, because then the front-end will not work
```

## License

This project is licensed under the [NPOSL 3.0](https://opensource.org/licenses/NPOSL-3.0) License.