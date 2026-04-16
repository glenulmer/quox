# Excel Quote Generation Notes (`~/redef/code/xl-quote.go`)

## Scope
This documents how `CreateExcelQuote` works in `/home/glen/redef/code/xl-quote.go`, what it depends on, and how to make it easier to maintain.

Entry point from request flow:
- `DownloadReview` calls `CreateExcelQuote(client, user, slim != 0)` in `/home/glen/redef/code/pg-review.go:242-253`.

## High-level flow
`CreateExcelQuote` builds an `.xlsx` by loading a pre-existing template (`assets/masters/ExcelStyle.xlsx`), writing values into fixed cells, deleting rows/columns depending on client state, then saving to `assets/work/<customer> initial overview(.slim).xlsx`.

Main phases:
1. Open template, load named styles from template `formats` sheet.
2. Fill header and meta rows (name, DOB, sick-cover label).
3. Choose benefit section list (`fat` vs `slim`).
4. Write left-side benefit labels and section titles.
5. Write one plan per fixed column (`D`, `F`, `H`, `J`, `L`).
6. Delete unused plan columns.
7. Write footer note text.
8. Remove rows based on segment/slim/dependants.
9. Remove `formats` sheet, add top image, save file.

## Function-by-function behavior

### `CreateExcelQuote` (`xl-quote.go:12-80`)
- Uses hardcoded template path: `assets/masters/ExcelStyle.xlsx`.
- Calls `Excel(template)` from `/home/glen/redef/zlib/xl/xl.go`.
  - `Excel()` opens workbook and defaults to first sheet (`GetSheetName(0)`), so write target depends on workbook sheet order.
- Calls `ex.AddStyles("formats", "B2", "J100")`.
  - This scans that rectangle on sheet `formats` and builds `ex.Styles[<cell text>] = <style id>`.
- Customer header:
  - `A3 = custname` (fallback `"<user.name>s.Customer"` if blank).
  - `A4 = "Date of birth: <dob>"`.
  - `A8` gets sick-pay line only when `client.sickcover > 0`.
- Segment-specific row text:
  - If not employee, sets `A22 = "Your monthly cost"`.
- Section choice:
  - `sections = client.bens.sections.fat`, or `.slim` when `slim == true`.
- Left side labels:
  - `lastRow := WriteBenefitsAndTips(...)`.
- Plan columns:
  - Iterates `client.chosen` and writes up to 5 plans (fixed by `PlanColumns`).
- Trailing column deletion:
  - Computes first column after last used plan and deletes that same column index 10 times (shift-left behavior used intentionally).
- Footer note:
  - Writes at `A(lastRow+2)` using style `SlimNote`.
- Row deletions:
  - Non-employee: delete row 21 (`MonthlyWithEmpRow`).
  - Dependants: template has slots for 10, deletes unused slots.
  - Non-slim student with zero sick cover: delete rows 8-9.
  - Slim: delete rows 6-11.
- Cleanup and save:
  - Deletes sheet `formats`.
  - Adds header image at `A1` from `assets/klexpats.jpg`.
  - Saves file under `assets/work/`.

### `WriteBenefitsAndTips` (`xl-quote.go:82-103`)
- Starts at fixed `BenefitRow = 24`.
- For each section:
  - Writes section name in column A.
  - Merges section title block in column A over section height.
  - Writes benefit labels in column B.
- Title/value style keys are name-derived (`<section.name>Title/Label`).
- Special-case benefit id `1` uses style `DeductibleLabel`.
- Appends final merged `Helpful extra tips` block and returns last used row.

### `WritePlanColumn` (`xl-quote.go:105-186`)
For one chosen plan at one fixed column:
- Reads plan from `client.chosen[ix]` and `client.plans.byId`.
- Adds provider logo image: `assets/logos/<lower(provider name)>.jpg` at row 2.
- Optional top note at row 3 using style name from DB (`plan.topnote.style`).
- Rows 4-9: provider/plan title, reference+commission, HIC/PVN/sick split, sick-after text.
- Rows 11-20: per-dependant cost lines; also rewrites column A labels with dependant name+age.
- Row 21/22: total and “you pay”.
- Rows 24+: benefit offers:
  - Calculated benefits (magic ids `1` deductible, `2` no-claims) or family map values.
  - Product-specific override via `FindOverride`.
- Tip lines (`client.bens.tips[family.id]`) written under benefit table with `TipText` style.

### `FindOverride` (`xl-quote.go:188-195`)
- Scans active products and returns first override string found in `client.bens.overrides[(benefit, product)]`.

## Template contract (hard dependencies)

### Workbook structure dependencies
Hardcoded assumptions:
- Template exists at `assets/masters/ExcelStyle.xlsx`.
- First worksheet is the data sheet (currently `Sheet1`).
- A helper sheet named `formats` exists and is deletable.
- Template is designed around plan columns `D/F/H/J/L` with spacer columns between them.

Observed workbook sheets:
- `Sheet1`
- `formats`

