package clientDb

const FullWorkloadSQL = `
WITH InRegBase AS(
    SELECT em.EmployeeName as OpName, a.DeptCode as Dept, a.InStoreNo as billno, b.ProductInfoID as prod_id, b.Amount as qty, b.PurchasePrice as BuyPrice
    FROM TB_Warehouse a JOIN TB_WarehouseDetail b ON a.InStoreNo = b.InStoreNo JOIN TB_Employee em ON em.HRCode = a.Auditor
    WHERE a.AuditorDate >= ? AND a.AuditorDate <= ? AND a.Status = 21 AND a.WarehouseType IN (0,1) AND a.Source = 1
),
InboundBase AS (
    SELECT a.MEnName as OpName, CASE WHEN A.DeptCode='200404' THEN '200346' ELSE A.DeptCode END as Dept, a.billno, b.prod_id, b.qty, b.BuyPrice
    FROM T_Prod_Enter a INNER JOIN T_ProdEnter_Detail b ON a.Reg_id = b.Reg_id
    WHERE a.billstate IN ('41','51') AND b.IsVoid = 0 AND a.EnterDate >= ? AND a.EnterDate <= ?
    UNION ALL
    SELECT em.EmployeeName, a.DeptCode, a.InStoreNo, b.ProductInfoID, b.Amount, b.PurchasePrice
    FROM TB_Warehouse a JOIN TB_WarehouseDetail b ON a.InStoreNo = b.InStoreNo JOIN TB_Employee em ON em.HRCode = a.Confirmer
    WHERE a.Status = 21 AND a.WarehouseType IN (0,1) AND a.Source = 1 AND a.ConfirmerDate >= ? AND a.ConfirmerDate <= ?
),
OutboundBase AS (
    SELECT a.BLMakerName as OpName, CASE WHEN a.TreasuryDepartment='200404' THEN '200346' ELSE a.TreasuryDepartment END as Dept, a.DepartmentCollarID as billno, b.ProductInfoID as prod_id, b.Amount, b.RealUnitPrice as Price
    FROM TB_DepartmentCollar a INNER JOIN TB_DepartmentCollarDetail b ON a.DepartmentCollarID = b.DepartmentCollarID
    WHERE a.Status IN ('21','51','61') AND b.IsVoid = 0 AND a.TreasuryDepartment IN ('200346','200418','200632','200624') AND a.BLDate >= ? AND a.BLDate <= ?
    UNION ALL
    SELECT em.EmployeeName, dr.StoreCode, CONVERT(VARCHAR, dr.DeliveryID) + '-' + CONVERT(VARCHAR, dr.DetailSort), dr.ProductInfoID, dr.ThisAmount, dr.UnitPrice
    FROM TB_DeliveryApplyDetailRecord dr JOIN TB_DeliveryApply d ON dr.DeliveryID = d.DeliveryID JOIN TB_Employee em ON em.HRCode = dr.StoreAuditor
    WHERE dr.IsVoid = 0 AND d.Source = '1' AND d.IsStockGoods = '0' AND d.[Type] IN ('1','2') AND d.[Status] IN (61,71,41,81,22,91,19,29,99) AND dr.CreateTime >= ? AND dr.CreateTime <= ?
),
ReturnBase AS (
    SELECT d.EmployeeName as OpName, CASE WHEN a.StorehouseID='200404' THEN '200346' ELSE a.StorehouseID END as Dept, a.ReturnID as billno, b.ProductInfoID as prod_id, b.Amount, b.UnitPrice as Price
    FROM TB_ReturnPurchase a INNER JOIN TB_ReturnPurchaseDetail b ON a.ReturnID = b.ReturnID INNER JOIN TB_Employee d ON a.BLMaker = d.HRCode
    WHERE a.BLDate >= ? AND a.BLDate <= ? AND a.StorehouseID IN ('200346','200418','200632','200624') AND a.Status = 21 AND b.IsVoid = 0
    UNION ALL
    SELECT em.EmployeeName, a.StorehouseID, a.ReturnID, b.ProductInfoID, b.Amount, b.UnitPrice
    FROM TB_ReturnPurchase a INNER JOIN TB_ReturnPurchaseDetail b ON a.ReturnID = b.ReturnID INNER JOIN TB_Employee em ON em.HRCode = a.Auditor
    WHERE a.Source in (0,3) AND a.Status in (21,61,51) AND AuditorDate >= ? AND AuditorDate <= ?
),
PurchaseBase AS (
    SELECT em.EmployeeName as OpName, '采购业务' as Dept, a.PurchaseSummaryID as billno, 0 as prod_id, 0 as qty, 0 as Price
    FROM TB_PurchaseSummary a LEFT JOIN TB_Employee em ON em.HRCode = Auditor
    WHERE a.Status in (91,71,61,81,92,51) AND a.AuditorDate >= ? AND a.AuditorDate < ? AND HRcODE <> '0000900000'
),
PushBase AS (
    SELECT em.EmployeeName as OpName, '采购业务' as Dept, a.BillNo as billno, 0 as prod_id, 0 as qty, 0 as Price
    FROM TB_PurchaseSummaryPush a LEFT JOIN TB_Employee em ON em.HRCode = a.Pusher
    WHERE a.PushDate >= ? AND a.PushDate < ?
),
RefundBase AS (
    SELECT em.EmployeeName as OpName, a.TargetStorehouseID as Dept, a.RetWarhouID as billno, b.ProductInfoID as prod_id, b.Amount as qty, b.UnitPrice as Price
    FROM TB_Refund A INNER JOIN TB_RefundDetail B ON A.RetWarhouID = B.RetWarhouID INNER JOIN TB_Employee em ON em.HRCode = A.SourceAuditor
    WHERE A.Status IN (51) AND A.SourceAuditorDate >= ? AND A.SourceAuditorDate <= ?
),
FinalResults AS (
    -- 0:验收, 1:出库, 2:退货, 3:登记, 4:采购, 5:催货, 6:二级库退库
    SELECT OpName, 0 as OpType, Dept, COUNT(DISTINCT billno) as Bills, SUM(qty*BuyPrice) as Amt, (SELECT SUM(c) FROM (SELECT COUNT(DISTINCT prod_id) as c FROM InboundBase t2 WHERE t2.OpName=t1.OpName AND t2.Dept=t1.Dept GROUP BY billno) x) as Specs FROM InboundBase t1 GROUP BY OpName, Dept
    UNION ALL
    SELECT OpName, 1, Dept, COUNT(DISTINCT billno), SUM(Amount*Price), (SELECT SUM(c) FROM (SELECT COUNT(DISTINCT prod_id) as c FROM OutboundBase t2 WHERE t2.OpName=t1.OpName AND t2.Dept=t1.Dept GROUP BY billno) x) FROM OutboundBase t1 GROUP BY OpName, Dept
    UNION ALL
    SELECT OpName, 2, Dept, COUNT(DISTINCT billno), SUM(Amount*Price), (SELECT SUM(c) FROM (SELECT COUNT(DISTINCT prod_id) as c FROM ReturnBase t2 WHERE t2.OpName=t1.OpName AND t2.Dept=t1.Dept GROUP BY billno) x) FROM ReturnBase t1 GROUP BY OpName, Dept
    UNION ALL
    SELECT OpName, 3, Dept, COUNT(DISTINCT billno), SUM(qty*BuyPrice), (SELECT SUM(c) FROM (SELECT COUNT(DISTINCT prod_id) as c FROM InRegBase t2 WHERE t2.OpName=t1.OpName AND t2.Dept=t1.Dept GROUP BY billno) x) FROM InRegBase t1 GROUP BY OpName, Dept
    UNION ALL
    SELECT OpName, 4, Dept, COUNT(DISTINCT billno), 0, 0 FROM PurchaseBase GROUP BY OpName, Dept
    UNION ALL
    SELECT OpName, 5, Dept, COUNT(DISTINCT billno), 0, 0 FROM PushBase GROUP BY OpName, Dept
    UNION ALL
    SELECT OpName, 6, Dept, COUNT(DISTINCT billno), SUM(qty*Price), (SELECT SUM(c) FROM (SELECT COUNT(DISTINCT prod_id) as c FROM RefundBase t2 WHERE t2.OpName=t1.OpName AND t2.Dept=t1.Dept GROUP BY billno) x) FROM RefundBase t1 GROUP BY OpName, Dept
)
SELECT 
    f.OpName as OperatorName, f.OpType as OperationType, f.Dept as StorehouseCode,
    f.Specs as SpecCount, f.Bills as BillCount, f.Amt as TotalAmount,
    d.DepartmentName as FallbackName
FROM FinalResults f
LEFT JOIN TB_Department d ON f.Dept = d.DeptCode
ORDER BY  OperatorName,OperationType
`
