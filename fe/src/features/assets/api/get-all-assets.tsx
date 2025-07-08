import { httpRequest, type ApiResponse } from '@/utils'
import type { AssetsType } from '../view-all-assets'

export const getAllAssets = async (): Promise<ApiResponse<AssetsType[]>> => {
  return await httpRequest.get('/assets/filter')
}
