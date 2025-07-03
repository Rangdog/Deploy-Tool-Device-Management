import { z } from 'zod'

export const createBillSchema = z.object({
  assetId: z.string({ required_error: 'Asset is required' }),
  assetName: z.string({ required_error: 'Asset name is required' }),
  categoryName: z.string({ required_error: 'Category name is required' }),
  cost: z.coerce.number({ required_error: 'Cost is required' }),
  description: z.string({ required_error: 'Description is required' }).min(1, 'Description is required'),
  status: z.enum(['Unpaid', 'Paid']),
  fileAttachment: z.instanceof(File).optional().or(z.literal('')),
  imageUpload: z.instanceof(File).optional().or(z.literal('')),
})

export type CreateBillFormType = z.infer<typeof createBillSchema>
