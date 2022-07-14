-- view to aggregate census data from county to the state level
CREATE OR REPLACE VIEW state_metrics AS
    WITH agg_metrics AS (
        SELECT 
            state.state_id,
            SUM(pop) as pop,
            SUM(male_pop) as male_pop,
            SUM(female_pop) as female_pop,
            SUM(median_income * pop) / SUM(pop) as average_median_income, 
            SUM(average_rent * pop) / SUM(pop) as average_rent, 
            SUM(commute * pop) / SUM(pop) as commute
        FROM county INNER JOIN state ON county.state_id = state.state_id
        GROUP BY state.state_id
    ) 

    SELECT 
        state.state_id,
        state.state_name,
        pop,
        male_pop,
        female_pop,
        average_median_income,
        average_rent, 
        commute
    FROM state INNER JOIN agg_metrics ON state.state_id = agg_metrics.state_id;
