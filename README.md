# gt-crawler

Search and crawl IGC tracks from the Web and store them
in a local filesystem, the metadata goes to postgres DB.

# Database: postgis

The metadata, track information, and actual geo data are stored
in Postgresql with GIS extension: postgis database.

## Data organisation

to be described


# Development

## Setup

* You need to have installed:
   * `Go`
   * `Docker`
* The project is organised into folders:
   * Top level folder contains this `README.md` file, `docker-compose.yml`. Here you run `docker-compose` commands.
   * `service_postgis` contains the `Dockerfile` needed to build and startup the `postgis` database.
   * `service_crawler` contains the Golang sources for the actual crawler, as well as the `Dockerfile` to build the image
   * `data_` folders are mapped to container volumes and represent:
      * `data_files` - data folder where crawler stores all crawled `.icg` files.
      * `data_pgadmin` - data folder for `pgAdmin 4`
      * `data_postgis` - data folder for `postgis` database


# Deployment

* Copy `.env.sample` to `.env` and fill that up with appropriate data for your deployment. **Note**: Do not store `.env` file in Git.
* Run once `docker network create net-internal` to create internal network for docker containers. This needs to be done only once.
* To build and start/stop all the services run:
   * `docker-compose up`
   * `docker-compose down`
* To connect to `pgAdmin 4` navigate your browser to `localhost:8080` on the deployment environment.
* The database is available on `postgis:5432`
* Use the credentials that have been set in `.env`


# Useful commands

## PostGIS

* `pg_ctl -D /usr/local/var/postgres start`
* `pg_ctl -D /usr/local/var/postgres stop`
* `createdb <dbname>` -- to create new DB from cmd line
* `dropdb <dbname>` -- to drop the db from the cmd line
* `psql <dbname>` -- to connect to a given DB


