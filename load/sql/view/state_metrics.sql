-- view to aggregate census data from county to the state level
DROP VIEW IF EXISTS state_metrics;

CREATE VIEW state_metrics AS
    WITH agg_metrics AS (
        SELECT 
            states.state_id,
            SUM(pop) AS pop,
            SUM(male_pop) AS male_pop,
            SUM(female_pop) AS female_pop,
            CAST(ROUND(SUM(median_income * pop) / SUM(pop)) AS BIGINT) AS average_median_income, 
            CAST(ROUND(SUM(average_rent * pop) / SUM(pop)) AS BIGINT) AS average_rent, 
            SUM(commute * pop) / SUM(pop) AS commute
        FROM county INNER JOIN states ON county.state_id = states.state_id
        GROUP BY states.state_id
    ) 

    SELECT 
        states.state_id,
        states.state_name,
        pop,
        male_pop,
        female_pop,
        average_median_income,
        average_rent, 
        commute
    FROM states INNER JOIN agg_metrics ON states.state_id = agg_metrics.state_id
    -- no null record
    WHERE states.state_id != 32767;
