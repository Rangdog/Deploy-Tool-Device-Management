import { httpRequest } from '@/utils'
import type { CreateBillRequest } from '../model/bill-types'

export const createBill = async (data: CreateBillRequest) => {
  return await httpRequest.post('/bills', data)
}