### Style-key dependencies
`AddStyles("formats", "B2", "J100")` expects these keys in `formats` sheet cell text (exact spelling):
- `YouPay`
- `BasicsTitle`, `BasicsLabel`, `BasicsValue`
- `DeductibleLabel`, `DeductibleValue`
- `Hospital StayTitle`, `Hospital StayLabel`, `Hospital StayValue`
- `OutpatientTitle`, `OutpatientLabel`, `OutpatientValue`
- `ExtrasTitle`, `ExtrasLabel`, `ExtrasValue`
- `DentalTitle`, `DentalLabel`, `DentalValue`
- `TipTitle`, `TipText`
- `Commission`
- `TopNote`
- `SlimNote`

If any style key is absent or renamed in template, writes still happen but style id becomes `0` (silent formatting loss).

### Cell/row/column dependencies
Hardcoded coordinates/constants:
- Header/meta rows: `A3`, `A4`, `A8`, `A22`
- Dependants area starts at row `11`, max 10 rows
- Totals rows `21` and `22`
- Benefits start row `24`
- Plan columns fixed to 5 slots: `D`, `F`, `H`, `J`, `L`

### Asset dependencies
- Top banner: `assets/klexpats.jpg`
- Provider logos: `assets/logos/<lower(provider)>.jpg`
  - Depends on exact provider naming convention mapping to filenames.

### Data dependencies from DB/model
- `client.bens.sections.fat/slim` ordering and names drive table layout and style lookup.
- Section names must align with template style prefixes (`<section name>Title/Label/Value`).
- Benefit ids `1` and `2` are magic values with special logic.
- `plan.topnote.style` must match a style key loaded from template.
- `family.offers`, `bens.overrides`, `bens.tips` must be populated for expected output.

## Risks and brittle points
1. Hidden coupling to template internals
- Much behavior depends on sheet names, row numbers, column letters, style key strings, and image filenames.
- Most of these are not validated before writing.

2. Silent failures
- `AddPicture` logs errors but generation continues.
- `Save` return value is ignored.
- Missing style keys silently degrade formatting.

3. Dead/ineffective check
- `if ex == nil` after `Excel(template)` is effectively dead for missing file, because `Excel()` panics on open failure.

4. Sheet-order coupling
- No explicit sheet selection in `CreateExcelQuote`; it writes to workbook sheet index 0.
- Reordering sheets in template can break output target.

5. Magic numbers and IDs
- Row/column constants and benefit ids are hardcoded and distributed.
- Harder to review when layout changes.

6. Truncation without explicit signal
- More than 5 chosen plans are silently ignored.

7. Filename safety
- Output filename includes raw `client.custname`; risky for invalid filename characters.

## Dependency-tracking concerns (template usage)
Your concern is valid: this implementation has many implicit template contracts and no single place that declares them.

Practical improvements without changing architecture style:
1. Add a template contract document next to template
- Example: `assets/masters/ExcelStyle.contract.md` listing required sheets, style keys, critical cells, and required image naming.

2. Add explicit runtime validator (single function)
- Validate once before writing:
  - required sheets exist (`Sheet1`, `formats`)
  - required style keys loaded
  - critical cells exist/are reachable
  - banner/logo files exist for selected providers
- Fail fast with one concise error summary.

3. Add a lightweight golden test
- Generate one known quote fixture.
- Assert key cells/style keys/sheet count in resulting workbook XML.
- This catches accidental template edits early.

4. Use template named ranges for anchors
- Replace some raw addresses with named anchors (`ClientNameCell`, `BenefitsStart`, etc.).
- Keeps layout maintainable when rows move in Excel.

## Organization/readability improvements
Suggested target organization for `xl-quote.go`:

1. Split into clear phases with small functions
- `initQuoteWorkbook`
- `writeHeader`
- `writeBenefitsLabels`
- `writePlanColumns`
- `pruneLayout`
- `writeFooter`
- `finalizeAndSave`

2. Centralize layout contract in one struct
- One `quoteLayout` value containing all rows/columns/style keys.
- Removes scattered magic numbers and improves diff readability.

3. Replace magic benefit IDs with named constants near domain model
- `BenefitDeductible = 1`, `BenefitNoClaim = 2` in one canonical place.

4. Make style key access explicit and checked
- Helper like `mustStyle("SlimNote")` or prevalidated style map.
- Avoid silent style id `0` fallback.

5. Separate data computation from sheet writes
- Compute a `quoteViewModel` per plan first.
- Then render model to sheet in one pass.
- Makes logic testable without workbook I/O.

6. Isolate destructive layout ops
- Keep all row/column deletions in one function with comments on ordering assumptions.

7. Make truncation explicit
- If plans > 5, add one note in workbook/log so truncation is visible.

## Bottom line
The current code works by mutating a tightly preformatted workbook and depends on many implicit contracts. The main maintainability issue is not Excel itself; it is that those contracts are scattered and mostly unchecked. Centralizing and validating the template contract would remove most fragility without changing the overall approach.
