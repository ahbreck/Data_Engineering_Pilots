CREATE TABLE IF NOT EXISTS covid_zip_weekly (
  zip_code CHAR(9),
  week_start DATE,
  week_end DATE,
  case_rate_weekly NUMERIC,
  percent_tested_positive_weekly NUMERIC,
  records_fetched_at TIMESTAMP DEFAULT now(),
  PRIMARY KEY (zip_code, week_start)
);