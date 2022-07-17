CREATE TABLE county (
	county_id INTEGER PRIMARY KEY,
    county_name VARCHAR ( 50 ) NOT NULL,
    -- the state id is a foriegn key for the state table
    CONSTRAINT fk_state
        FOREIGN KEY(state_id) 
	    REFERENCES states(state_id)
        ON DELETE CASCADE,

    state_id SMALLINT NOT NULL,
	pop INTEGER NOT NULL,
    male_pop INTEGER NOT NULL,
    female_pop INTEGER NOT NULL,
    -- median income and average rent must be big int
    -- for state aggregation view
    median_income BIGINT NOT NULL,
    average_rent BIGINT NOT NULL, 
    -- all metrics are not null. Use zero value in load if not applicable.
    commute INTEGER NOT NULL
);