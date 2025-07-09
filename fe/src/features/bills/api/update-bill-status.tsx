import { httpRequest } from '@/utils'

export const updateBillStatus = async (billNumber: string, status: 'Paid' | 'Unpaid') => {
  return await httpRequest.patch(`/bills/${billNumber}`, status)
}
