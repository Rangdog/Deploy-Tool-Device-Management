import { httpRequest } from '@/utils'

export const createBill = async (data: FormData) => {
  console.log('🚀 ~ createBill ~ data:', data)
  return await httpRequest.post('/bills', data, {
    'Content-Type': 'multipart/form-data',
  })
}
