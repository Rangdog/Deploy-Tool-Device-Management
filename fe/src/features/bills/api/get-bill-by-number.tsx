import { httpRequest } from '@/utils'

export const getBillByNumber = async (billNumber: string) => {
  return await httpRequest.get(`/bills/${billNumber}`)
}
