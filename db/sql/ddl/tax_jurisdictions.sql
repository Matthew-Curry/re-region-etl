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
    resident_rate DECIMAL,
    resident_month_fee  DECIMAL,
    resident_year_fee  DECIMAL,
    resident_pay_period_fee  DECIMAL,
    resident_state_rate DECIMAL,
    -- non-resident fields
    nonresident_desc VARCHAR( 50 ) NOT NULL,
    nonresident_rate DECIMAL,
    nonresident_month_fee  DECIMAL,
    nonresident_year_fee  DECIMAL,
    nonresident_pay_period_fee  DECIMAL,
    nonresident_state_rate DECIMAL
);
