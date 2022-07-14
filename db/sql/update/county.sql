-- assume ids and names are static
ON CONFLICT (county_id) DO UPDATE SET
	pop = EXCLUDED.pop,
    male_pop = EXCLUDED.male_pop,
    female_pop = EXCLUDED.female_pop,
    median_income = EXCLUDED.median_income,
    average_rent = EXCLUDED.average_rent,
    commute = EXCLUDED.commute;
    