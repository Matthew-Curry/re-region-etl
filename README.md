# re-region-etl
Dockerized ETL CLI tool to load source data for the re-region-api (https://github.com/Matthew-Curry/re-region-api) into a Postgres DB from the Census Bureau Data API and excel files produced by the Tax Foundation.

## Usage
First, build the docker container with the database parameters passed in as environment variables:

```docker build -t re-region-etl:latest --build-arg RE_REGION_ETL_USER='{USER}' --build-arg RE_REGION_ETL_PASSWORD='{PASSWORD}' --build-arg RE_REGION_DB='{DB_NAME}' --build-arg DB_PORT='{DB_PORT}' --build-arg DB_HOST='{DB_HOST}' .```

Running the container will output the options for the CLI:
```docker run re-region-etl:latest``` 

Output:

```Usage of ./re-region-etl:
  -c    c, Clears existing tables resulting from the ETL stages to run. Will only take effect if l option is provided to trigger the ETL
  -l    l, Runs the ETL code to load the tables
  -v    v, Runs SQL to define the views.
 ```
In addition to the flags, one or more stages can be provided to define which stages of the ETL run. Other than the first stage (federal tax data) each stage is dependent on the previous (i.e passing stage 4 will load 2, 3 and 4 out of necessity)

Pass in the flags and stages to run the ETL as needed.

## Project Structure and Data Processing
**data:** Holds source excel files from the Tax Foundation <br>
**extract:** Holds extractors that take data from sources, transform and load to in memory structures. Those sources are the afformentioned data files as well as the Census Bureau Data API. <br>
**load:** Holds database engine with functionality to create tables, insert data, and define views on the re-region database. Also holds "sql" folder with all DDL, insert, update and create view statements. <br>
**logging:** Package holds my implementation of an aggregated logger with public methods for different log levels that is used throughout the app <br>
**sourceFileUtils:** Package holds method used to read in the source excel files. <br>
**main.go:** Defines the CLI interface. Holds a core "runETL" method that uses the extractors and the DB engine to load the database. The ETL will be processed as per the provided args and stages.

The most interesting part of the process is in the localTaxExtractor. A tax jurisdiction (sourced from the local tax info files from the Tax Foundation) oftentimes, but not always, is a county, or the identifier for the jursidiction contains some or all of the county name. So, a tax jursidction is linked to a county by way of fuzzy matching using the open source package github.com/paul-mannino/go-fuzzywuzzy. Further, the core loop in the extractor spawns goroutines to perform this linking in parallel, which drastically improved the runtime.

## Source Data
Taxation information is sourced to the app's database from datasets published by the Tax Foundation. It is also from these datasets that the app sources local tax jurisdictions. The taxation estimates the API provides are based on the information given by these data sets, but it is the application building those estimates. The estimates are a simplification and should not be taken as definitive taxation information or advice. The linking between the federal, state, and local tax data sets is done by the applicaiton. Notably, the application matches tax jurisdictions to counties using an open source package implementing fuzzy matching functionality. Those links are not provided by any source dataset and are not guarenteed to be accurate. This application is in no way affiliated or endorsed by the Tax Foundation.

Tax foundation works are licensed under a Creative Commons Attribution NonCommercial 4.0 International License.

https://taxfoundation.org/copyright-notice/

Links to the original source data sets:

Published in 2022: https://taxfoundation.org/publications/federal-tax-rates-and-tax-brackets/

Published in 2022: https://taxfoundation.org/publications/state-individual-income-tax-rates-and-brackets/

Published in 2019: https://taxfoundation.org/local-income-taxes-2019/

This application uses the Census Bureau Data API to access data from the 2019 American Community Survey to source survery statistic information to the API. The app is not endorsed or certified by the Census Bureau. Data is accessed from the Census API at the county level; this applicaiton does the aggregation of those metrics to the state level.

