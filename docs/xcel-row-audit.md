# Excel Plan-Column Audit (Definitive, Redef Parity)

Standard used for all answers: **correct = matches redef behavior/output contract**.

## Row 2
1. Is the value correct?  
No.
2. If not, did you compute it with new code?  
No.
3. Or did you read it from qvars?  
No.
4. Is it correctly formatted?  
No (content is missing).
5. If not, how does redef format it?  
Redef inserts provider logo image in row 2 of each plan column.
6. What is current output?  
Nothing is written to row 2.

## Row 3
1. Is the value correct?  
No (not redef-parity).
2. If not, did you compute it with new code?  
Yes.
3. Or did you read it from qvars?  
No.
4. Is it correctly formatted?  
No (style source differs from redef).
5. If not, how does redef format it?  
Redef writes top note and applies the note's own style id.
6. What is current output?  
Top note text from lookup with fixed style key `TopNote`.

## Row 4
1. Is the value correct?  
Yes.
2. If not, did you compute it with new code?  
N/A.
3. Or did you read it from qvars?  
No.
4. Is it correctly formatted?  
Yes.
5. If not, how does redef format it?  
N/A.
6. What is current output?  
`<provider> / <plan name>`.

## Row 5
1. Is the value correct?  
No (commission source differs from redef).
2. If not, did you compute it with new code?  
Yes.
3. Or did you read it from qvars?  
No.
4. Is it correctly formatted?  
Partly (plain `Ref x/y` shape is right; value source is not redef-parity).
5. If not, how does redef format it?  
`Ref <plan.id>/<AnnualComm(...)>`.
6. What is current output?  
`Ref <planId>/<EuroFlatFromCent(row.commission).Int64()>`.

## Row 6
1. Is the value correct?  
No (value path differs from redef).
2. If not, did you compute it with new code?  
Yes.
3. Or did you read it from qvars?  
No.
4. Is it correctly formatted?  
Yes (shape matches redef: amount or amount + preex).
5. If not, how does redef format it?  
`<hic>` or `<hic> + <preex>`, all as plain cent-strings.
6. What is current output?  
`row.price - pvn - sick`, optional ` + pX`.

## Row 7
1. Is the value correct?  
No (value path differs from redef, and preex handling differs).
2. If not, did you compute it with new code?  
Yes.
3. Or did you read it from qvars?  
No.
4. Is it correctly formatted?  
Yes (plain cent-string).
5. If not, how does redef format it?  
`Str(pvn)` from chosen-price structure.
6. What is current output?  
`Str(addon.base + addon.surcharge)` from current row addon.

## Row 8
1. Is the value correct?  
No (value path differs from redef, and preex handling differs).
2. If not, did you compute it with new code?  
Yes.
3. Or did you read it from qvars?  
No.
4. Is it correctly formatted?  
Yes (plain cent-string).
5. If not, how does redef format it?  
`Str(sick)` from chosen-price structure.
6. What is current output?  
`Str(addon.base + addon.surcharge)` from current row addon.

## Row 9
1. Is the value correct?  
No.
2. If not, did you compute it with new code?  
Yes.
3. Or did you read it from qvars?  
No.
4. Is it correctly formatted?  
No.
5. If not, how does redef format it?  
Redef writes sick waiting-period text (`After()` output).
6. What is current output?  
Addon pick label or addon category text.

## Row 10
1. Is the value correct?  
Yes (redef also leaves this row unused in plan top block).
2. If not, did you compute it with new code?  
N/A.
3. Or did you read it from qvars?  
N/A.
4. Is it correctly formatted?  
N/A.
5. If not, how does redef format it?  
N/A.
6. What is current output?  
No write.

## Rows 11-20
1. Is the value correct?  
No (not redef-parity for data path guarantees).
2. If not, did you compute it with new code?  
Yes.
3. Or did you read it from qvars?  
No direct precomputed chosen-dependant price field.
4. Is it correctly formatted?  
Yes (plain cent-string now matches redef text format).
5. If not, how does redef format it?  
Redef writes `Str(dep price)` from precomputed dependant choice prices.
6. What is current output?  
Per dependant, recomputed row with swapped birth/vision, then `Str(depRow.price)`.

## Row 21
1. Is the value correct?  
No for employee parity (value path differs).  
For non-employee row is removed elsewhere as requested behavior.
2. If not, did you compute it with new code?  
Yes.
3. Or did you read it from qvars?  
No.
4. Is it correctly formatted?  
Yes (plain cent-string).
5. If not, how does redef format it?  
Redef writes `Str(t)` where `t` is chosen monthly total.
6. What is current output?  
Employee only: `Str(row.price)`.

## Row 22
1. Is the value correct?  
No (employer-match path and accumulation path are not redef-parity).
2. If not, did you compute it with new code?  
Yes.
3. Or did you read it from qvars?  
No.
4. Is it correctly formatted?  
Yes (plain cent-string).
5. If not, how does redef format it?  
Redef uses `youpay = t - employerPays`, then adds dependant totals.
6. What is current output?  
`yourPay` from local employer-share formula + dependant recompute totals.

## Section 0 - Value 1 (Deductible, benefit id 1)
1. Is the value correct?  
No (not guaranteed redef-parity because upstream row source differs).
2. If not, did you compute it with new code?  
Yes.
3. Or did you read it from qvars?  
No.
4. Is it correctly formatted?  
Yes (`Str(...)` format matches redef style).
5. If not, how does redef format it?  
Redef writes `Str(chosen.deductible)`.
6. What is current output?  
`Str(row.deduct)`.

## Section 0 - Value 2 (No-claim bonus, benefit id 2)
1. Is the value correct?  
No (not guaranteed redef-parity because upstream row source differs).
2. If not, did you compute it with new code?  
Yes.
3. Or did you read it from qvars?  
No.
4. Is it correctly formatted?  
Yes (note+amount style now matches redef shape).
5. If not, how does redef format it?  
Redef writes `Str(plan.noclaims.note, chosen.ncbonus)`.
6. What is current output?  
`Str(plan.nc.note, row.noClaims)` or `Str(row.noClaims)` when note is empty.
