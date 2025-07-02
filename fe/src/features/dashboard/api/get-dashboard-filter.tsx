import { httpRequest } from '@/utils'

export interface DashboardFilterType {
  categoryId?: string | null
  departmentId?: string | null
  status?: string | null
  month?: number
  year?: number
  startDate?: string
  endDate?: string
  dateField?: 'purchase' | 'warranty' | 'created'
}

export const getDashboardAssetsFilter = async (value: DashboardFilterType) => {
  return await httpRequest.get('/assets/filter', {
    params: {
      ...value,
    },
  })
}
