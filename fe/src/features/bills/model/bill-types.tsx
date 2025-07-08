export interface BillType {
  id: number
  billNumber: string
  assetId: number
  description: string
  status: 'Unpaid' | 'Paid'
  // categoryId?: number
  companyId: number
  amount: number
  createdBy: number
  createAt: string
  updateAt: string
  purchaseDate?: string
  warrantyExpiry?: string
  qrUrl?: string
  fileAttachment?: string
  imageUpload?: string

  creator?: {
    id: number
    fullName: string
    email: string
    avatar?: string
  }

  assets?: {
    id: number
    assetName: string
    serialNumber?: string
    cost: number
    status?: string
    category: {
      id: number
      categoryName: string
    }
    department?: {
      id: number
      departmentName: string
    }
    purchaseDate?: string
    warrantyExpiry?: string
  }
}

export type CreateBillRequest = {
  assetId: number
  // cost: number
  description: string
  // categoryId?: number
  status?: string
  fileAttachment?: File
  imageUpload?: File
}

export interface BillFilterType {
  billNumber: string | ''
  categoryId: string | null
  companyId: string | null
  status: 'Unpaid' | 'Paid' | null
}

export interface MonthlyBillSummary {
  month: string
  year: number
  totalBills: number
  totalAmount: number
  statusBreakdown: {
    status: string
    count: number
    amount: number
  }[]
  categoryBreakdown: {
    categoryName: string
    count: number
    amount: number
  }[]
}
