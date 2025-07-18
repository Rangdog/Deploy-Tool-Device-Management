import { httpRequest } from '@/utils'

export const updateReadNotification = async (id: string) => {
  return await httpRequest.put(`/notifications/${id}`)
}
