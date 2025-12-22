package clientDb

const QueryRefundBillno = `SELECT RetWarhouCode as yddh,'02' as rkfs FROM TB_Refund WHERE ISNULL(SendStatus,'') = '' 
                                      and Status = 51  and  CreateTime >= ?  and createTime < ? `

const UpdateRefund_Sql = ` UPDATE TB_Refund SET SendStatus = 1  where RetWarhouCode = ?`
