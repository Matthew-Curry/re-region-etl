# re-region-etl
Dockerized ETL CLI tool to load source data for re-region-api into a Postgres DB from US census API and excel files produced by the Tax Foundation.

## Source Data
Taxation information is sourced to the app's database from datasets published by the Tax Foundation. This application is in no way affiliated or endorsed by the Tax Foundation. This data has been transformed and cleaned for storage in the database so it does not match the orginial form. For instance, the linkage of counties to local tax jurisdictions is performed by this application and is not a part of the Tax Foundation's orginial datasets. Also, the application performs the taxation estimates based on this data.

Tax foundation works are licensed under a Creative Commons Attribution NonCommercial 4.0 International License.

https://taxfoundation.org/copyright-notice/

Links to the original source data sets:

https://taxfoundation.org/publications/federal-tax-rates-and-tax-brackets/

https://taxfoundation.org/publications/state-individual-income-tax-rates-and-brackets/

https://taxfoundation.org/local-income-taxes-2019/

This application uses the Census Bureau Data API to source lifestyle and demographic data on regions, but is not endorsed or certified by the Census Bureau. This is accessed from the Census API at the county level, this applicaiton does the aggregation of those metrics to the state level.

