# re-region-etl
Dockerized ETL CLI tool to load source data for re-region-api into a Postgres DB from US census API and excel files produced by the Tax Foundation.

## Usage
First, build the docker container with the database parameters passed in as environment variables:

```docker build -t re-region-etl:latest --build-arg RE_REGION_ETL_USER='{USER}' --build-arg RE_REGION_ETL_PASSWORD='{PASSWORD}' --build-arg RE_REGION_DB='{DB_NAME}' --build-arg DB_PORT='{DB_PORT}' --build-arg DB_HOST='{DB_HOST}' .```

Running the container will output the options to the CLI:
```docker run re-region-etl:latest``` 

Output:

```Usage of ./re-region-etl:
  -c    c, Clears existing tables resulting from the ETL stages to run. Will only take effect if l option is provided to trigger the ETL
  -l    l, Runs the ETL code to load the tables
  -v    v, Runs SQL to define the views.
 ```
Pass in the flags to run the ETL as needed.

## Project Structure
data: Holds source excel files from the Tax Foundation <br>
extract: Holds extractors that take data from sources, transform and load to in memory structures. Those sources are the afformentioned data files as well as the Census Bureau Data API.
load: Holds database engine with functionality to create tables, insert data, and define views on the re-region database. Also holds "sql" folder with all DDL, insert, update and create view statements.
logging: Package holds my implementation of an aggregated with public methods for different log levels that is used throughout the app
sourceFileUtils: Package holds method used to read in the source excel files.
main.go: Defines the CLI interface. Holds a core "runETL" method that will process the ETL as per the provided args and stages.



## Source Data
Taxation information is sourced to the app's database from datasets published by the Tax Foundation. This application is in no way affiliated or endorsed by the Tax Foundation. This data has been transformed and cleaned for storage in the database so it does not match the orginial form. For instance, the linkage of counties to local tax jurisdictions is performed by this application and is not a part of the Tax Foundation's orginial datasets. Also, the application performs the taxation estimates based on this data.

Tax foundation works are licensed under a Creative Commons Attribution NonCommercial 4.0 International License.

https://taxfoundation.org/copyright-notice/

Links to the original source data sets:

https://taxfoundation.org/publications/federal-tax-rates-and-tax-brackets/

https://taxfoundation.org/publications/state-individual-income-tax-rates-and-brackets/

https://taxfoundation.org/local-income-taxes-2019/

This application uses the Census Bureau Data API to source lifestyle and demographic data on regions, but is not endorsed or certified by the Census Bureau. This is accessed from the Census API at the county level, this applicaiton does the aggregation of those metrics to the state level.

