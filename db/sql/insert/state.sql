INSERT INTO state(
    state_id, 
    state_name,
    single_deduction,
    married_deduction,
    single_exemption,
    married_exemption,
    dependent_exemption
    ) 
-- initial row of values for the default state record
VALUES (?, ?, ?, ?, ?, ?, ?), 