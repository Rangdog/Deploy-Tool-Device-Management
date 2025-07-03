import { httpRequest } from '@/utils'
import type { BillFilterType } from '../model/bill-types'

export const getBillsFilter = async (filters: BillFilterType) => {
  return await httpRequest.get('/bills/filter', {
    params: {
      ...filters,
    },
  })
}
