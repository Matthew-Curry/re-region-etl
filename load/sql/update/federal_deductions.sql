ON CONFLICT (single_deduction, married_deduction, head_deduction) DO UPDATE SET
    single_deduction = EXCLUDED.single_deduction,
    married_deduction = EXCLUDED.married_deduction,
    head_deduction = EXCLUDED.head_deduction;
