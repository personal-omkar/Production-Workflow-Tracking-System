package dashboard

import (
	"fmt"
	"strings"
	"time"
)

func BuildVendorOrderStatusQuery(month, userType string) string {
	if strings.TrimSpace(month) == "" {
		month = time.Now().Format("2006-01")
	}

	entityColumn := "v.vendor_name"
	joinClauses := map[string]string{
		"pending":   "JOIN kb_extension ke ON ke.id = kd.kb_extension_id JOIN vendors v ON v.id = ke.vendor_id",
		"latestTx":  "JOIN vendors v ON v.id = ke.vendor_id",
		"approved":  "JOIN vendors v ON v.id = ke.vendor_id",
		"rejected1": "JOIN vendors v ON v.id = ke.vendor_id",
		"rejected2": "JOIN vendors v ON v.id = ke.vendor_id",
	}

	if strings.EqualFold(userType, "operator") || strings.EqualFold(userType, "customer") {
		entityColumn = "c.compound_name"
		joinClauses["pending"] = "JOIN kb_extension ke ON ke.id = kd.kb_extension_id JOIN compounds c ON c.id = kd.compound_id"
		joinClauses["latestTx"] = "JOIN kb_data kd ON kd.kb_extension_id = ke.id JOIN compounds c ON c.id = kd.compound_id"
		joinClauses["approved"] = "JOIN kb_data kd ON kd.kb_extension_id = ke.id JOIN compounds c ON c.id = kd.compound_id"
		joinClauses["rejected1"] = "JOIN kb_data kd ON kd.kb_extension_id = ke.id JOIN compounds c ON c.id = kd.compound_id"
		joinClauses["rejected2"] = "JOIN kb_data kd ON kd.kb_extension_id = ke.id JOIN compounds c ON c.id = kd.compound_id"
	}

	query := fmt.Sprintf(`
WITH pending_cte AS (
  SELECT %[1]s AS vendor_name, COUNT(*) AS pending
  FROM kb_data kd
  %[2]s
  WHERE TO_CHAR(kd.created_on, 'YYYY-MM') = '%[3]s' AND ke.status = 'pending'
  GROUP BY vendor_name
),

latest_tx AS (
  SELECT DISTINCT ON (kt.kb_root_id)
    kt.kb_root_id, ke.status AS extension_status, kr.status AS root_status,
    %[1]s AS vendor_name, kt.created_on
  FROM kb_transaction kt
  JOIN kb_root kr ON kt.kb_root_id = kr.id
  JOIN kb_data kg ON kg.id = kr.kb_data_id
  JOIN kb_extension ke ON ke.id = kg.kb_extension_id
  %[4]s
  WHERE TO_CHAR(kt.created_on, 'YYYY-MM') = '%[3]s'
  ORDER BY kt.kb_root_id, kt.created_on DESC
),

approved_cte AS (
  SELECT %[1]s AS vendor_name, COUNT(*) AS approved
  FROM kb_extension ke
  %[5]s
  WHERE ke.status = 'approved'
  GROUP BY vendor_name
),

rejected_cte AS (
  SELECT vendor_name, COUNT(*) AS rejected
  FROM (
    SELECT %[1]s AS vendor_name, ke.id AS reject_id
    FROM kb_extension ke
    %[6]s
    WHERE ke.status = 'reject' AND TO_CHAR(ke.created_on, 'YYYY-MM') = '%[3]s'

    UNION

    SELECT %[1]s AS vendor_name, kr.id AS reject_id
    FROM kb_root kr
    JOIN kb_data kg ON kg.id = kr.kb_data_id
    JOIN kb_extension ke ON ke.id = kg.kb_extension_id
    %[7]s
    WHERE kr.status = '-1' AND TO_CHAR(kr.created_on, 'YYYY-MM') = '%[3]s'
  ) AS all_rejections
  GROUP BY vendor_name
),

all_names_cte AS (
  SELECT vendor_name FROM pending_cte
  UNION SELECT vendor_name FROM approved_cte
  UNION SELECT vendor_name FROM rejected_cte
  UNION SELECT vendor_name FROM latest_tx
),

final AS (
  SELECT
    an.vendor_name,
    COALESCE(p.pending, 0) AS pending,
    COALESCE(a.approved, 0) AS approved,
    COUNT(lt.*) FILTER (WHERE lt.extension_status = 'InProductionProcess') AS in_progress,
    COUNT(lt.*) FILTER (WHERE lt.root_status = '2') AS quality,
    COUNT(lt.*) FILTER (WHERE lt.root_status = '3') AS packing,
    COUNT(lt.*) FILTER (WHERE lt.root_status = '4') AS completed,
    COALESCE(r.rejected, 0) AS rejected
  FROM all_names_cte an
  LEFT JOIN pending_cte p ON p.vendor_name = an.vendor_name
  LEFT JOIN approved_cte a ON a.vendor_name = an.vendor_name
  LEFT JOIN rejected_cte r ON r.vendor_name = an.vendor_name
  LEFT JOIN latest_tx lt ON lt.vendor_name = an.vendor_name
  GROUP BY an.vendor_name, p.pending, a.approved, r.rejected
)

SELECT * FROM final
WHERE pending > 0 OR approved > 0 OR rejected > 0 OR in_progress > 0 OR quality > 0 OR packing > 0 OR completed > 0
ORDER BY vendor_name;
`, entityColumn, joinClauses["pending"], month,
		joinClauses["latestTx"], joinClauses["approved"],
		joinClauses["rejected1"], joinClauses["rejected2"])

	return query
}

func BuildInProgressOrderQuery(prodLineID int) string {
	return fmt.Sprintf(`
	WITH latest_tx AS (
		SELECT DISTINCT ON (kt.kb_root_id)
			kt.kb_root_id,
			kt.prod_process_line_id,
			ppl.prod_process_id,
			kr.running_no,
			kr.status
		FROM kb_transaction kt
		JOIN kb_root kr ON kt.kb_root_id = kr.id
		JOIN prod_process_line ppl ON ppl.id = kt.prod_process_line_id
		WHERE ppl.prod_line_id = %d
		ORDER BY kt.kb_root_id, kt.created_on DESC
		)
		SELECT 
		COUNT(*) FILTER (WHERE prod_process_id = 1) AS scheduled_count,
		COUNT(*) FILTER (
			WHERE CAST(running_no AS INTEGER) = 0 
			AND prod_process_id NOT IN (1, 2)
		) AS in_progress_count
		FROM latest_tx;
	`, prodLineID)
}
