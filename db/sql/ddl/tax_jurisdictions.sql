CREATE TABLE tax_jurisdiction (
    tax_jurisdiction_id INTEGER PRIMARY KEY,
    tax_jurisdiction VARCHAR( 50 ) NOT NULL,
    -- the county id is a foriegn key for the tax jurisdiction table
    CONSTRAINT fk_county
        FOREIGN KEY(county_id) 
	    REFERENCES county(county_id)
        ON DELETE CASCADE,

    county_id INTEGER NOT NULL,
    -- resident fields
    resident_desc VARCHAR( 50 ) NOT NULL,
    resident_rate DECIMAL ( 4 ),
    resident_month_fee  DECIMAL ( 4 ),
    resident_year_fee  DECIMAL ( 4 ),
    resident_pay_period_fee  DECIMAL ( 4 ),
    resident_state_rate DECIMAL ( 4 ),
    -- non-resident fields
    nonresident_desc VARCHAR( 50 ) NOT NULL,
    nonresident_rate DECIMAL ( 4 ),
    nonresident_month_fee  DECIMAL ( 4 ),
    nonresident_year_fee  DECIMAL ( 4 ),
    nonresident_pay_period_fee  DECIMAL ( 4 ),
    nonresident_state_rate DECIMAL ( 4 )
);
