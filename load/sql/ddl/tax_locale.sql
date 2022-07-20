CREATE TABLE tax_locale (
    tax_locale_id INTEGER PRIMARY KEY,
    tax_locale VARCHAR( 50 ) NOT NULL,
    -- the county id is a foriegn key for the tax jurisdiction table
    CONSTRAINT fk_county
        FOREIGN KEY(county_id) 
	    REFERENCES county(county_id)
        ON DELETE CASCADE,

    county_id INTEGER NOT NULL,
    -- all metrics are not null. Use zero value in load if not applicable.
    -- resident fields
    resident_desc VARCHAR( 50 ) NOT NULL,
    resident_rate DECIMAL NOT NULL,
    resident_month_fee  DECIMAL NOT NULL,
    resident_year_fee  DECIMAL NOT NULL,
    resident_pay_period_fee  DECIMAL NOT NULL,
    resident_state_rate DECIMAL NOT NULL,
    -- non-resident fields
    nonresident_desc VARCHAR( 50 ) NOT NULL,
    nonresident_rate DECIMAL NOT NULL,
    nonresident_month_fee DECIMAL NOT NULL,
    nonresident_year_fee DECIMAL NOT NULL,
    nonresident_pay_period_fee DECIMAL NOT NULL,
    nonresident_state_rate DECIMAL NOT NULL
);
