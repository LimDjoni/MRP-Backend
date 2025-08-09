package stockfuel

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	ListStockFuel(sortFilter StockFuelSummary) ([]StockFuelSummary, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ListStockFuel(sortFilter StockFuelSummary) ([]StockFuelSummary, error) {
	var rawResults []StockFuelSummary

	if sortFilter.Month == "" {
		return nil, fmt.Errorf("month is required")
	}

	// Load Asia/Jakarta location (UTC+7)
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, err
	}

	// Parse the month in UTC first
	parsed, err := time.Parse("2006-01", sortFilter.Month) // e.g. "2025-07"
	if err != nil {
		return nil, err
	}

	// Convert parsed time to Asia/Jakarta timezone with day 1, time zeroed
	parsedInLoc := time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, loc)

	// Start date: 1st day of month in UTC+7
	startDate := parsedInLoc.Format("2006-01-02")

	// End date: last day of month in UTC+7
	endDate := parsedInLoc.AddDate(0, 1, -1).Format("2006-01-02")
	prevMonthStart := parsedInLoc.AddDate(0, -1, 0).Format("2006-01-02") // 1st day prev month
	prevMonthEnd := parsedInLoc.AddDate(0, 0, -1).Format("2006-01-02")   // last day prev month

	const stockFuelSQL = `
		WITH RECURSIVE
	dates AS (
		SELECT generate_series(?::date, ?::date, interval '1 day') AS report_date
	),
	fuel_data AS (
		SELECT
			TO_DATE(COALESCE(fr.tanggal, fr.tanggal_awal), 'YYYY-MM-DD') AS report_date,
			SUM(CASE WHEN fr.shift = 'Shift 1' THEN fr.total_refill ELSE 0 END) AS mrp_day,
			SUM(CASE WHEN fr.shift = 'Shift 2' THEN fr.total_refill ELSE 0 END) AS mrp_night,
			SUM(fr.total_refill) AS total_fuel_out
		FROM fuel_ratios fr
		WHERE fr.status = true
			AND TO_DATE(COALESCE(fr.tanggal, fr.tanggal_awal), 'YYYY-MM-DD') BETWEEN ? AND ?
		GROUP BY report_date
	),
	fuel_in_data AS (
		SELECT
			TO_DATE(fi.date, 'YYYY-MM-DD') AS report_date,
			SUM(fi.qty_now) AS total_fuel_in,
			SUM(CASE WHEN fi.vendor = 'MJSU' THEN fi.qty_now ELSE 0 END) AS mjsu,
			SUM(CASE WHEN fi.vendor = 'PPP' THEN fi.qty_now ELSE 0 END) AS ppp,
			SUM(CASE WHEN fi.vendor = 'SADP' THEN fi.qty_now ELSE 0 END) AS sadp
		FROM fuel_ins fi
		WHERE TO_DATE(fi.date, 'YYYY-MM-DD') BETWEEN ? AND ?
		GROUP BY report_date
	),
	explicit_first_stock AS (
		SELECT date::date AS report_date, stock FROM adjust_stocks WHERE deleted_at IS NULL
	),
	prev_month_end_stock AS (
		SELECT
			COALESCE((SELECT stock FROM adjust_stocks WHERE date = ?), 0)
			+ COALESCE((SELECT SUM(qty_now) FROM fuel_ins WHERE date >= ? AND date <= ?), 0)
			- COALESCE((
				SELECT SUM(total_refill)
				FROM fuel_ratios
				WHERE status = true
				AND TO_DATE(COALESCE(tanggal, tanggal_awal), 'YYYY-MM-DD') >= ?
				AND TO_DATE(COALESCE(tanggal, tanggal_awal), 'YYYY-MM-DD') <= ?
			), 0) AS stock
	),
	daily_data AS (
		SELECT
			d.report_date,
			COALESCE(f.mrp_day, 0) AS mrp_day,
			COALESCE(f.mrp_night, 0) AS mrp_night,
			COALESCE(f.total_fuel_out, 0) AS total_fuel_out,
			COALESCE(fi.total_fuel_in, 0) AS fuel_in,
			COALESCE(fi.mjsu, 0) AS mjsu,
			COALESCE(fi.ppp, 0) AS ppp,
			COALESCE(fi.sadp, 0) AS sadp,
			SUM(COALESCE(f.total_fuel_out, 0)) OVER (ORDER BY d.report_date) AS mtd_fuel_out,
			efs.stock AS explicit_stock
		FROM dates d
		LEFT JOIN fuel_data f ON f.report_date = d.report_date
		LEFT JOIN fuel_in_data fi ON fi.report_date = d.report_date
		LEFT JOIN explicit_first_stock efs ON efs.report_date = d.report_date
	),
	recursive_stock AS (
		-- First day: use explicit stock or prev month end stock
		SELECT
			dd.report_date,
			dd.mrp_day,
			dd.mrp_night,
			dd.total_fuel_out,
			dd.fuel_in,
			dd.explicit_stock,
			dd.mtd_fuel_out,
			dd.mjsu,
			dd.ppp,
			dd.sadp,
			COALESCE(dd.explicit_stock, pmes.stock, 0) AS first_stock,
			COALESCE(dd.explicit_stock, pmes.stock, 0) + dd.fuel_in - dd.total_fuel_out AS end_stock
		FROM daily_data dd, prev_month_end_stock pmes
		WHERE dd.report_date = (SELECT MIN(report_date) FROM daily_data)

		UNION ALL

		-- Subsequent days: use explicit stock or previous day end_stock
		SELECT
			dd.report_date,
			dd.mrp_day,
			dd.mrp_night,
			dd.total_fuel_out,
			dd.fuel_in,
			dd.explicit_stock,
			dd.mtd_fuel_out,
			dd.mjsu,
			dd.ppp,
			dd.sadp,
			COALESCE(dd.explicit_stock, rs.end_stock) AS first_stock,
			COALESCE(dd.explicit_stock, rs.end_stock) + dd.fuel_in - dd.total_fuel_out AS end_stock
		FROM daily_data dd
		JOIN recursive_stock rs ON dd.report_date = rs.report_date + INTERVAL '1 day'
	)
	SELECT
		report_date,
		first_stock,
		mrp_day AS day,
		mrp_night AS night,
		total_fuel_out AS total,
		total_fuel_out AS grand_total,
		fuel_in,
		end_stock,
		mtd_fuel_out AS mtd_consump,
		mjsu,
		ppp,
		sadp
	FROM recursive_stock
	ORDER BY report_date;
		`

	r.db.Exec("SET TIME ZONE 'Asia/Jakarta'")

	err = r.db.Raw(stockFuelSQL,
		startDate, endDate,
		startDate, endDate,
		startDate, endDate,
		prevMonthStart,
		prevMonthStart, prevMonthEnd,
		prevMonthStart, prevMonthEnd,
	).Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Format the date into "Selasa, 01 Juli 2025"
	months := map[time.Month]string{
		time.January:   "Januari",
		time.February:  "Februari",
		time.March:     "Maret",
		time.April:     "April",
		time.May:       "Mei",
		time.June:      "Juni",
		time.July:      "Juli",
		time.August:    "Agustus",
		time.September: "September",
		time.October:   "Oktober",
		time.November:  "November",
		time.December:  "Desember",
	}

	days := map[time.Weekday]string{
		time.Sunday:    "Minggu",
		time.Monday:    "Senin",
		time.Tuesday:   "Selasa",
		time.Wednesday: "Rabu",
		time.Thursday:  "Kamis",
		time.Friday:    "Jumat",
		time.Saturday:  "Sabtu",
	}

	var results []StockFuelSummary
	for _, rraw := range rawResults {
		rraw.Date = fmt.Sprintf("%s, %02d %s %d",
			days[rraw.ReportDate.Weekday()],
			rraw.ReportDate.Day(),
			months[rraw.ReportDate.Month()],
			rraw.ReportDate.Year(),
		)
		results = append(results, rraw)
	}

	return results, nil
}
