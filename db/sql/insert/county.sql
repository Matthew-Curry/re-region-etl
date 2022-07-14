INSERT INTO county 
    (
    county_id, 
    county_name,
    state_id,
    pop,
    male_pop,
    female_pop,
    median_income,
    average_rent,
    commute
    ) 
-- initial row of values for the default county record
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?), 