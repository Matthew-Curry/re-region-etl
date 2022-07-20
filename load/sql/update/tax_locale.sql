-- assume ids and names are static
ON CONFLICT (tax_locale_id) DO UPDATE SET
    resident_desc = EXCLUDED.resident_desc,
    resident_rate = EXCLUDED.resident_rate,
    resident_month_fee = EXCLUDED.resident_month_fee,
    resident_year_fee = EXCLUDED.resident_year_fee,
    resident_pay_period_fee = EXCLUDED.resident_pay_period_fee,
    resident_state_rate = EXCLUDED.resident_state_rate,
    -- non-resident fields
    nonresident_desc = EXCLUDED.nonresident_desc,
    nonresident_rate = EXCLUDED.nonresident_rate,
    nonresident_month_fee = EXCLUDED.nonresident_month_fee,
    nonresident_year_fee = EXCLUDED.nonresident_year_fee,
    nonresident_pay_period_fee = EXCLUDED.nonresident_pay_period_fee,
    nonresident_state_rate = EXCLUDED.nonresident_state_rate;
