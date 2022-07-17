-- assume id and name are constant
ON CONFLICT (state_id) DO UPDATE SET
    single_deduction = EXCLUDED.single_deduction,
    married_deduction = EXCLUDED.married_deduction,
    single_exemption = EXCLUDED.single_exemption,
    married_exemption = EXCLUDED.married_exemption,
    dependent_exemption = EXCLUDED.dependent_exemption;
