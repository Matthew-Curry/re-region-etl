CREATE TABLE state_brackets (
    -- the state id is a foriegn key for the state table
    CONSTRAINT fk_state
        FOREIGN KEY(state_id) 
	    REFERENCES state(state_id)
        ON DELETE CASCADE,
    
    state_id SMALLINT NOT NULL,
    single_rate DECIMAL,
    single_bracket INTEGER,
    married_rate DECIMAL,
    married_bracket INTEGER,
    CONSTRAINT ux_state_brackets UNIQUE (state_id, single_bracket, married_bracket)
);