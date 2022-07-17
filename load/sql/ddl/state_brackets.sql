CREATE TABLE state_brackets (
    -- the state id is a foriegn key for the state table
    CONSTRAINT fk_state
        FOREIGN KEY(state_id) 
	    REFERENCES states(state_id)
        ON DELETE CASCADE,
    
    state_id SMALLINT NOT NULL,
    -- all metrics are not null. Use zero value in load if not applicable.
    single_rate DECIMAL NOT NULL,
    single_bracket INTEGER NOT NULL,
    married_rate DECIMAL NOT NULL,
    married_bracket INTEGER NOT NULL,
    CONSTRAINT ux_state_brackets UNIQUE (state_id, single_bracket, married_bracket)
);